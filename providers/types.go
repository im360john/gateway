package providers

import "context"

type ConversationRole string

const (
	UserRole      ConversationRole = "user"
	AssistantRole ConversationRole = "assistant"
)

type ContentBlock interface {
	isContentBlock()
}

func (*ContentBlockText) isContentBlock() {}

type ContentBlockText struct {
	Value string `json:"value"`
}

type Message struct {
	Role    ConversationRole `json:"role"`
	Content []ContentBlock   `json:"content"`
}

type ConversationRequest struct {
	ModelId     string    `json:"modelId"`
	System      string    `json:"system,omitempty"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"maxTokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	Reasoning   bool      `json:"reasoning,omitempty"`
}

type StopReason string

const (
	StopReasonStop      StopReason = "stop"
	StopReasonToolCalls StopReason = "toolCalls"
	StopReasonLength    StopReason = "length"
)

type ConversationResponse struct {
	ProviderName string         `json:"providerName,omitempty"`
	ModelId      string         `json:"modelId,omitempty"`
	Content      []ContentBlock `json:"content"`
	StopReason   StopReason     `json:"stopReason,omitempty"`
	Usage        *ModelUsage    `json:"usage,omitempty"`
}

type StreamChunk interface {
	isStreamChunk()
}

func (*StreamChunkError) isStreamChunk() {}

type StreamChunkError struct {
	Error string `json:"error,omitempty"`
}

func (*StreamChunkContent) isStreamChunk() {}

type StreamChunkContent struct {
	Content ContentBlock `json:"content"`
}

func (*StreamChunkStop) isStreamChunk() {}

type StreamChunkStop struct {
	StopReason StopReason `json:"stopReason,omitempty"`
}

func (*StreamChunkUsage) isStreamChunk() {}

type StreamChunkUsage struct {
	ModelId string      `json:"modelId,omitempty"`
	Usage   *ModelUsage `json:"usage,omitempty"`
}

type ModelUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

type ChatStreamOutput interface {
	GetStream() ChatStream
}

type ChatStream interface {
	Events() <-chan StreamChunk
}

type ModelProvider interface {
	GetName() string
	Chat(ctx context.Context, req *ConversationRequest) (*ConversationResponse, error)
	ChatStream(ctx context.Context, req *ConversationRequest) (ChatStreamOutput, error)
}

type ModelProviderConfig struct {
	Name            string
	Endpoint        string
	APIKey          string
	BedrockRegion   string
	VertexAIRegion  string
	VertexAIProject string
}
