package lrucache

import (
	_ "embed"

	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/plugins"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

//go:embed README.md
var docString string

func init() {
	plugins.Register(New)
}

type PluginBundle interface {
	plugins.Wrapper
}

func New(cfg Config) (PluginBundle, error) {
	return &Plugin{
		config: cfg,
	}, nil
}

type Plugin struct {
	config Config
}

func (p Plugin) Doc() string {
	return docString
}

func (p Plugin) Wrap(connector connectors.Connector) (connectors.Connector, error) {
	cache := expirable.NewLRU[string, []map[string]any](p.config.MaxSize, nil, p.config.TTL)

	return &Connector{
		Connector: connector,
		config:    p.config,
		lru:       cache,
	}, nil
}
