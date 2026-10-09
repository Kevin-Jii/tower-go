package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSwaggerDocEndpointServesOpenAPI3(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerDocumentationRoutes(r)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	r.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json; charset=utf-8", response.Header().Get("Content-Type"))

	var document map[string]any
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &document))
	require.Equal(t, "3.0.3", document["openapi"])
	require.NotContains(t, document, "swagger")
}

func TestSwaggerUIStillAvailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerDocumentationRoutes(r)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	r.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Body.String(), "Swagger UI")
}
