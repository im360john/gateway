package openai

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/centralmind/gateway/providers"
	openai "github.com/sashabaranov/go-openai"
)

const (
	defaultOpenAIModelId          = "o3-mini"
	defaultGeminiModelId          = "gemini-2.0-flash-thinking-exp-01-21"
	defaultOpenAIMaxTokens        = 100_000
	defaultOpenAIStreamBufferSize = 100
)

var openAIModelsWithoutTemperatureSupport = []string{
	"o3",
	"o1",
}

var (
	ErrNoAPIKey      = errors.New("OpenAI API key not provided")
	ErrClientNotInit = errors.New("OpenAI client is not initialized")
	ErrEmptyChoices  = errors.New("unexpected empty response from OpenAI")
)

type OpenAIProvider struct {
	Client   *openai.Client
	Endpoint string
	Gemini   bool
}

var _ providers.ModelProvider = (*OpenAIProvider)(nil)

func init() {
	providers.RegisterModelProvider("openai", NewOpenAIProvider)
	providers.RegisterModelProvider("gemini", NewGeminiProvider)
}

func (op *OpenAIProvider) GetName() string {
	if op.Gemini {
		return "Gemini"
	}

	return "OpenAI"
}

func (ap *OpenAIProvider) CostEstimate(modelId string, usage providers.ModelUsage) float64 {
	var inputPrice, outputPrice float64
	const oneMillion = 1_000_000.0

	switch {
	case strings.HasPrefix(modelId, "gpt-4o-mini"):
		inputPrice = 0.15 / oneMillion
		outputPrice = 0.60 / oneMillion
	case strings.HasPrefix(modelId, "gpt-4o"):
		inputPrice = 2.5 / oneMillion
		outputPrice = 10.0 / oneMillion
	case strings.HasPrefix(modelId, "o3-mini"):
		inputPrice = 1.1 / oneMillion
		outputPrice = 4.4 / oneMillion
	case strings.HasPrefix(modelId, "o1"):
		inputPrice = 15 / oneMillion
		outputPrice = 60 / oneMillion
	case strings.HasPrefix(modelId, "gemini"):
		inputPrice = 0.10 / oneMillion
		outputPrice = 0.40 / oneMillion
	default:
		return 0.0
	}

	inputCost := float64(usage.InputTokens) * inputPrice
	outputCost := float64(usage.OutputTokens) * outputPrice
	totalCost := inputCost + outputCost

	return totalCost
}

func NewOpenAIProvider(providerConfig providers.ModelProviderConfig) (providers.ModelProvider, error) {
	return NewOpenAIProviderIntl(providerConfig, false)
}

func NewGeminiProvider(providerConfig providers.ModelProviderConfig) (providers.ModelProvider, error) {
	return NewOpenAIProviderIntl(providerConfig, true)
}

func NewOpenAIProviderIntl(providerConfig providers.ModelProviderConfig, gemini bool) (providers.ModelProvider, error) {
	effectiveAPIKey := providerConfig.APIKey
	if effectiveAPIKey == "" {
		if gemini {
			effectiveAPIKey = os.Getenv("GEMINI_API_KEY")
		} else {
			effectiveAPIKey = os.Getenv("OPENAI_API_KEY")
		}

		if effectiveAPIKey == "" {
			return nil, ErrNoAPIKey
		}
	}

	effectiveEndpoint := providerConfig.Endpoint
	if gemini {
		effectiveEndpoint = "https://generativelanguage.googleapis.com/v1beta/openai/"
	} else {
		envEndpoint := os.Getenv("OPENAI_ENDPOINT")
		if effectiveEndpoint == "" && envEndpoint != "" {
			effectiveEndpoint = envEndpoint
		}
	}

	config := openai.DefaultConfig(effectiveAPIKey)
	if effectiveEndpoint != "" {
		config.BaseURL = effectiveEndpoint
	}
	client := openai.NewClientWithConfig(config)

	return &OpenAIProvider{
		Client:   client,
		Endpoint: effectiveEndpoint,
		Gemini:   gemini,
	}, nil
}

func (op *OpenAIProvider) Chat(ctx context.Context, req *providers.ConversationRequest) (*providers.ConversationResponse, error) {
	if op.Client == nil {
		return nil, ErrClientNotInit
	}

	modelId := req.ModelId
	if modelId == "" {
		var envModelId string
		if op.Gemini {
			envModelId = os.Getenv("GEMINI_MODEL_ID")
		} else {
			envModelId = os.Getenv("OPENAI_MODEL_ID")
		}

		if envModelId != "" {
			modelId = envModelId
		} else {
			if op.Gemini {
				modelId = defaultGeminiModelId
			} else {
				modelId = defaultOpenAIModelId
			}
		}
	}

	messages := prepareOpenAIMessages(req.Messages, req.System)

	maxTokens := defaultOpenAIMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	request := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            messages,
		MaxCompletionTokens: maxTokens,
	}

	supportsTemperature := true
	for _, model := range openAIModelsWithoutTemperatureSupport {
		if strings.HasPrefix(modelId, model) {
			supportsTemperature = false
			break
		}
	}

	if supportsTemperature {
		request.Temperature = max(req.Temperature, 0.0)
	}

	if req.Reasoning && !op.Gemini {
		request.ReasoningEffort = "high"
	}

	if req.JsonResponse && op.Endpoint == "" {
		request.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		}
	}

	resp, err := op.Client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, ErrEmptyChoices
	}

	var responseContentBlocks []providers.ContentBlock
	if req.JsonResponse {
		responseContentBlocks = append(responseContentBlocks, &providers.ContentBlockText{
			Value: providers.ExtractJSON(resp.Choices[0].Message.Content),
		})
	} else {
		responseContentBlocks = append(responseContentBlocks, &providers.ContentBlockText{
			Value: resp.Choices[0].Message.Content,
		})
	}

	stopReason := convertOpenAIStopReason(resp.Choices[0].FinishReason)
	usage := convertOpenAIUsage(resp.Usage)

	return &providers.ConversationResponse{
		ProviderName: "OpenAI",
		ModelId:      modelId,
		Content:      responseContentBlocks,
		StopReason:   stopReason,
		Usage:        usage,
	}, nil
}

