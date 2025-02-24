package api_keys

import (
	"context"
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/errors"
	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/xcontext"
	"golang.org/x/xerrors"
)

type Connector struct {
	connectors.Connector

	config Config
}

func (c Connector) Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error) {
	authToken := xcontext.Header(ctx, c.config.Name)
	if authToken == "" {
		return nil, xerrors.Errorf("empty token: %w", errors.ErrNotAuthorized)
	}
	found := false
	for _, token := range c.config.Keys {
		if token.Key == authToken {
			found = true
			if !token.Allowed(endpoint.MCPMethod) {
				return nil, xerrors.Errorf("method: %s is not authorized for this token: %w", endpoint.MCPMethod, errors.ErrNotAuthorized)
			}
			break
		}
	}
	if !found {
		return nil, xerrors.Errorf("unknown token: %w", errors.ErrNotAuthorized)
	}
	return c.Connector.Query(ctx, endpoint, params)
}
