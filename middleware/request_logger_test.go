package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestLoggerSkipsBinaryUploadBodies(t *testing.T) {
	tests := []struct {
		path        string
		contentType string
		want        bool
	}{
		{path: "/api/v1/orders", contentType: "application/json", want: true},
		{path: "/api/v1/galleries/upload", contentType: "multipart/form-data; boundary=test", want: false},
		{path: "/api/v1/galleries/multipart/session-1/parts/1", contentType: "application/octet-stream", want: false},
	}
	for _, tt := range tests {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request = httptest.NewRequest(http.MethodPost, tt.path, nil)
		ctx.Request.Header.Set("Content-Type", tt.contentType)
		if got := shouldReadRequestLogBody(ctx); got != tt.want {
			t.Fatalf("shouldReadRequestLogBody(%q, %q) = %v, want %v", tt.path, tt.contentType, got, tt.want)
		}
	}
}
