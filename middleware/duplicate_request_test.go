package middleware

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	httpx "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

func TestDuplicateRequestMiddlewareRejectsConcurrentRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	guard := newDuplicateRequestGuard(time.Second, 10*time.Millisecond)
	router := gin.New()
	router.Use(duplicateRequestMiddleware(guard))

	entered := make(chan struct{}, 1)
	release := make(chan struct{})
	var handled atomic.Int32
	router.POST("/api/v1/orders", func(c *gin.Context) {
		handled.Add(1)
		entered <- struct{}{}
		<-release
		httpx.Success(c, gin.H{"created": true})
	})

	firstDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		firstDone <- performDuplicateRequest(router, `{"name":"same"}`, "Bearer user-1")
	}()
	<-entered

	duplicate := performDuplicateRequest(router, `{"name":"same"}`, "Bearer user-1")
	if duplicate.Code != http.StatusOK {
		t.Fatalf("duplicate HTTP status = %d, want %d", duplicate.Code, http.StatusOK)
	}
	var body httpx.Response
	if err := json.Unmarshal(duplicate.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode duplicate response: %v", err)
	}
	if body.Code != apicode.DuplicateOperation.Num {
		t.Fatalf("duplicate code = %d, want %d", body.Code, apicode.DuplicateOperation.Num)
	}
	if body.Message != "相同请求正在处理中，请勿重复操作" {
		t.Fatalf("duplicate message = %q", body.Message)
	}
	if handled.Load() != 1 {
		t.Fatalf("handler called %d times before release, want 1", handled.Load())
	}

	close(release)
	first := <-firstDone
	if first.Code != http.StatusOK {
		t.Fatalf("first HTTP status = %d, want %d", first.Code, http.StatusOK)
	}
	time.Sleep(20 * time.Millisecond)

	afterCooldown := performDuplicateRequest(router, `{"name":"same"}`, "Bearer user-1")
	if afterCooldown.Code != http.StatusOK {
		t.Fatalf("request after cooldown HTTP status = %d, want %d", afterCooldown.Code, http.StatusOK)
	}
	if handled.Load() != 2 {
		t.Fatalf("handler called %d times after cooldown, want 2", handled.Load())
	}
}

func TestDuplicateRequestFingerprintIncludesBodyAndIdentity(t *testing.T) {
	first, restored := fingerprintForTest(t, `{"name":"first"}`, "Bearer user-1")
	same, _ := fingerprintForTest(t, `{"name":"first"}`, "Bearer user-1")
	differentBody, _ := fingerprintForTest(t, `{"name":"second"}`, "Bearer user-1")
	differentUser, _ := fingerprintForTest(t, `{"name":"first"}`, "Bearer user-2")

	if first != same {
		t.Fatal("identical requests produced different fingerprints")
	}
	if first == differentBody {
		t.Fatal("different request bodies produced the same fingerprint")
	}
	if first == differentUser {
		t.Fatal("different users produced the same fingerprint")
	}
	if restored != `{"name":"first"}` {
		t.Fatalf("restored request body = %q", restored)
	}
}

func TestDuplicateRequestFingerprintNormalizesMultipartBoundary(t *testing.T) {
	firstBody, firstType := multipartBodyForTest(t, "customer", "order.txt", "same file")
	secondBody, secondType := multipartBodyForTest(t, "customer", "order.txt", "same file")

	first := multipartFingerprintForTest(t, firstBody, firstType)
	second := multipartFingerprintForTest(t, secondBody, secondType)
	if first != second {
		t.Fatal("equivalent multipart requests produced different fingerprints")
	}
}

func performDuplicateRequest(handler http.Handler, body, authorization string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders?store_id=1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorization)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func fingerprintForTest(t *testing.T, body, authorization string) (string, string) {
	t.Helper()
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/orders?store_id=1", strings.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("Authorization", authorization)
	fingerprint, err := duplicateRequestFingerprint(ctx)
	if err != nil {
		t.Fatalf("fingerprint request: %v", err)
	}
	restored := new(bytes.Buffer)
	_, _ = restored.ReadFrom(ctx.Request.Body)
	return fingerprint, restored.String()
}

func multipartBodyForTest(t *testing.T, fieldValue, filename, fileContent string) ([]byte, string) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("customer", fieldValue); err != nil {
		t.Fatalf("write multipart field: %v", err)
	}
	file, err := writer.CreateFormFile("attachment", filename)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := file.Write([]byte(fileContent)); err != nil {
		t.Fatalf("write multipart file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	return body.Bytes(), writer.FormDataContentType()
}

func multipartFingerprintForTest(t *testing.T, body []byte, contentType string) string {
	t.Helper()
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/files/upload", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", contentType)
	ctx.Request.Header.Set("Authorization", "Bearer user-1")
	fingerprint, err := duplicateRequestFingerprint(ctx)
	if err != nil {
		t.Fatalf("fingerprint multipart request: %v", err)
	}
	return fingerprint
}
