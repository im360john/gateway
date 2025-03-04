package mcp

// Resource represents a resource that can be read by the client.
type Resource struct {
	Annotated
	// The URI of this resource.
	URI string `json:"uri"`
	// A human-readable name for this resource.
	//
	// This can be used by clients to populate UI elements.
	Name string `json:"name"`
	// A description of what this resource represents.
	//
	// This can be used by clients to improve the LLM's understanding of
	// available resources. It can be thought of like a "hint" to the model.
	Description string `json:"description,omitempty"`
	// The MIME type of this resource, if known.
	MIMEType string `json:"mimeType,omitempty"`
}

// ResourceTemplate represents a template for constructing resource URIs.
type ResourceTemplate struct {
	Annotated
	// A URI template (according to RFC 6570) that can be used to construct
	// resource URIs.
	URITemplate string `json:"uriTemplate"`
	// A human-readable name for the type of resource this template refers to.
	//
	// This can be used by clients to populate UI elements.
	Name string `json:"name"`
	// A description of what this template is for.
	//
	// This can be used by clients to improve the LLM's understanding of
	// available resources. It can be thought of like a "hint" to the model.
	Description string `json:"description,omitempty"`
	// The MIME type for all resources that match this template. This should only
	// be included if all resources matching this template have the same type.
	MIMEType string `json:"mimeType,omitempty"`
}

// ResourceContents is an interface for different types of resource contents.
type ResourceContents interface {
	isResourceContents()
}

// TextResourceContents represents text-based resource contents.
type TextResourceContents struct {
	// The URI of this resource.
	URI string `json:"uri"`
	// The MIME type of this resource, if known.
	MIMEType string `json:"mimeType,omitempty"`
	// The text of the item. This must only be set if the item can actually be
	// represented as text (not binary data).
	Text string `json:"text"`
}

func (TextResourceContents) isResourceContents() {}

// BlobResourceContents represents binary resource contents.
type BlobResourceContents struct {
	// The URI of this resource.
	URI string `json:"uri"`
	// The MIME type of this resource, if known.
	MIMEType string `json:"mimeType,omitempty"`
	// A base64-encoded string representing the binary data of the item.
	Blob string `json:"blob"`
}

func (BlobResourceContents) isResourceContents() {}

// ListResourcesRequest is a request to list available resources.
type ListResourcesRequest struct {
	PaginatedRequest
}

// ListResourcesResult is the response to a list resources request.
type ListResourcesResult struct {
	PaginatedResult
	Resources []Resource `json:"resources"`
}

// ListResourceTemplatesRequest is a request to list available resource templates.
type ListResourceTemplatesRequest struct {
	PaginatedRequest
}

// ListResourceTemplatesResult is the response to a list resource templates request.
type ListResourceTemplatesResult struct {
	PaginatedResult
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
}

// ReadResourceRequest is a request to read a resource.
type ReadResourceRequest struct {
	Request
	Params struct {
		// The URI of the resource to read. The URI can use any protocol; it is up
		// to the server how to interpret it.
		URI string `json:"uri"`
		// Arguments to pass to the resource handler
		Arguments map[string]interface{} `json:"arguments,omitempty"`
	} `json:"params"`
}

// ReadResourceResult is the response to a read resource request.
type ReadResourceResult struct {
	Result
	Contents []ResourceContents `json:"contents"` // Can be TextResourceContents or BlobResourceContents
}

// ResourceListChangedNotification is sent when the list of available resources changes.
type ResourceListChangedNotification struct {
	Notification
}

// SubscribeRequest is a request to subscribe to resource updates.
type SubscribeRequest struct {
	Request
	Params struct {
		// The URI of the resource to subscribe to. The URI can use any protocol; it
		// is up to the server how to interpret it.
		URI string `json:"uri"`
	} `json:"params"`
}

// UnsubscribeRequest is a request to unsubscribe from resource updates.
type UnsubscribeRequest struct {
	Request
	Params struct {
		// The URI of the resource to unsubscribe from.
		URI string `json:"uri"`
	} `json:"params"`
}

// ResourceUpdatedNotification is sent when a subscribed resource is updated.
type ResourceUpdatedNotification struct {
	Notification
	Params struct {
		// The URI of the resource that has been updated. This might be a sub-
		// resource of the one that the client actually subscribed to.
		URI string `json:"uri"`
	} `json:"params"`
}
