package gen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudy-sky-software/pulumi-provider-framework/openapi"
)

func TestPulumiSchema(t *testing.T) {
	b, err := os.ReadFile(filepath.Join("..", "..", "cmd", "pulumi-gen-xyz", "openapi.yml"))
	if err != nil {
		t.Fatalf("Failed reading openapi.yml: %v", err)
	}

	oaSpec := openapi.GetOpenAPISpec(b)

	PulumiSchema(*oaSpec)
}
