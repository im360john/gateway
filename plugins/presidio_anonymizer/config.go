package presidioanonymizer

// Config represents the configuration for the Presidio Anonymizer plugin
type Config struct {
	// PresidioURL is the URL of the Presidio Anonymizer API
	PresidioURL string `json:"presidio_url" yaml:"presidio_url"`
	// AnonymizerRules defines the anonymization rules for specific fields
	AnonymizerRules map[string][]AnonymizerRule `json:"anonymizer_rules" yaml:"anonymizer_rules"`
}

// AnonymizerRule defines how to anonymize a specific type of PII
type AnonymizerRule struct {
	// Type of PII to detect (e.g. "PERSON", "PHONE_NUMBER", etc.)
	Type string `json:"type" yaml:"type"`
	// Operator defines the anonymization operation ("replace" or "mask")
	Operator string `json:"operator" yaml:"operator"`
	// NewValue is used with "replace" operator
	NewValue string `json:"new_value,omitempty" yaml:"new_value,omitempty"`
	// MaskingChar is used with "mask" operator
	MaskingChar string `json:"masking_char,omitempty" yaml:"masking_char,omitempty"`
	// CharsToMask is used with "mask" operator
	CharsToMask int `json:"chars_to_mask,omitempty" yaml:"chars_to_mask,omitempty"`
}

func (c Config) Doc() string {
	return docString
}

func (c Config) Tag() string {
	return "presidio_anonymizer"
} 