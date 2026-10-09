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

func TestCustomDocumentationPage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerDocumentationRoutes(r)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/docs", nil)
	r.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Header().Get("Content-Type"), "text/html")
	require.Contains(t, response.Body.String(), "Tower Go API 文档")
	require.Contains(t, response.Body.String(), "src=\"/docs-assets/tower-logo.svg\"")
	require.Contains(t, response.Body.String(), "id=\"navigation\"")
	require.Contains(t, response.Body.String(), "id=\"operation-view\"")
	require.Contains(t, response.Body.String(), "Try it out")
	require.Contains(t, response.Body.String(), "id=\"global-config-button\"")
	require.Contains(t, response.Body.String(), "id=\"global-bearer-token\"")
	require.Contains(t, response.Body.String(), "data-global-header-name")
	require.Contains(t, response.Body.String(), "tower-go-api-docs-global-request-config")
	require.Contains(t, response.Body.String(), "data-parameter-tab=\"header\"")
	require.Contains(t, response.Body.String(), "data-parameter-pane=\"body\"")
	require.Contains(t, response.Body.String(), "fetch('/swagger/doc.json')")
}

func TestCustomDocumentationAssetsAvailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	registerDocumentationRoutes(r)

	for _, path := range []string{
		"/docs-assets/tower-logo.svg",
		"/docs-assets/tower-logo.svg",
	} {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, path, nil)
		r.ServeHTTP(response, request)
		require.Equalf(t, http.StatusOK, response.Code, "asset %s", path)
		require.NotEmptyf(t, response.Body.String(), "asset %s", path)
		if path == "/docs-assets/tower-logo.svg" {
			require.Contains(t, response.Body.String(), "<svg")
		}
	}
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
