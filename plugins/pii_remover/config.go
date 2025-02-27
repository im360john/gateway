package piiremover

// Config represents PII removal configuration
type Config struct {
	// Fields specifies which fields should be checked for PII
	// Can use wildcards, e.g., "user.*" or "*_email"
	Fields []string `yaml:"fields"`

	// Replacement is the string to use instead of PII values
	// Default: "[REDACTED]"
	Replacement string `yaml:"replacement"`

	// DetectionRules defines custom regex patterns for PII detection
	DetectionRules map[string]string `yaml:"detection_rules"`
}

func (c Config) Tag() string {
	return "pii_remover"
}

func (c Config) Doc() string {
	return docString
}
