package connectors

import (
	"context"
	"github.com/centralmind/gateway/model"
	"github.com/centralmind/gateway/remapper"
	"github.com/pkg/errors"
)

type Config interface {
	Type() string
}

type Connector interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, endpoint model.Endpoint, params map[string]any) ([]map[string]any, error)
}

var interceptors = map[string]func(any) (Connector, error){}

func Register[TConfig Config](f func(cfg TConfig) (Connector, error)) {
	var t TConfig
	interceptors[t.Type()] = func(a any) (Connector, error) {
		cfg, err := remapper.Remap[TConfig](a)
		if err != nil {
			return nil, errors.Errorf("unable to rempa: %w", err)
		}
		return f(cfg)
	}
}

func New(tag string, config any) (Connector, error) {
	f, ok := interceptors[tag]
	if !ok {
		return nil, errors.Errorf("connector: %s not found", tag)
	}
	return f(config)
}
