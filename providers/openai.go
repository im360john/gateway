package providers

import (
	"context"
	"errors"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	defaultOpenAIModelId          = "o3-mini"
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

func (ap *OpenAIProvider) CostEstimate(modelId string, usage ModelUsage) float64 {
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
	default:
		return 0.0
	}

	inputCost := float64(usage.InputTokens) * inputPrice
	outputCost := float64(usage.OutputTokens) * outputPrice
	totalCost := inputCost + outputCost

	return totalCost
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

	maxTokens := defaultOpenAIMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	request := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            messages,
		Temperature:         req.Temperature,
		MaxCompletionTokens: maxTokens,
	}

	if req.Reasoning {
		request.ReasoningEffort = "high"
	}

	if req.JsonResponse {
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

	var responseContentBlocks []ContentBlock
	if req.JsonResponse {
		responseContentBlocks = append(responseContentBlocks, &ContentBlockText{
			Value: ExtractJSON(resp.Choices[0].Message.Content),
		})
	} else {
		responseContentBlocks = append(responseContentBlocks, &ContentBlockText{
			Value: resp.Choices[0].Message.Content,
		})
	}

	stopReason := convertOpenAIStopReason(resp.Choices[0].FinishReason)
	usage := convertOpenAIUsage(resp.Usage)

	return &ConversationResponse{
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

	maxTokens := defaultOpenAIMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	request := openai.ChatCompletionRequest{
		Model:               modelId,
		Messages:            messages,
		Temperature:         req.Temperature,
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
		return nil, err
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
						ModelId: modelId,
						Usage:   convertOpenAIUsage(*response.Usage),
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
