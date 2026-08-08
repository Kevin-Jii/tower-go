package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInternalServiceAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		expectedToken string
		header        string
		wantStatus    int
		wantHandled   bool
	}{
		{name: "valid bearer token", expectedToken: "service-secret", header: "Bearer service-secret", wantStatus: http.StatusNoContent, wantHandled: true},
		{name: "bearer scheme is case insensitive", expectedToken: "service-secret", header: "bearer service-secret", wantStatus: http.StatusNoContent, wantHandled: true},
		{name: "missing header", expectedToken: "service-secret", wantStatus: http.StatusUnauthorized},
		{name: "wrong token", expectedToken: "service-secret", header: "Bearer user-jwt", wantStatus: http.StatusUnauthorized},
		{name: "malformed header", expectedToken: "service-secret", header: "Bearer token with spaces", wantStatus: http.StatusUnauthorized},
		{name: "missing server configuration", header: "Bearer service-secret", wantStatus: http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handled := false
			router := gin.New()
			router.Use(InternalServiceAuthMiddleware(tt.expectedToken))
			router.GET("/internal", func(ctx *gin.Context) {
				handled = true
				ctx.Status(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/internal", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if handled != tt.wantHandled {
				t.Fatalf("handler called = %v, want %v", handled, tt.wantHandled)
			}
			if got := response.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("Cache-Control = %q, want no-store", got)
			}
		})
	}
}
