// Copyright 2022, Cloudy Sky Software LLC.

package main

import (
	_ "embed"

	"github.com/cloudy-sky-software/pulumi-xyz/provider/pkg/provider"
	"github.com/cloudy-sky-software/pulumi-xyz/provider/pkg/version"
)

var providerName = "xyz"

//go:embed schema.json
var pulumiSchema []byte

//go:embed openapi_generated.yml
var openapiDocBytes []byte

//go:embed metadata.json
var metadataBytes []byte

func main() {
	provider.Serve(providerName, version.Version, pulumiSchema, openapiDocBytes, metadataBytes)
}
