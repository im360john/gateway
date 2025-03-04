package mcp

// Content is an interface for different types of message content.
type Content interface {
	isContent()
}

// TextContent represents text-based message content.
type TextContent struct {
	Annotated
	Type string `json:"type"` // Must be "text"
	// The text content of the message.
	Text string `json:"text"`
}

func (TextContent) isContent() {}

// ImageContent represents image-based message content.
type ImageContent struct {
	Annotated
	Type string `json:"type"` // Must be "image"
	// The base64-encoded image data.
	Data string `json:"data"`
	// The MIME type of the image. Different providers may support different image types.
	MIMEType string `json:"mimeType"`
}

func (ImageContent) isContent() {}

// EmbeddedResource represents a resource embedded in a message.
type EmbeddedResource struct {
	Annotated
	Type     string           `json:"type"`
	Resource ResourceContents `json:"resource"`
}

func (EmbeddedResource) isContent() {}

// SamplingMessage represents a message in a sampling conversation.
type SamplingMessage struct {
	Role    Role        `json:"role"`
	Content interface{} `json:"content"` // Can be TextContent or ImageContent
}

// CreateMessageRequest is a request to create a new message.
type CreateMessageRequest struct {
	Request
	Params struct {
		Messages         []SamplingMessage `json:"messages"`
		ModelPreferences *ModelPreferences `json:"modelPreferences,omitempty"`
		SystemPrompt     string            `json:"systemPrompt,omitempty"`
		IncludeContext   string            `json:"includeContext,omitempty"`
		Temperature      float64           `json:"temperature,omitempty"`
		MaxTokens        int               `json:"maxTokens"`
		StopSequences    []string          `json:"stopSequences,omitempty"`
		Metadata         interface{}       `json:"metadata,omitempty"`
	} `json:"params"`
}

// CreateMessageResult is the response to a create message request.
type CreateMessageResult struct {
	Result
	SamplingMessage
	// The name of the model that generated the message.
	Model string `json:"model"`
	// The reason why sampling stopped, if known.
	StopReason string `json:"stopReason,omitempty"`
}
