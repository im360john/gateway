package mcp

// NotificationContext represents the context for server notifications.
type NotificationContext struct {
	// ClientID is the unique identifier of the client.
	ClientID string `json:"clientId"`

	// SessionID is the unique identifier of the session.
	SessionID string `json:"sessionId"`
}

// LogNotification represents a log message notification from the server.
type LogNotification struct {
	Notification
	Params struct {
		// Level is the severity level of the log message.
		Level string `json:"level"`

		// Message is the log message content.
		Message string `json:"message"`
	} `json:"params"`
}

// ResourcesListChangedNotification represents a notification that the list of available resources has changed.
type ResourcesListChangedNotification struct {
	Notification
}

// PromptsListChangedNotification represents a notification that the list of available prompts has changed.
type PromptsListChangedNotification struct {
	Notification
}

// ToolsListChangedNotification represents a notification that the list of available tools has changed.
type ToolsListChangedNotification struct {
	Notification
}

// ResourceChangedNotification represents a notification that a resource has changed.
type ResourceChangedNotification struct {
	Notification
	Params struct {
		// URI is the URI of the resource that changed.
		URI string `json:"uri"`
	} `json:"params"`
}
