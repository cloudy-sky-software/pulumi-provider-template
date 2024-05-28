package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	fwCallback "github.com/cloudy-sky-software/pulumi-provider-framework/callback"
	fwRest "github.com/cloudy-sky-software/pulumi-provider-framework/rest"
)

type xyzProvider struct {
	name    string
	version string

	apiKey string
}

var (
	handler  *fwRest.Provider
	callback fwCallback.ProviderCallback
)

func makeProvider(host *provider.HostClient, name, version string, pulumiSchemaBytes, openapiDocBytes, metadataBytes []byte) (pulumirpc.ResourceProviderServer, error) {
	p := &xyzProvider{
		name:    name,
		version: version,
	}

	callback = p
	rp, err := fwRest.MakeProvider(host, name, version, pulumiSchemaBytes, openapiDocBytes, metadataBytes, callback)

	handler = rp.(*fwRest.Provider)

	return rp, err
}

func (p *xyzProvider) GetAuthorizationHeader() string {
	return fmt.Sprintf("%s %s", authSchemePrefix, p.apiKey)
}

func (p *xyzProvider) OnPreInvoke(ctx context.Context, req *pulumirpc.InvokeRequest, httpReq *http.Request) error {
	return nil
}

func (p *xyzProvider) OnPostInvoke(ctx context.Context, req *pulumirpc.InvokeRequest, outputs interface{}) (map[string]interface{}, error) {
	return outputs.(map[string]interface{}), nil
}

// OnConfigure is called by the provider framework when Pulumi calls Configure on
// the resource provider server.
func (p *xyzProvider) OnConfigure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	apiKey, ok := req.GetVariables()["xyz:config:apiKey"]
	if !ok {
		// Check if it's set as an env var.
		envVarNames := handler.GetSchemaSpec().Provider.InputProperties["apiKey"].DefaultInfo.Environment
		for _, n := range envVarNames {
			v := os.Getenv(n)
			if v != "" {
				apiKey = v
			}
		}

		// Return an error if the API key is still empty.
		if apiKey == "" {
			return nil, errors.New("api key is required")
		}
	}

	logging.V(3).Info("Configuring Xyz API key")
	p.apiKey = apiKey

	return &pulumirpc.ConfigureResponse{
		AcceptSecrets: true,
	}, nil
}

// OnDiff checks what impacts a hypothetical update will have on the resource's properties.
func (p *xyzProvider) OnDiff(ctx context.Context, req *pulumirpc.DiffRequest, resourceTypeToken string, diff *resource.ObjectDiff, jsonReq *openapi3.MediaType) (*pulumirpc.DiffResponse, error) {
	return nil, nil
}

func (p *xyzProvider) OnPreCreate(ctx context.Context, req *pulumirpc.CreateRequest, httpReq *http.Request) error {
	return nil
}

// OnPostCreate allocates a new instance of the provided resource and returns its unique ID afterwards.
func (p *xyzProvider) OnPostCreate(ctx context.Context, req *pulumirpc.CreateRequest, outputs interface{}) (map[string]interface{}, error) {
	return outputs.(map[string]interface{}), nil
}

func (p *xyzProvider) OnPreRead(ctx context.Context, req *pulumirpc.ReadRequest, httpReq *http.Request) error {
	return nil
}

func (p *xyzProvider) OnPostRead(ctx context.Context, req *pulumirpc.ReadRequest, outputs interface{}) (map[string]interface{}, error) {
	return outputs.(map[string]interface{}), nil
}

func (p *xyzProvider) OnPreUpdate(ctx context.Context, req *pulumirpc.UpdateRequest, httpReq *http.Request) error {
	return nil
}

func (p *xyzProvider) OnPostUpdate(ctx context.Context, req *pulumirpc.UpdateRequest, httpReq http.Request, outputs interface{}) (map[string]interface{}, error) {
	return outputs.(map[string]interface{}), nil
}

func (p *xyzProvider) OnPreDelete(ctx context.Context, req *pulumirpc.DeleteRequest, httpReq *http.Request) error {
	return nil
}

func (p *xyzProvider) OnPostDelete(ctx context.Context, req *pulumirpc.DeleteRequest) error {
	return nil
}