type OpenAIStreamOutput struct {
	stream *OpenAIStream
}

func (o *OpenAIStreamOutput) GetStream() providers.ChatStream {
	return o.stream
}

type OpenAIStream struct {
	eventCh chan providers.StreamChunk
}

func (s *OpenAIStream) Events() <-chan providers.StreamChunk {
	return s.eventCh
}

func (op *OpenAIProvider) ChatStream(ctx context.Context, req *providers.ConversationRequest) (providers.ChatStreamOutput, error) {
	if op.Client == nil {
		return nil, ErrClientNotInit
	}

	modelId := req.ModelId
	if modelId == "" {
		var envModelId string
		if op.Gemini {
			envModelId = os.Getenv("GEMINI_MODEL_ID")
		} else {
			envModelId = os.Getenv("OPENAI_MODEL_ID")
		}

		if envModelId != "" {
			modelId = envModelId
		} else {
			if op.Gemini {
				modelId = defaultGeminiModelId
			} else {
				modelId = defaultOpenAIModelId
			}
		}
	}

	messages := prepareOpenAIMessages(req.Messages, req.System)

	maxTokens := defaultOpenAIMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	request := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            messages,
		MaxCompletionTokens: maxTokens,
		Stream:              true,
		StreamOptions: &openai.StreamOptions{
			IncludeUsage: true,
		},
	}

	supportsTemperature := true
	for _, model := range openAIModelsWithoutTemperatureSupport {
		if strings.HasPrefix(modelId, model) {
			supportsTemperature = false
			break
		}
	}

	if supportsTemperature {
		request.Temperature = max(req.Temperature, 0.0)
	}

	if req.Reasoning && !op.Gemini {
		request.ReasoningEffort = "high"
	}

	if req.JsonResponse && op.Endpoint == "" {
		request.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		}
	}

	stream, err := op.Client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return nil, err
	}

	eventCh := make(chan providers.StreamChunk, defaultOpenAIStreamBufferSize)
	openaiStream := &OpenAIStream{
		eventCh: eventCh,
	}

	go func() {
		defer close(eventCh)
		defer stream.Close()

		var stopReason providers.StopReason = providers.StopReasonStop

		for {
			select {
			case <-ctx.Done():
				eventCh <- &providers.StreamChunkError{
					Error: ctx.Err().Error(),
				}
				return
			default:
				response, err := stream.Recv()

				if err != nil {
					eventCh <- &providers.StreamChunkError{
						Error: err.Error(),
					}
					return
				}

				if response.Usage != nil {
					eventCh <- &providers.StreamChunkUsage{
						ModelId: modelId,
						Usage:   convertOpenAIUsage(*response.Usage),
					}
				}

				if len(response.Choices) > 0 {
					choice := response.Choices[0]

					if choice.Delta.Content != "" {
						eventCh <- &providers.StreamChunkContent{
							Content: &providers.ContentBlockText{
								Value: choice.Delta.Content,
							},
						}
					}

					if choice.FinishReason != "" {
						stopReason = convertOpenAIStopReason(choice.FinishReason)
						eventCh <- &providers.StreamChunkStop{
							StopReason: stopReason,
						}
					}
				}
			}
		}
	}()

	return &OpenAIStreamOutput{
		stream: openaiStream,
	}, nil
}

func prepareOpenAIMessages(messages []providers.Message, systemContent string) []openai.ChatCompletionMessage {
	var openaiMessages []openai.ChatCompletionMessage

	if systemContent != "" {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemContent,
		})
	}

	for _, msg := range messages {
		role := openai.ChatMessageRoleUser
		if msg.Role == providers.AssistantRole {
			role = openai.ChatMessageRoleAssistant
		}

		var contentText string
		for _, content := range msg.Content {
			if textBlock, ok := content.(*providers.ContentBlockText); ok {
				contentText += textBlock.Value
			}
		}

		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: contentText,
		})
	}

	return openaiMessages
}

func convertOpenAIStopReason(reason openai.FinishReason) providers.StopReason {
	switch reason {
	case openai.FinishReasonStop:
		return providers.StopReasonStop
	case openai.FinishReasonToolCalls:
		return providers.StopReasonToolCalls
	case openai.FinishReasonLength:
		return providers.StopReasonLength
	default:
		return providers.StopReasonStop
	}
}

func convertOpenAIUsage(usage openai.Usage) *providers.ModelUsage {
	return &providers.ModelUsage{
		InputTokens:  usage.PromptTokens,
		OutputTokens: usage.CompletionTokens,
		TotalTokens:  usage.TotalTokens,
	}
}
