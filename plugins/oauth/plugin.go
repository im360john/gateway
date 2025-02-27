package oauth

import (
	_ "embed"
	"github.com/danielgtaylor/huma/v2"
	"net/http"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/plugins"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

//go:embed README.md
var docString string

func init() {
	plugins.Register(New)
}

type PluginBundle interface {
	plugins.Wrapper
	plugins.Swaggerer
	plugins.HTTPServer
}

func New(cfg Config) (PluginBundle, error) {
	cfg.WithDefaults()
	oauthConfig := cfg.GetOAuthConfig()
	if oauthConfig == nil {
		return nil, xerrors.New("failed to create OAuth config")
	}

	plugin := &Plugin{
		config:      cfg,
		oauthConfig: oauthConfig,
	}

	return plugin, nil
}

type Plugin struct {
	config      Config
	oauthConfig *oauth2.Config
}

func (p *Plugin) RegisterRoutes(mux *http.ServeMux) {
	if p.config.AuthURL == "" || p.config.CallbackURL == "" {
		return
	}
	// Register HTTP handlers
	mux.HandleFunc(p.config.AuthURL, p.HandleAuthorize)
	mux.HandleFunc(p.config.CallbackURL, p.HandleCallback)
}

func (p *Plugin) Doc() string {
	return docString
}

func (p *Plugin) Wrap(connector connectors.Connector) (connectors.Connector, error) {
	return &Connector{
		Connector:   connector,
		config:      p.config,
		oauthConfig: p.oauthConfig,
	}, nil
}

func (p *Plugin) Enrich(swag *huma.OpenAPI) *huma.OpenAPI {
	// Add OAuth2 security definition
	if swag.Components.SecuritySchemes == nil {
		swag.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}

	swag.Components.SecuritySchemes["OAuth2"] = &huma.SecurityScheme{
		Type:        "oauth2",
		Description: "OAuth2 authentication",
		Flows: &huma.OAuthFlows{
			AuthorizationCode: &huma.OAuthFlow{
				AuthorizationURL: p.oauthConfig.Endpoint.AuthURL,
				TokenURL:         p.oauthConfig.Endpoint.TokenURL,
				Scopes:           make(map[string]string),
			},
		},
	}

	// Add security requirements to all paths
	for _, pathItem := range swag.Paths {
		for _, op := range []*huma.Operation{
			pathItem.Get,
			pathItem.Post,
			pathItem.Put,
			pathItem.Delete,
			pathItem.Patch,
		} {
			if op != nil {
				op.Security = []map[string][]string{
					{
						"OAuth2": []string{},
					},
				}
			}
		}
	}

	return swag
}
