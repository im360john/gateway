package mcp

// LoggingLevel represents the severity level of a log message.
type LoggingLevel string

const (
	// LoggingLevelTrace is for very detailed debugging information.
	LoggingLevelTrace LoggingLevel = "trace"
	// LoggingLevelDebug is for debugging information.
	LoggingLevelDebug LoggingLevel = "debug"
	// LoggingLevelInfo is for general information.
	LoggingLevelInfo LoggingLevel = "info"
	// LoggingLevelWarn is for warnings.
	LoggingLevelWarn LoggingLevel = "warn"
	// LoggingLevelError is for errors.
	LoggingLevelError LoggingLevel = "error"
)

// SetLevelRequest is a request to set the logging level.
type SetLevelRequest struct {
	Request
	Params struct {
		// The level of logging that the client wants to receive from the server.
		// The server should send all logs at this level and higher (i.e., more severe) to
		// the client as notifications/logging/message.
		Level LoggingLevel `json:"level"`
	} `json:"params"`
}

// LoggingMessageNotification is sent when a log message is generated.
type LoggingMessageNotification struct {
	Notification
	Params struct {
		// The severity of this log message.
		Level LoggingLevel `json:"level"`
		// An optional name of the logger issuing this message.
		Logger string `json:"logger,omitempty"`
		// The data to be logged, such as a string message or an object. Any JSON
		// serializable type is allowed here.
		Data interface{} `json:"data"`
	} `json:"params"`
}
