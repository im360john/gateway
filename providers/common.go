package providers

import (
	"errors"
)

var (
	ErrUnknownProvider = errors.New("unknown provider")
)

func NewModelProvider(providerConfig ModelProviderConfig) (ModelProvider, error) {
	if providerConfig.Name == "bedrock" {
		return NewBedrockProvider(providerConfig)
	} else if providerConfig.Name == "openai" {
		return NewOpenAIProvider(providerConfig)
	} else if providerConfig.Name == "anthropic" {
		return NewAnthropicProvider(providerConfig, false)
	} else if providerConfig.Name == "anthropic-vertexai" {
		return NewAnthropicProvider(providerConfig, true)
	}

	return nil, ErrUnknownProvider
}
