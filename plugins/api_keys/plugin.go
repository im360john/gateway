package api_keys

import (
	_ "embed"
	"github.com/danielgtaylor/huma/v2"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/plugins"
)

//go:embed README.md
var docString string

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

func MaybeSetSecurity(op *huma.Operation, name string) *huma.Operation {
	if op == nil {
		return op
	}
	op.Security = []map[string][]string{
		{
			name: []string{},
		},
	}
	return op
}

func (p Plugin) Enrich(swag *huma.OpenAPI) *huma.OpenAPI {
	securityName := "BearerAuth"
	if swag.Components.SecuritySchemes == nil {
		swag.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	swag.Components.SecuritySchemes[securityName] = &huma.SecurityScheme{
		Type: "apiKey",
		In:   p.config.Location,
		Name: p.config.Name,
	}
	for _, v := range swag.Paths {
		v.Get = MaybeSetSecurity(v.Get, securityName)
		v.Delete = MaybeSetSecurity(v.Delete, securityName)
		v.Post = MaybeSetSecurity(v.Post, securityName)
		v.Put = MaybeSetSecurity(v.Put, securityName)
		v.Patch = MaybeSetSecurity(v.Patch, securityName)
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
	return docString
}
