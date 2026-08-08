package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"strings"

	httpx "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

// InternalServiceAuthMiddleware authenticates calls made by trusted backend services.
// It deliberately does not accept user JWTs or store/user identity headers.
func InternalServiceAuthMiddleware(expectedToken string) gin.HandlerFunc {
	expectedToken = strings.TrimSpace(expectedToken)

	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-store")

		if expectedToken == "" {
			abortInternalAuth(ctx, http.StatusServiceUnavailable, "internal service authentication is not configured")
			return
		}

		providedToken, ok := internalBearerToken(ctx.GetHeader("Authorization"))
		if !ok || !secureTokenEqual(providedToken, expectedToken) {
			ctx.Header("WWW-Authenticate", "Bearer")
			abortInternalAuth(ctx, http.StatusUnauthorized, "unauthorized")
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

func abortInternalAuth(ctx *gin.Context, status int, message string) {
	ctx.AbortWithStatusJSON(status, httpx.Response{
		Code:    status,
		Message: message,
	})
}
