package api_keys

import (
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/plugins"
	"github.com/getkin/kin-openapi/openapi3"
)

func init() {
	plugins.Register(New)
}

type PluginBundle interface {
	plugins.Wrapper
	plugins.Swaggerer
}

func New(cfg Config) (PluginBundle, error) {
	return &Plugin{
		config: cfg,
	}, nil
}

type Plugin struct {
	config Config
}

func MaybeSetSecurity(op *openapi3.Operation, name string) *openapi3.Operation {
	if op == nil {
		return op
	}
	op.Security = &openapi3.SecurityRequirements{
		openapi3.SecurityRequirement{
			name: []string{}, // Reference the security scheme
		},
	}
	return op
}

func (p Plugin) Enrich(swag *openapi3.T) *openapi3.T {
	securityName := "BearerAuth"
	if swag.Components.SecuritySchemes == nil {
		swag.Components.SecuritySchemes = map[string]*openapi3.SecuritySchemeRef{}
	}
	swag.Components.SecuritySchemes[securityName] = &openapi3.SecuritySchemeRef{
		Value: &openapi3.SecurityScheme{
			Type: "apiKey",
			In:   p.config.Location,
			Name: p.config.Name,
		},
	}
	for key, v := range swag.Paths.Map() {
		v.Get = MaybeSetSecurity(v.Get, securityName)
		v.Delete = MaybeSetSecurity(v.Delete, securityName)
		v.Post = MaybeSetSecurity(v.Post, securityName)
		v.Put = MaybeSetSecurity(v.Put, securityName)
		v.Patch = MaybeSetSecurity(v.Patch, securityName)
		swag.Paths.Set(key, v)
	}
	return swag
}

func (p Plugin) Wrap(connector connectors.Connector) (connectors.Connector, error) {
	return &Connector{
		Connector: connector,
		config:    p.config,
	}, nil
}

func (p Plugin) Doc() string {
	return `
Add auth-check for api-keys, api key is located in headers
`
}
