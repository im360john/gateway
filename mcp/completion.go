package mcp

// CompleteRequest is a request to complete an argument value.
type CompleteRequest struct {
	Request
	Params struct {
		Ref      interface{} `json:"ref"` // Can be PromptReference or ResourceReference
		Argument struct {
			// The name of the argument
			Name string `json:"name"`
			// The value of the argument to use for completion matching.
			Value string `json:"value"`
		} `json:"argument"`
	} `json:"params"`
}

// CompleteResult is the response to a complete request.
type CompleteResult struct {
	Result
	Completion struct {
		// An array of completion values. Must not exceed 100 items.
		Values []string `json:"values"`
		// The total number of completion options available. This can exceed the
		// number of values actually sent in the response.
		Total int `json:"total,omitempty"`
		// Indicates whether there are additional completion options beyond those
		// provided in the current response, even if the exact total is unknown.
		HasMore bool `json:"hasMore,omitempty"`
	} `json:"completion"`
}

// ResourceReference is a reference to a resource.
type ResourceReference struct {
	Type string `json:"type"`
	// The URI or URI template of the resource.
	URI string `json:"uri"`
}

// PromptReference is a reference to a prompt.
type PromptReference struct {
	Type string `json:"type"`
	// The name of the prompt or prompt template
	Name string `json:"name"`
}
