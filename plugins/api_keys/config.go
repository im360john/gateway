package api_keys

// Config represents API key authentication plugin configuration
type Config struct {
	// Name specifies the header or query parameter name for the API key
	// Example: "X-API-Key" or "api_key"
	Name string `yaml:"name"`

	// Keys is a list of API keys and their permissions
	Keys []Key `yaml:"keys"`

	// KeysFile is a path to file containing API keys (optional)
	KeysFile string `yaml:"keys_file"`

	// Location specifies where to look for the API key: "header" or "query"
	Location string `yaml:"location"`
}

func (c Config) Tag() string {
	return "api_keys"
}

func (c Config) Doc() string {
	return docString
}

// Key represents a single API key configuration with its permissions
type Key struct {
	// Key is the API key value
	Key string `yaml:"key"`

	// AllowedMethods specifies which HTTP methods this key can use
	// If empty, all methods are allowed
	AllowedMethods []string `yaml:"allowed_methods"`
}

func (t *Key) Allowed(method string) bool {
	if len(t.AllowedMethods) == 0 {
		return true
	}
	for _, m := range t.AllowedMethods {
		if m == method {
			return true
		}
	}
	return false
}
