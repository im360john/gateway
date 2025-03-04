// Package mcp defines the core types and interfaces for the Model Control Protocol (MCP).
// MCP is a protocol for communication between LLM-powered applications and their supporting services.
package mcp

import "encoding/json"

// Request is the base type for all requests.
type Request struct {
	Method string `json:"method"`
	Params struct {
		Meta *struct {
			// If specified, the caller is requesting out-of-band progress
			// notifications for this request (as represented by
			// notifications/progress). The value of this parameter is an
			// opaque token that will be attached to any subsequent
			// notifications. The receiver is not obligated to provide these
			// notifications.
			ProgressToken ProgressToken `json:"progressToken,omitempty"`
		} `json:"_meta,omitempty"`
	} `json:"params,omitempty"`
}

// Params represents a map of parameter values.
type Params map[string]interface{}

// Notification is the base type for all notifications.
type Notification struct {
	Method string             `json:"method"`
	Params NotificationParams `json:"params,omitempty"`
}

// NotificationParams represents parameters for a notification.
type NotificationParams struct {
	// This parameter name is reserved by MCP to allow clients and
	// servers to attach additional metadata to their notifications.
	Meta map[string]interface{} `json:"_meta,omitempty"`

	// Additional fields can be added to this map
	AdditionalFields map[string]interface{} `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for NotificationParams.
func (p NotificationParams) MarshalJSON() ([]byte, error) {
	// Create a map to hold all fields
	m := make(map[string]interface{})

	// Add Meta if it exists
	if p.Meta != nil {
		m["_meta"] = p.Meta
	}

	// Add all additional fields
	for k, v := range p.AdditionalFields {
		// Ensure we don't override the _meta field
		if k != "_meta" {
			m[k] = v
		}
	}

	return json.Marshal(m)
}

// UnmarshalJSON implements custom JSON unmarshaling for NotificationParams.
func (p *NotificationParams) UnmarshalJSON(data []byte) error {
	// Create a map to hold all fields
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	// Initialize maps if they're nil
	if p.Meta == nil {
		p.Meta = make(map[string]interface{})
	}
	if p.AdditionalFields == nil {
		p.AdditionalFields = make(map[string]interface{})
	}

	// Process all fields
	for k, v := range m {
		if k == "_meta" {
			// Handle Meta field
			if meta, ok := v.(map[string]interface{}); ok {
				p.Meta = meta
			}
		} else {
			// Handle additional fields
			p.AdditionalFields[k] = v
		}
	}

	return nil
}

// Result is the base type for all results.
type Result struct {
	// This result property is reserved by the protocol to allow clients and
	// servers to attach additional metadata to their responses.
	Meta map[string]interface{} `json:"_meta,omitempty"`
}

// EmptyResult represents a response that indicates success but carries no data.
type EmptyResult Result

// PaginatedRequest is the base type for requests that support pagination.
type PaginatedRequest struct {
	Request
	Params struct {
		// An opaque token representing the current pagination position.
		// If provided, the server should return results starting after this cursor.
		Cursor Cursor `json:"cursor,omitempty"`
	} `json:"params,omitempty"`
}

// PaginatedResult is the base type for paginated results.
type PaginatedResult struct {
	Result
	// An opaque token representing the pagination position after the last
	// returned result.
	// If present, there may be more results available.
	NextCursor Cursor `json:"nextCursor,omitempty"`
}

// Annotated is the base type for objects that can be annotated.
type Annotated struct {
	Annotations *struct {
		// Describes who the intended customer of this object or data is.
		//
		// It can include multiple entries to indicate content useful for multiple
		// audiences (e.g., `["user", "assistant"]`).
		Audience []Role `json:"audience,omitempty"`

		// Describes how important this data is for operating the server.
		//
		// A value of 1 means "most important," and indicates that the data is
		// effectively required, while 0 means "least important," and indicates that
		// the data is entirely optional.
		Priority float64 `json:"priority,omitempty"`
	} `json:"annotations,omitempty"`
}

// ClientRequest represents a request that can be sent by a client.
type ClientRequest interface{}

// ClientNotification represents a notification that can be sent by a client.
type ClientNotification interface{}

// ClientResult represents a result that can be sent by a client.
type ClientResult interface{}

// ServerRequest represents a request that can be sent by a server.
type ServerRequest interface{}

// ServerNotification represents a notification that can be sent by a server.
type ServerNotification interface{}

// ServerResult represents a result that can be sent by a server.
type ServerResult interface{}
