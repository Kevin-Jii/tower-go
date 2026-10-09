// Package apidocs converts the swag-generated Swagger 2 document to OpenAPI 3.
//
// swag remains the annotation parser, so existing controller annotations do not
// need to be rewritten. The document exposed to clients is always OpenAPI 3.
package apidocs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/Kevin-Jii/tower-go/docs"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
)

const Version = "3.0.3"

var currentDocument atomic.Pointer[[]byte]

// Document returns the current generated API description as validated OpenAPI 3 JSON.
func Document() ([]byte, error) {
	if current := currentDocument.Load(); current != nil {
		return append([]byte(nil), (*current)...), nil
	}
	return Convert([]byte(docs.SwaggerInfo.ReadDoc()))
}

// ConvertFile converts a file, writes the OpenAPI 3 artifact, and activates it
// for the current process. Activating it ensures startup regeneration is also
// reflected by /swagger/doc.json without recompiling the application.
func ConvertFile(sourcePath, destinationPath string) error {
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("read Swagger 2 document: %w", err)
	}
	document, err := Convert(source)
	if err != nil {
		return err
	}
	if err := os.WriteFile(destinationPath, document, 0o644); err != nil {
		return fmt.Errorf("write OpenAPI 3 document: %w", err)
	}
	active := append([]byte(nil), document...)
	currentDocument.Store(&active)
	return nil
}

func Convert(source []byte) ([]byte, error) {
	var swagger openapi2.T
	if err := json.Unmarshal(source, &swagger); err != nil {
		return nil, fmt.Errorf("decode Swagger 2 document: %w", err)
	}
	if swagger.Swagger != "2.0" {
		return nil, fmt.Errorf("expected Swagger 2.0 source, got %q", swagger.Swagger)
	}

	document, err := openapi2conv.ToV3(&swagger)
	if err != nil {
		return nil, fmt.Errorf("convert Swagger 2 document to OpenAPI 3: %w", err)
	}
	document.OpenAPI = Version
	if err := document.Validate(context.Background()); err != nil {
		return nil, fmt.Errorf("validate OpenAPI 3 document: %w", err)
	}

	result, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode OpenAPI 3 document: %w", err)
	}
	return append(result, '\n'), nil
}
