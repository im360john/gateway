package plugins

import (
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/remapper"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/xerrors"
)

type Config interface {
	Tag() string
}

type Plugin interface {
	Doc() string
}

type Interceptor interface {
	Plugin
	Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool)
}

type Wrapper interface {
	Plugin
	Wrap(connector connectors.Connector) (connectors.Connector, error)
}

type Swaggerer interface {
	Plugin
	Enrich(swag *openapi3.T) *openapi3.T
}

var (
	plugins = map[string]func(any) (Plugin, error){}
)

func Register[TConfig Config, TPlugin Plugin](f func(cfg TConfig) (TPlugin, error)) {
	var t TConfig
	plugins[t.Tag()] = func(a any) (Plugin, error) {
		cfg, err := remapper.Remap[TConfig](a)
		if err != nil {
			return nil, xerrors.Errorf("unable to rempa: %w", err)
		}
		return f(cfg)
	}
}

func New(tag string, config any) (Plugin, error) {
	f, ok := plugins[tag]
	if !ok {
		return nil, xerrors.Errorf("plugin: %s not found", tag)
	}
	return f(config)
}

func Enrich(pluginsCfg map[string]any, schema *openapi3.T) (*openapi3.T, error) {
	for k, v := range pluginsCfg {
		if _, ok := plugins[k]; !ok {
			continue
		}
		plugin, err := plugins[k](v)
		if err != nil {
			return nil, xerrors.Errorf("unable to construct: %s: %w", k, err)
		}
		wrapper, ok := plugin.(Swaggerer)
		if !ok {
			continue
		}
		schema = wrapper.Enrich(schema)
	}
	return schema, nil
}

func Wrap(pluginsCfg map[string]any, connector connectors.Connector) (connectors.Connector, error) {
	for k, v := range pluginsCfg {
		if _, ok := plugins[k]; !ok {
			continue
		}
		plugin, err := plugins[k](v)
		if err != nil {
			return nil, xerrors.Errorf("unable to construct: %s: %w", k, err)
		}
		wrapper, ok := plugin.(Wrapper)
		if !ok {
			continue
		}
		connector, err = wrapper.Wrap(connector)
		if err != nil {
			return nil, xerrors.Errorf("unable to wrap: %s: %w", k, err)
		}
	}
	return connector, nil
}
