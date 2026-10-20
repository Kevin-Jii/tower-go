package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	httpx "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

// InternalServiceAuthMiddleware 验证受信任后端服务发起的调用。
// 此中间件不接受用户 JWT 或门店、用户身份请求头。
func InternalServiceAuthMiddleware(expectedToken string) gin.HandlerFunc {
	expectedToken = strings.TrimSpace(expectedToken)

	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-store")

		if expectedToken == "" {
			abortInternalAuth(ctx, http.StatusServiceUnavailable, apicode.ConfigMissing)
			return
		}

		providedToken, ok := internalBearerToken(ctx.GetHeader("Authorization"))
		if !ok || !secureTokenEqual(providedToken, expectedToken) {
			ctx.Header("WWW-Authenticate", "Bearer")
			abortInternalAuth(ctx, http.StatusUnauthorized, apicode.TokenInvalid)
			return
		}

		ctx.Next()
	}
}

func internalBearerToken(value string) (string, bool) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(value), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	if token == "" || strings.ContainsAny(token, " \t\r\n") {
		return "", false
	}
	return token, true
}

func secureTokenEqual(provided, expected string) bool {
	providedHash := sha256.Sum256([]byte(provided))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(providedHash[:], expectedHash[:]) == 1
}

func abortInternalAuth(ctx *gin.Context, status int, code apicode.Code) {
	ctx.AbortWithStatusJSON(status, httpx.Response{
		Code:    code.Num,
		Message: code.Msg,
	})
}
