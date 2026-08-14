package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuditBodyBufferRespectsLimit(t *testing.T) {
	buffer := &auditBodyBuffer{limit: 4}
	if n, err := buffer.Write([]byte("123456")); err != nil || n != 6 {
		t.Fatalf("Write() = (%d, %v), want (6, nil)", n, err)
	}
	if got := buffer.String(); got != "1234" {
		t.Fatalf("buffer = %q, want %q", got, "1234")
	}

	if n, err := buffer.Write([]byte("78")); err != nil || n != 2 {
		t.Fatalf("Write() after limit = (%d, %v), want (2, nil)", n, err)
	}
	if got := buffer.String(); got != "1234" {
		t.Fatalf("buffer grew after limit: %q", got)
	}
}

func TestMultipartPartUploadSkipsAuditLog(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(
		http.MethodPut,
		"/api/v1/galleries/multipart/session-1/parts/1",
		nil,
	)
	if shouldAuditRequest(ctx) {
		t.Fatal("multipart part upload should not create an audit log")
	}

	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/galleries/multipart/init", nil)
	if !shouldAuditRequest(ctx) {
		t.Fatal("multipart initialization should remain auditable")
	}
}
