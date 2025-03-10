package anthropic

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/vertex"
	"github.com/centralmind/gateway/providers"
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

var _ providers.ModelProvider = (*AnthropicProvider)(nil)

func init() {
	providers.RegisterModelProvider("anthropic", NewAnthropicProvider)
	providers.RegisterModelProvider("anthropic-vertexai", NewAnthropicVertexAIProvider)
}

func (ap *AnthropicProvider) GetName() string {
	if ap.VertexAI {
		return "Anthropic VertexAI"
	}

	return "Anthropic"
}

func (ap *AnthropicProvider) CostEstimate(modelId string, usage providers.ModelUsage) float64 {
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

func NewAnthropicProvider(providerConfig providers.ModelProviderConfig) (providers.ModelProvider, error) {
	return NewAnthropicProviderIntl(providerConfig, false)
}

func NewAnthropicVertexAIProvider(providerConfig providers.ModelProviderConfig) (providers.ModelProvider, error) {
	return NewAnthropicProviderIntl(providerConfig, true)
}

func NewAnthropicProviderIntl(providerConfig providers.ModelProviderConfig, vertexAI bool) (providers.ModelProvider, error) {
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

func (ap *AnthropicProvider) Chat(ctx context.Context, req *providers.ConversationRequest) (*providers.ConversationResponse, error) {
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
	temperature = max(temperature, 0.0)
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

	var responseContentBlocks []providers.ContentBlock
	for _, block := range resp.Content {
		if block.Type == "text" {
			if req.JsonResponse {

				responseContentBlocks = append(responseContentBlocks, &providers.ContentBlockText{
					Value: providers.ExtractJSON(block.Text),
				})
			} else {
				responseContentBlocks = append(responseContentBlocks, &providers.ContentBlockText{
					Value: block.Text,
				})
			}
		}
	}

	stopReason := convertAnthropicStopReason(string(resp.StopReason))
	usage := &providers.ModelUsage{
		InputTokens:  int(resp.Usage.InputTokens),
		OutputTokens: int(resp.Usage.OutputTokens),
		TotalTokens:  int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
	}

	return &providers.ConversationResponse{
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

func (o *AnthropicStreamOutput) GetStream() providers.ChatStream {
	return o.stream
}

type AnthropicStream struct {
	eventCh chan providers.StreamChunk
}

func (s *AnthropicStream) Events() <-chan providers.StreamChunk {
	return s.eventCh
}

func (ap *AnthropicProvider) ChatStream(ctx context.Context, req *providers.ConversationRequest) (providers.ChatStreamOutput, error) {
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
	temperature = max(temperature, 0.0)
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

	eventCh := make(chan providers.StreamChunk, defaultAnthropicStreamBufferSize)
	anthropicStream := &AnthropicStream{
		eventCh: eventCh,
	}

	go func() {
		defer close(eventCh)
		message := anthropic.Message{}

		for stream.Next() {
			select {
			case <-ctx.Done():
				eventCh <- &providers.StreamChunkError{
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
						eventCh <- &providers.StreamChunkContent{
							Content: &providers.ContentBlockText{
								Value: delta.Text,
							},
						}
					}
				case anthropic.MessageStopEvent:
					eventCh <- &providers.StreamChunkStop{
						StopReason: convertAnthropicStopReason(string(message.StopReason)),
					}

					eventCh <- &providers.StreamChunkUsage{
						ModelId: modelId,
						Usage: &providers.ModelUsage{
							InputTokens:  int(message.Usage.InputTokens),
							OutputTokens: int(message.Usage.OutputTokens),
							TotalTokens:  int(message.Usage.InputTokens + message.Usage.OutputTokens),
						},
					}
				}
			}
		}

		if err := stream.Err(); err != nil {
			eventCh <- &providers.StreamChunkError{
				Error: err.Error(),
			}
		}
	}()

	return &AnthropicStreamOutput{
		stream: anthropicStream,
	}, nil
}

func prepareAnthropicMessages(messages []providers.Message) []anthropic.MessageParam {
	var anthropicMessages []anthropic.MessageParam
	for _, msg := range messages {
		var contentBlocks []anthropic.ContentBlockParamUnion
		for _, content := range msg.Content {
			if textBlock, ok := content.(*providers.ContentBlockText); ok {
				contentBlocks = append(contentBlocks, anthropic.NewTextBlock(textBlock.Value))
			}
		}

		if msg.Role == providers.UserRole {
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

func convertAnthropicStopReason(reason string) providers.StopReason {
	switch reason {
	case "end_turn":
		return providers.StopReasonStop
	case "max_tokens":
		return providers.StopReasonLength
	case "tool_use":
		return providers.StopReasonToolCalls
	default:
		return providers.StopReasonStop
	}
}
