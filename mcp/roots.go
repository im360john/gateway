package mcp

// ListRootsRequest is a request to list available roots.
type ListRootsRequest struct {
	Request
}

// ListRootsResult is the response to a list roots request.
type ListRootsResult struct {
	Result
	Roots []Root `json:"roots"`
}

// Root represents a root directory.
type Root struct {
	// The URI identifying the root. This *must* start with file:// for now.
	// This restriction may be relaxed in future versions of the protocol to allow
	// other URI schemes.
	URI string `json:"uri"`
	// An optional name for the root. This can be used to provide a human-readable
	// identifier for the root, which may be useful for display purposes or for
	// referencing the root in other parts of the application.
	Name string `json:"name,omitempty"`
}

// RootsListChangedNotification is sent when the list of available roots changes.
type RootsListChangedNotification struct {
	Notification
}
