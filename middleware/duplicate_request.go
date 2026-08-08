package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	httpx "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	redisutil "github.com/Kevin-Jii/tower-go/utils/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	redislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	duplicateRequestKeyPrefix = "tower:duplicate-request:"
	duplicateRequestLockTTL   = 60 * time.Second
	duplicateRequestCooldown  = time.Second
	duplicateRedisTimeout     = 500 * time.Millisecond
)

var releaseDuplicateRequestScript = redislib.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
end
return 0
`)

// DuplicateRequestMiddleware rejects an identical write request while the
// first request is still running. A short cooldown also covers very fast APIs
// whose first response finishes between the two click events.
func DuplicateRequestMiddleware() gin.HandlerFunc {
	guard := newDuplicateRequestGuard(duplicateRequestLockTTL, duplicateRequestCooldown)
	return duplicateRequestMiddleware(guard)
}

func duplicateRequestMiddleware(guard *duplicateRequestGuard) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldRejectDuplicateRequest(c.Request) {
			c.Next()
			return
		}

		fingerprint, err := duplicateRequestFingerprint(c)
		if err != nil {
			logging.LogWarn("计算请求防重指纹失败，已跳过本次防重", zap.Error(err))
			c.Next()
			return
		}

		lockKey := duplicateRequestKeyPrefix + fingerprint
		token, distributed, acquired := guard.acquire(c.Request.Context(), lockKey)
		if !acquired {
			c.Header("Retry-After", "1")
			c.Header("X-Duplicate-Request", "true")
			httpx.ErrorApp(c, apicode.DuplicateOperation.WithMessage("相同请求正在处理中，请勿重复操作"))
			c.Abort()
			return
		}

		defer guard.releaseAfter(lockKey, token, distributed)
		c.Next()
	}
}

func shouldRejectDuplicateRequest(r *http.Request) bool {
	if r == nil || !strings.HasPrefix(r.URL.Path, "/api/v1/") {
		return false
	}
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func duplicateRequestFingerprint(c *gin.Context) (string, error) {
	bodyDigest, err := duplicateRequestBodyDigest(c)
	if err != nil {
		return "", err
	}

	identity := strings.TrimSpace(c.GetHeader("Authorization"))
	if identity == "" {
		identity = "ip:" + c.ClientIP()
	}

	digest := sha256.New()
	writeFingerprintPart(digest, c.Request.Method)
	writeFingerprintPart(digest, c.Request.URL.Path)
	writeFingerprintPart(digest, c.Request.URL.Query().Encode())
	writeFingerprintPart(digest, identity)
	writeFingerprintPart(digest, c.GetHeader("x-tenant-id"))
	writeFingerprintPart(digest, c.GetHeader("X-Client-Source"))
	writeFingerprintPart(digest, bodyDigest)
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func duplicateRequestBodyDigest(c *gin.Context) (string, error) {
	if c.Request.Body == nil {
		return "", nil
	}

	body, err := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("read request body: %w", err)
	}
	if len(body) == 0 {
		return "", nil
	}

	mediaType, params, parseErr := mime.ParseMediaType(c.GetHeader("Content-Type"))
	if parseErr == nil && strings.HasPrefix(mediaType, "multipart/") && params["boundary"] != "" {
		if digest, multipartErr := duplicateMultipartDigest(body, params["boundary"]); multipartErr == nil {
			return digest, nil
		}
	}

	digest := sha256.Sum256(body)
	return hex.EncodeToString(digest[:]), nil
}

func duplicateMultipartDigest(body []byte, boundary string) (string, error) {
	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	digest := sha256.New()
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read multipart part: %w", err)
		}

		content, readErr := io.ReadAll(part)
		_ = part.Close()
		if readErr != nil {
			return "", fmt.Errorf("read multipart content: %w", readErr)
		}

		writeFingerprintPart(digest, part.FormName())
		writeFingerprintPart(digest, part.FileName())
		writeFingerprintPart(digest, part.Header.Get("Content-Type"))
		writeFingerprintPart(digest, string(content))
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func writeFingerprintPart(digest hash.Hash, value string) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = digest.Write(length[:])
	_, _ = digest.Write([]byte(value))
}

type duplicateRequestGuard struct {
	mu               sync.Mutex
	localLocks       map[string]duplicateRequestLocalLock
	lockTTL          time.Duration
	cooldown         time.Duration
	redisWarningOnce sync.Once
}

type duplicateRequestLocalLock struct {
	token     string
	expiresAt time.Time
}

func newDuplicateRequestGuard(lockTTL, cooldown time.Duration) *duplicateRequestGuard {
	return &duplicateRequestGuard{
		localLocks: make(map[string]duplicateRequestLocalLock),
		lockTTL:    lockTTL,
		cooldown:   cooldown,
	}
}

func (g *duplicateRequestGuard) acquire(ctx context.Context, key string) (token string, distributed bool, acquired bool) {
	token = uuid.NewString()
	if !g.acquireLocal(key, token) {
		return "", false, false
	}

	client := redisutil.GetClient()
	if client == nil {
		return token, false, true
	}

	redisCtx, cancel := context.WithTimeout(ctx, duplicateRedisTimeout)
	defer cancel()
	ok, err := client.SetNX(redisCtx, key, token, g.lockTTL).Result()
	if err != nil {
		g.redisWarningOnce.Do(func() {
			logging.LogWarn("Redis 请求防重锁不可用，已降级为单机防重", zap.Error(err))
		})
		return token, false, true
	}
	if !ok {
		g.releaseLocal(key, token)
		return "", false, false
	}
	return token, true, true
}

func (g *duplicateRequestGuard) acquireLocal(key, token string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now()
	if lock, exists := g.localLocks[key]; exists && now.Before(lock.expiresAt) {
		return false
	}
	g.localLocks[key] = duplicateRequestLocalLock{
		token:     token,
		expiresAt: now.Add(g.lockTTL),
	}
	return true
}

func (g *duplicateRequestGuard) releaseAfter(key, token string, distributed bool) {
	release := func() {
		g.releaseLocal(key, token)
		if !distributed {
			return
		}
		client := redisutil.GetClient()
		if client == nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), duplicateRedisTimeout)
		defer cancel()
		_, _ = releaseDuplicateRequestScript.Run(ctx, client, []string{key}, token).Result()
	}
	if g.cooldown <= 0 {
		release()
		return
	}
	time.AfterFunc(g.cooldown, release)
}

func (g *duplicateRequestGuard) releaseLocal(key, token string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.localLocks[key].token == token {
		delete(g.localLocks, key)
	}
}
