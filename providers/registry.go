package providers

import (
	"errors"
)

type ModelProviderFactory func(ModelProviderConfig) (ModelProvider, error)

var (
	ErrUnknownProvider = errors.New("unknown ai provider")
	providerRegistry   = make(map[string]ModelProviderFactory)
)

func RegisterModelProvider(name string, factory ModelProviderFactory) {
	providerRegistry[name] = factory
}

func NewModelProvider(config ModelProviderConfig) (ModelProvider, error) {
	if factory, ok := providerRegistry[config.Name]; ok {
		return factory(config)
	}

	return nil, ErrUnknownProvider
}
