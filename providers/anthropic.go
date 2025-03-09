package providers

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/vertex"
)

const (
	defaultVertexAIRegion            = "us-east1"
	defaultAnthropicModelId          = "claude-3-7-sonnet-20250219"
	defaultVertexAIModelId           = "claude-3-7-sonnet@20250219"
	defaultAnthropicMaxTokens        = 64000
	defaultAnthropicStreamBufferSize = 100
)

var (
	ErrNoAnthropicAPIKey      = errors.New("anthropic API key not provided")
	ErrAnthropicClientNotInit = errors.New("anthropic client is not initialized")
	ErrAnthropicEmptyResponse = errors.New("unexpected empty response from Anthropic")
)

type AnthropicProvider struct {
	Client   *anthropic.Client
	Endpoint string
	VertexAI bool
}

var _ ModelProvider = (*AnthropicProvider)(nil)

func (ap *AnthropicProvider) GetName() string {
	if ap.VertexAI {
		return "Anthropic VertexAI"
	}

	return "Anthropic"
}

func (ap *AnthropicProvider) CostEstimate(modelId string, usage ModelUsage) float64 {
	var inputPrice, outputPrice float64
	const oneMillion = 1_000_000.0

	switch {
	case strings.Contains(modelId, "sonnet"):
		inputPrice = 3.75 / oneMillion
		outputPrice = 15.0 / oneMillion
	default:
		return 0.0
	}

	inputCost := float64(usage.InputTokens) * inputPrice
	outputCost := float64(usage.OutputTokens) * outputPrice
	totalCost := inputCost + outputCost

	return totalCost
}

func NewAnthropicProvider(providerConfig ModelProviderConfig, vertexAI bool) (*AnthropicProvider, error) {
	var client *anthropic.Client

	effectiveEndpoint := providerConfig.Endpoint
	envEndpoint := os.Getenv("ANTHROPIC_ENDPOINT")
	if effectiveEndpoint == "" && envEndpoint != "" {
		effectiveEndpoint = envEndpoint
	}

	if !vertexAI {
		effectiveAPIKey := providerConfig.APIKey
		if effectiveAPIKey == "" {
			effectiveAPIKey = os.Getenv("ANTHROPIC_API_KEY")
			if effectiveAPIKey == "" {
				return nil, ErrNoAnthropicAPIKey
			}
		}

		clientOpts := []option.RequestOption{
			option.WithAPIKey(effectiveAPIKey),
		}

		if effectiveEndpoint != "" {
			clientOpts = append(clientOpts, option.WithBaseURL(effectiveEndpoint))
		}

		client = anthropic.NewClient(clientOpts...)
	} else {
		effectiveVertexAIRegion := providerConfig.VertexAIRegion
		if effectiveVertexAIRegion == "" {
			effectiveVertexAIRegion = os.Getenv("ANTHROPIC_VERTEXAI_REGION")

			if effectiveVertexAIRegion == "" {
				effectiveVertexAIRegion = defaultVertexAIRegion
			}
		}

		effectiveVertexAIProject := providerConfig.VertexAIProject
		if effectiveVertexAIProject == "" {
			effectiveVertexAIProject = os.Getenv("ANTHROPIC_VERTEXAI_PROJECT")
		}

		client = anthropic.NewClient(
			vertex.WithGoogleAuth(context.Background(), effectiveVertexAIRegion, effectiveVertexAIProject),
		)
	}

	return &AnthropicProvider{
		Client:   client,
		Endpoint: effectiveEndpoint,
		VertexAI: vertexAI,
	}, nil
}

func (ap *AnthropicProvider) Chat(ctx context.Context, req *ConversationRequest) (*ConversationResponse, error) {
	if ap.Client == nil {
		return nil, ErrAnthropicClientNotInit
	}

	modelId := req.ModelId
	if modelId == "" {
		if envModelId := os.Getenv("ANTHROPIC_MODEL_ID"); envModelId != "" {
			modelId = envModelId
		} else {
			if ap.VertexAI {
				modelId = defaultVertexAIModelId
			} else {
				modelId = defaultAnthropicModelId
			}
		}
	}

	temperature := req.Temperature
	if req.Reasoning {
		temperature = 1.0
	}

	maxTokens := defaultAnthropicMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	messages := prepareAnthropicMessages(req.Messages)

	params := anthropic.MessageNewParams{
		Model:       anthropic.F(modelId),
		MaxTokens:   anthropic.F(int64(maxTokens)),
		Temperature: anthropic.F(float64(temperature)),
		Messages:    anthropic.F(messages),
	}

	if req.Reasoning {
		params.Thinking = anthropic.F[anthropic.ThinkingConfigParamUnion](anthropic.ThinkingConfigEnabledParam{
			BudgetTokens: anthropic.F(int64(4096)),
			Type:         anthropic.F(anthropic.ThinkingConfigEnabledTypeEnabled),
		})
	}

	if req.System != "" {
		params.System = anthropic.F([]anthropic.TextBlockParam{
			{
				Type: anthropic.F(anthropic.TextBlockParamTypeText),
				Text: anthropic.F(req.System),
			},
		})
	}

	resp, err := ap.Client.Messages.New(ctx, params, option.WithRequestTimeout(15*60*time.Second))
	if err != nil {
		return nil, err
	}

	if len(resp.Content) == 0 {
		return nil, ErrAnthropicEmptyResponse
	}

	var responseContentBlocks []ContentBlock
	for _, block := range resp.Content {
		if block.Type == "text" {
			if req.JsonResponse {

				responseContentBlocks = append(responseContentBlocks, &ContentBlockText{
					Value: ExtractJSON(block.Text),
				})
			} else {
				responseContentBlocks = append(responseContentBlocks, &ContentBlockText{
					Value: block.Text,
				})
			}
		}
	}

	stopReason := convertAnthropicStopReason(string(resp.StopReason))
	usage := &ModelUsage{
		InputTokens:  int(resp.Usage.InputTokens),
		OutputTokens: int(resp.Usage.OutputTokens),
		TotalTokens:  int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
	}

	return &ConversationResponse{
		ProviderName: "Anthropic",
		ModelId:      modelId,
		Content:      responseContentBlocks,
		StopReason:   stopReason,
		Usage:        usage,
	}, nil
}

