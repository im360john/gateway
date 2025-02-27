package plugins

import (
	"github.com/centralmind/gateway/connectors"
	"github.com/centralmind/gateway/remapper"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/xerrors"
)

// Config defines the basic configuration interface that all plugin configs must implement
type Config interface {
	// Tag returns the unique identifier for the plugin
	Tag() string
	Doc() string
}

// Plugin is the base interface that all plugins must implement
type Plugin interface {
	// Doc returns the documentation string describing plugin's purpose and configuration
	Doc() string
}

// Interceptor represents a plugin that can process and modify data before it reaches the connector
type Interceptor interface {
	Plugin
	// Process handles the data transformation and returns processed data and a skip flag
	// data: the input data to be processed
	// context: additional context information as key-value pairs
	// Returns: processed data and a boolean indicating if further processing should be skipped
	Process(data map[string]any, context map[string][]string) (procesed map[string]any, skipped bool)
}

// Wrapper represents a plugin that can wrap and enhance a connector's functionality
type Wrapper interface {
	Plugin
	// Wrap takes a connector and returns an enhanced version of it
	// Returns: wrapped connector or error if wrapping fails
	Wrap(connector connectors.Connector) (connectors.Connector, error)
}

// Swaggerer represents a plugin that can modify OpenAPI documentation
type Swaggerer interface {
	Plugin
	// Enrich enhances the OpenAPI documentation with additional specifications
	// Returns: modified OpenAPI documentation
	Enrich(swag *openapi3.T) *openapi3.T
}

var (
	plugins = map[string]func(any) (Plugin, error){}
	configs = map[string]Config{}
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
	configs[t.Tag()] = t
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

// KnownPlugins returns a list of all registered plugin configurations
func KnownPlugins() []Config {
	result := make([]Config, 0, len(plugins))
	for tag := range plugins {
		result = append(result, configs[tag])
	}
	return result
}

// KnownPlugin returns configuration for a specific plugin by tag
func KnownPlugin(tag string) (Config, bool) {
	cfg, ok := configs[tag]
	if !ok {
		return nil, false
	}
	return cfg, true
}
