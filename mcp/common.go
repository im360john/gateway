package mcp

// Role represents the role of a message sender.
type Role string

const (
	// RoleSystem represents system messages.
	RoleSystem Role = "system"
	// RoleUser represents user messages.
	RoleUser Role = "user"
	// RoleAssistant represents assistant messages.
	RoleAssistant Role = "assistant"
	// RoleTool represents tool messages.
	RoleTool Role = "tool"
)

// ProgressToken is used to associate progress notifications with the original request.
type ProgressToken interface{}

// Cursor is an opaque token used to represent a cursor for pagination.
type Cursor string

// RequestId is a uniquely identifying ID for a request in JSON-RPC.
// It can be any JSON-serializable value, typically a number or string.
type RequestId interface{}

// PingRequest is a request to check if the server is alive.
type PingRequest struct {
	Request
}

// ProgressNotification is sent to indicate progress on a long-running request.
type ProgressNotification struct {
	Notification
	Params struct {
		// The progress token which was given in the initial request, used to
		// associate this notification with the request that is proceeding.
		ProgressToken ProgressToken `json:"progressToken"`
		// The progress thus far. This should increase every time progress is made,
		// even if the total is unknown.
		Progress float64 `json:"progress"`
		// Total number of items to process (or total progress required), if known.
		Total float64 `json:"total,omitempty"`
	} `json:"params"`
}

// CancelledNotification can be sent by either side to indicate that it is
// cancelling a previously-issued request.
//
// The request SHOULD still be in-flight, but due to communication latency, it
// is always possible that this notification MAY arrive after the request has
// already finished.
//
// This notification indicates that the result will be unused, so any
// associated processing SHOULD cease.
//
// A client MUST NOT attempt to cancel its `initialize` request.
type CancelledNotification struct {
	Notification
	Params struct {
		// The ID of the request to cancel.
		//
		// This MUST correspond to the ID of a request previously issued
		// in the same direction.
		RequestId RequestId `json:"requestId"`

		// An optional string describing the reason for the cancellation. This MAY
		// be logged or presented to the user.
		Reason string `json:"reason,omitempty"`
	} `json:"params"`
}
