package presidioanonymizer

// Config represents the configuration for the Presidio Anonymizer plugin
type Config struct {
	// AnonymizeURL is the URL of the Presidio Anonymizer API
	AnonymizeURL string `json:"anonymize_url" yaml:"anonymize_url"`
	// AnalyzerURL is the URL of the Presidio Analyzer API
	AnalyzerURL string `json:"analyzer_url" yaml:"analyzer_url"`
	// Language is the language used for analysis (default: "en")
	Language string `json:"language" yaml:"language"`
	// HashType is the hash algorithm used for hash operator (e.g., "md5", "sha256")
	HashType string `json:"hash_type" yaml:"hash_type"`
	// EncryptKey is the key used for encrypt operator
	EncryptKey string `json:"encrypt_key" yaml:"encrypt_key"`
	// AnonymizerRules defines the anonymization rules that apply to detected entities
	AnonymizerRules []AnonymizerRule `json:"anonymizer_rules" yaml:"anonymizer_rules"`
}

// AnonymizerRule defines how to anonymize a specific type of PII
type AnonymizerRule struct {
	// Type of PII to detect (e.g. "PERSON", "PHONE_NUMBER", etc.)
	Type string `json:"type" yaml:"type"`
	// Operator defines the anonymization operation ("mask", "replace", "hash", "encrypt")
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