type AnthropicStreamOutput struct {
	stream *AnthropicStream
}

func (o *AnthropicStreamOutput) GetStream() ChatStream {
	return o.stream
}

type AnthropicStream struct {
	eventCh chan StreamChunk
}

func (s *AnthropicStream) Events() <-chan StreamChunk {
	return s.eventCh
}

func (ap *AnthropicProvider) ChatStream(ctx context.Context, req *ConversationRequest) (ChatStreamOutput, error) {
	if ap.Client == nil {
		return nil, ErrAnthropicClientNotInit
	}

	modelId := req.ModelId
	if modelId == "" {
		if envModelId := os.Getenv("ANTHROPIC_MODEL_ID"); envModelId != "" {
			modelId = envModelId
		} else {
			if ap.VertexAI {
				modelId = defaultVertexAIModelId
			} else {
				modelId = defaultAnthropicModelId
			}
		}
	}

	temperature := req.Temperature
	if req.Reasoning {
		temperature = 1.0
	}

	maxTokens := defaultAnthropicMaxTokens
	if req.MaxTokens > 0 {
		maxTokens = req.MaxTokens
	}

	messages := prepareAnthropicMessages(req.Messages)

	params := anthropic.MessageNewParams{
		Model:       anthropic.F(modelId),
		MaxTokens:   anthropic.F(int64(maxTokens)),
		Temperature: anthropic.F(float64(temperature)),
		Messages:    anthropic.F(messages),
	}

	if req.Reasoning {
		params.Thinking = anthropic.F[anthropic.ThinkingConfigParamUnion](anthropic.ThinkingConfigEnabledParam{
			BudgetTokens: anthropic.F(int64(4096)),
			Type:         anthropic.F(anthropic.ThinkingConfigEnabledTypeEnabled),
		})
	}

	if req.System != "" {
		params.System = anthropic.F([]anthropic.TextBlockParam{
			{
				Type: anthropic.F(anthropic.TextBlockParamTypeText),
				Text: anthropic.F(req.System),
			},
		})
	}

	stream := ap.Client.Messages.NewStreaming(ctx, params)

	eventCh := make(chan StreamChunk, defaultAnthropicStreamBufferSize)
	anthropicStream := &AnthropicStream{
		eventCh: eventCh,
	}

	go func() {
		defer close(eventCh)
		message := anthropic.Message{}

		for stream.Next() {
			select {
			case <-ctx.Done():
				eventCh <- &StreamChunkError{
					Error: ctx.Err().Error(),
				}
				return
			default:
				event := stream.Current()
				message.Accumulate(event)
				//fmt.Println("Event: ", event)

				switch event := event.AsUnion().(type) {
				case anthropic.ContentBlockDeltaEvent:
					delta := event.Delta
					if delta.Text != "" {
						eventCh <- &StreamChunkContent{
							Content: &ContentBlockText{
								Value: delta.Text,
							},
						}
					}
				case anthropic.MessageStopEvent:
					eventCh <- &StreamChunkStop{
						StopReason: convertAnthropicStopReason(string(message.StopReason)),
					}

					eventCh <- &StreamChunkUsage{
						ModelId: modelId,
						Usage: &ModelUsage{
							InputTokens:  int(message.Usage.InputTokens),
							OutputTokens: int(message.Usage.OutputTokens),
							TotalTokens:  int(message.Usage.InputTokens + message.Usage.OutputTokens),
						},
					}
				}
			}
		}

		if err := stream.Err(); err != nil {
			eventCh <- &StreamChunkError{
				Error: err.Error(),
			}
		}
	}()

	return &AnthropicStreamOutput{
		stream: anthropicStream,
	}, nil
}

func prepareAnthropicMessages(messages []Message) []anthropic.MessageParam {
	var anthropicMessages []anthropic.MessageParam
	for _, msg := range messages {
		var contentBlocks []anthropic.ContentBlockParamUnion
		for _, content := range msg.Content {
			if textBlock, ok := content.(*ContentBlockText); ok {
				contentBlocks = append(contentBlocks, anthropic.NewTextBlock(textBlock.Value))
			}
		}

		if msg.Role == UserRole {
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(contentBlocks...))
		} else {
			anthropicMessages = append(anthropicMessages, anthropic.MessageParam{
				Role:    anthropic.F(anthropic.MessageParamRoleAssistant),
				Content: anthropic.F(contentBlocks),
			})
		}
	}

	return anthropicMessages
}

func convertAnthropicStopReason(reason string) StopReason {
	switch reason {
	case "end_turn":
		return StopReasonStop
	case "max_tokens":
		return StopReasonLength
	case "tool_use":
		return StopReasonToolCalls
	default:
		return StopReasonStop
	}
}
