package presidioanonymizer

// PresidioAnonymizer represents configuration for anonymization
type PresidioAnonymizer struct {
	Type        string `json:"type"`
	NewValue    string `json:"new_value,omitempty"`
	MaskingChar string `json:"masking_char,omitempty"`
	CharsToMask int    `json:"chars_to_mask,omitempty"`
	HashType    string `json:"hash_type,omitempty"`
	CryptoKey   string `json:"key,omitempty"`
	FromEnd     bool   `json:"from_end"`
}

// analyzerRequest represents request to Presidio Analyzer API
type analyzerRequest struct {
	Text             string            `json:"text"`
	AnalyzeTemplates []analyzeTemplate `json:"analyze_templates"`
	Language         string            `json:"language"`
}

// analyzeTemplate represents template for analysis
type analyzeTemplate struct {
	PII        string  `json:"pii_type"`
	Score      float64 `json:"score,omitempty"`
	EntityType string  `json:"entity_type,omitempty"`
}

// analyzerResponse represents response from Presidio Analyzer API
type analyzerResponse struct {
	Results []analyzeResult `json:"results"`
}

// analyzeResult represents single analysis result
type analyzeResult struct {
	Type       string  `json:"entity_type"`
	Score      float64 `json:"score"`
	StartIndex int     `json:"start"`
	EndIndex   int     `json:"end"`
	Text       string  `json:"text"`
}

// anonymizeRequest represents request to Presidio Anonymizer API
type anonymizeRequest struct {
	Text        string                        `json:"text"`
	Anonymizers map[string]PresidioAnonymizer `json:"anonymizers"`
	Analyzer    []analyzeResult               `json:"analyzer_results"`
}

// anonymizeResponse represents response from Presidio Anonymizer API
type anonymizeResponse struct {
	Text string `json:"text"`
}
