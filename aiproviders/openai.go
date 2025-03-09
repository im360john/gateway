package aiproviders

import (
	"context"
	"errors"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	defaultOpenAIModelId          = "o3-mini"
	defaultOpenAITemperature      = float32(0)
	defaultOpenAIMaxTokens        = 100_000
	defaultOpenAIStreamBufferSize = 100
)

var (
	ErrNoAPIKey      = errors.New("OpenAI API key not provided")
	ErrClientNotInit = errors.New("OpenAI client is not initialized")
	ErrEmptyChoices  = errors.New("unexpected empty response from OpenAI")
)

type OpenAIProvider struct {
	Client   *openai.Client
	Endpoint string
}

var _ ModelProvider = (*OpenAIProvider)(nil)

func (op *OpenAIProvider) GetName() string {
	return "OpenAI"
}

func NewOpenAIProvider(providerConfig ModelProviderConfig) (*OpenAIProvider, error) {
	effectiveAPIKey := providerConfig.APIKey
	if effectiveAPIKey == "" {
		effectiveAPIKey = os.Getenv("OPENAI_API_KEY")
		if effectiveAPIKey == "" {
			return nil, ErrNoAPIKey
		}
	}

	effectiveEndpoint := providerConfig.Endpoint
	envEndpoint := os.Getenv("OPENAI_ENDPOINT")
	if effectiveEndpoint == "" && envEndpoint != "" {
		effectiveEndpoint = envEndpoint
	}

	config := openai.DefaultConfig(effectiveAPIKey)
	if effectiveEndpoint != "" {
		config.BaseURL = effectiveEndpoint
	}
	client := openai.NewClientWithConfig(config)

	return &OpenAIProvider{
		Client:   client,
		Endpoint: effectiveEndpoint,
	}, nil
}

func (op *OpenAIProvider) Chat(ctx context.Context, req *ConversationRequest) (*ConversationResponse, error) {
	if op.Client == nil {
		return nil, ErrClientNotInit
	}

	modelId := req.ModelId
	if modelId == "" {
		if envModelId := os.Getenv("OPENAI_MODEL_ID"); envModelId != "" {
			modelId = envModelId
		} else {
			modelId = defaultOpenAIModelId
		}
	}

	messages := prepareOpenAIMessages(req.Messages, req.System)

	temperature := defaultOpenAITemperature
	if req.Temperature >= 0 {
		temperature = req.Temperature
	}

	maxTokens := defaultOpenAIMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	request := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            messages,
		Temperature:         temperature,
		MaxCompletionTokens: maxTokens,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		},
	}

	if req.Reasoning {
		request.ReasoningEffort = "high"
	}

	resp, err := op.Client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error during conversation: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, ErrEmptyChoices
	}

	var responseContentBlocks []ContentBlock
	responseContentBlocks = append(responseContentBlocks, &ContentBlockText{
		Value: resp.Choices[0].Message.Content,
	})

	stopReason := convertOpenAIStopReason(resp.Choices[0].FinishReason)
	usage := convertOpenAIUsage(resp.Usage)

	return &ConversationResponse{
		ProviderName: "OpenAI",
		Content:      responseContentBlocks,
		StopReason:   stopReason,
		ModelId:      modelId,
		Usage:        usage,
	}, nil
}

type OpenAIStreamOutput struct {
	stream *OpenAIStream
}

func (o *OpenAIStreamOutput) GetStream() ChatStream {
	return o.stream
}

type OpenAIStream struct {
	eventCh chan StreamChunk
}

func (s *OpenAIStream) Events() <-chan StreamChunk {
	return s.eventCh
}

func (op *OpenAIProvider) ChatStream(ctx context.Context, req *ConversationRequest) (ChatStreamOutput, error) {
	if op.Client == nil {
		return nil, ErrClientNotInit
	}

	modelId := req.ModelId
	if modelId == "" {
		if envModelId := os.Getenv("OPENAI_MODEL_ID"); envModelId != "" {
			modelId = envModelId
		} else {
			modelId = defaultOpenAIModelId
		}
	}

	messages := prepareOpenAIMessages(req.Messages, req.System)

	temperature := defaultOpenAITemperature
	if req.Temperature >= 0 {
		temperature = req.Temperature
	}

	maxTokens := defaultOpenAIMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	request := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            messages,
		Temperature:         temperature,
		MaxCompletionTokens: maxTokens,
		Stream:              true,
		StreamOptions: &openai.StreamOptions{
			IncludeUsage: true,
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		},
	}

	if req.Reasoning {
		request.ReasoningEffort = "high"
	}

	stream, err := op.Client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error creating stream: %w", err)
	}

	eventCh := make(chan StreamChunk, defaultOpenAIStreamBufferSize)
	openaiStream := &OpenAIStream{
		eventCh: eventCh,
	}

	go func() {
		defer close(eventCh)
		defer stream.Close()

		var stopReason StopReason = StopReasonStop

		for {
			select {
			case <-ctx.Done():
				eventCh <- &StreamChunkError{
					Error: ctx.Err().Error(),
				}
				return
			default:
				response, err := stream.Recv()

				if err != nil {
					eventCh <- &StreamChunkError{
						Error: err.Error(),
					}
					return
				}

				if response.Usage != nil {
					eventCh <- &StreamChunkUsage{
						Usage: convertOpenAIUsage(*response.Usage),
					}
				}

				if len(response.Choices) > 0 {
					choice := response.Choices[0]

					if choice.Delta.Content != "" {
						eventCh <- &StreamChunkContent{
							Content: &ContentBlockText{
								Value: choice.Delta.Content,
							},
						}
					}

					if choice.FinishReason != "" {
						stopReason = convertOpenAIStopReason(choice.FinishReason)
						eventCh <- &StreamChunkStop{
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

func prepareOpenAIMessages(messages []Message, systemContent string) []openai.ChatCompletionMessage {
	var openaiMessages []openai.ChatCompletionMessage

	if systemContent != "" {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemContent,
		})
	}

	for _, msg := range messages {
		role := openai.ChatMessageRoleUser
		if msg.Role == AssistantRole {
			role = openai.ChatMessageRoleAssistant
		}

		var contentText string
		for _, content := range msg.Content {
			if textBlock, ok := content.(*ContentBlockText); ok {
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

func convertOpenAIStopReason(reason openai.FinishReason) StopReason {
	switch reason {
	case openai.FinishReasonStop:
		return StopReasonStop
	case openai.FinishReasonToolCalls:
		return StopReasonToolCalls
	case openai.FinishReasonLength:
		return StopReasonLength
	default:
		return StopReasonStop
	}
}

func convertOpenAIUsage(usage openai.Usage) *ModelUsage {
	return &ModelUsage{
		InputTokens:  usage.PromptTokens,
		OutputTokens: usage.CompletionTokens,
		TotalTokens:  usage.TotalTokens,
	}
}
