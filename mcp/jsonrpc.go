package mcp

// JSONRPCMessage represents either a JSONRPCRequest, JSONRPCNotification, JSONRPCResponse, or JSONRPCError
type JSONRPCMessage interface{}

// JSONRPC_VERSION is the version of JSON-RPC used by MCP.
const JSONRPC_VERSION = "2.0"

// JSONRPCRequest represents a request that expects a response.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      RequestId   `json:"id"`
	Params  interface{} `json:"params,omitempty"`
	Request
}

// JSONRPCNotification represents a notification which does not expect a response.
type JSONRPCNotification struct {
	JSONRPC string `json:"jsonrpc"`
	Notification
}

// JSONRPCResponse represents a successful (non-error) response to a request.
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      RequestId   `json:"id"`
	Result  interface{} `json:"result"`
}

// JSONRPCError represents a non-successful (error) response to a request.
type JSONRPCError struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      RequestId `json:"id"`
	Error   struct {
		// The error type that occurred.
		Code int `json:"code"`
		// A short description of the error. The message SHOULD be limited
		// to a concise single sentence.
		Message string `json:"message"`
		// Additional information about the error. The value of this member
		// is defined by the sender (e.g. detailed error information, nested errors etc.).
		Data interface{} `json:"data,omitempty"`
	} `json:"error"`
}

// Standard JSON-RPC error codes
const (
	PARSE_ERROR      = -32700
	INVALID_REQUEST  = -32600
	METHOD_NOT_FOUND = -32601
	INVALID_PARAMS   = -32602
	INTERNAL_ERROR   = -32603
)
