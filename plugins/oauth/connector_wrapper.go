package oauth

import (
	"context"
	"strings"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/errors"
	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/xcontext"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

type Connector struct {
	connectors.Connector
	config      Config
	oauthConfig *oauth2.Config
}

func (c *Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	// Get token from header
	authHeader := xcontext.Header(ctx, c.config.TokenHeader)
	if authHeader == "" {
		return nil, xerrors.Errorf("empty authorization header: %w", errors.ErrNotAuthorized)
	}

	// Extract token from "Bearer <token>" header
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, xerrors.Errorf("invalid authorization header format: %w", errors.ErrNotAuthorized)
	}
	tokenString := parts[1]

	// Create token
	token := &oauth2.Token{
		AccessToken: tokenString,
	}

	// Validate token through provider
	client := c.oauthConfig.Client(ctx, token)

	// Here we can add scope validation for the method
	if scopes, ok := c.config.MethodScopes[endpoint.MCPMethod]; ok {
		// Validate scopes through provider's API
		// This depends on specific provider
		// For example, for GitHub:
		// resp, err := client.Get("https://api.github.com/user")
		// ...
		_ = client // Use client for scope validation

		if len(scopes) > 0 {
			// Scope validation
			// TODO: implement validation for specific provider
		}
	}

	return c.Connector.Query(ctx, endpoint, params)
}
