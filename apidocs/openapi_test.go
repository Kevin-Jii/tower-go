package apidocs

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/Kevin-Jii/tower-go/docs"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestDocumentIsOpenAPI3(t *testing.T) {
	data, err := Document()
	require.NoError(t, err)

	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw))
	require.Equal(t, Version, raw["openapi"])
	require.NotContains(t, raw, "swagger")

	document, err := openapi3.NewLoader().LoadFromData(data)
	require.NoError(t, err)
	require.Len(t, document.Servers, 2)
	require.Equal(t, "http://localhost:10024/api/v1", document.Servers[0].URL)

	bearer := document.Components.SecuritySchemes["Bearer"]
	require.NotNil(t, bearer)
	require.Equal(t, "apiKey", bearer.Value.Type)
	require.Equal(t, "Authorization", bearer.Value.Name)
	require.Equal(t, "header", bearer.Value.In)
}

func TestDocumentConvertsBodyParametersToRequestBodies(t *testing.T) {
	data, err := Document()
	require.NoError(t, err)

	document, err := openapi3.NewLoader().LoadFromData(data)
	require.NoError(t, err)

	login := document.Paths.Find("/auth/login")
	require.NotNil(t, login)
	require.NotNil(t, login.Post)
	require.NotNil(t, login.Post.RequestBody)
	require.Contains(t, login.Post.RequestBody.Value.Content, "application/json")

	upload := document.Paths.Find("/files/upload")
	require.NotNil(t, upload)
	require.NotNil(t, upload.Post)
	require.NotNil(t, upload.Post.RequestBody)
	require.Contains(t, upload.Post.RequestBody.Value.Content, "multipart/form-data")
}

func TestDocumentPreservesOperationSecurity(t *testing.T) {
	data, err := Document()
	require.NoError(t, err)

	document, err := openapi3.NewLoader().LoadFromData(data)
	require.NoError(t, err)

	path := document.Paths.Find("/dict-types")
	require.NotNil(t, path)
	operation := path.Get
	require.NotNil(t, operation)
	require.NotNil(t, operation.Security)
	require.Equal(t, "Bearer", firstSecurityScheme(*operation.Security))
}

func TestCheckedInDocumentMatchesGeneratedDocument(t *testing.T) {
	generated, err := Convert([]byte(docs.SwaggerInfo.ReadDoc()))
	require.NoError(t, err)
	checkedIn, err := os.ReadFile("../docs/openapi.json")
	require.NoError(t, err)
	require.True(t, bytes.Equal(generated, checkedIn), "docs/openapi.json is stale; run make docs")
}

func firstSecurityScheme(requirements openapi3.SecurityRequirements) string {
	for _, requirement := range requirements {
		for name := range requirement {
			return name
		}
	}
	return ""
}

func TestConvertFileActivatesGeneratedDocument(t *testing.T) {
	previous := currentDocument.Load()
	t.Cleanup(func() { currentDocument.Store(previous) })

	source := []byte(docs.SwaggerInfo.ReadDoc())
	var raw map[string]any
	require.NoError(t, json.Unmarshal(source, &raw))
	raw["info"].(map[string]any)["version"] = "runtime-refresh-test"
	source, err := json.Marshal(raw)
	require.NoError(t, err)

	sourcePath := t.TempDir() + "/swagger.json"
	destinationPath := t.TempDir() + "/openapi.json"
	require.NoError(t, os.WriteFile(sourcePath, source, 0o644))
	require.NoError(t, ConvertFile(sourcePath, destinationPath))

	active, err := Document()
	require.NoError(t, err)
	require.Contains(t, string(active), `"version": "runtime-refresh-test"`)
}

func TestConvertRejectsNonSwagger2Input(t *testing.T) {
	_, err := Convert([]byte(`{"openapi":"3.0.3"}`))
	require.ErrorContains(t, err, "expected Swagger 2.0 source")
}
