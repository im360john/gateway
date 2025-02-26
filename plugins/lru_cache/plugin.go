package lrucache

import (
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/plugins"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

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
	return `
LRU-based Cache

Suitable for scenarios where the cache holds a fixed number of entries and evicts the least recently used items when full.

## Example YAML configuration:

lru_cache:
  max_size: 1000
  ttl: "5m"  # 5 minutes
`
}

func (p Plugin) Wrap(connector connectors.Connector) (connectors.Connector, error) {
	cache := expirable.NewLRU[string, []map[string]any](p.config.MaxSize, nil, p.config.TTL)

	return &Connector{
		Connector: connector,
		config:    p.config,
		lru:       cache,
	}, nil
}
