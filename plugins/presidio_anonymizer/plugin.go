package presidioanonymizer

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/centralmind/gateway/plugins"
)

//go:embed README.md
var docString string

func init() {
	plugins.Register(func(cfg Config) (plugins.Interceptor, error) {
		return New(cfg)
	})
}

type Plugin struct {
	cfg Config
}

func (p *Plugin) Doc() string {
	return docString
}

type PresidioAnonymizer struct {
	Type        string `json:"type"`
	NewValue    string `json:"new_value,omitempty"`
	MaskingChar string `json:"masking_char,omitempty"`
	CharsToMask int    `json:"chars_to_mask,omitempty"`
	HashType    string `json:"hash_type,omitempty"`
	CryptoKey   string `json:"key,omitempty"`
	FromEnd     bool   `json:"from_end"`
}

type analyzerRequest struct {
	Text             string            `json:"text"`
	AnalyzeTemplates []analyzeTemplate `json:"analyze_templates"`
	Language         string            `json:"language"`
}

type analyzeTemplate struct {
	PII        string  `json:"pii_type"`
	Score      float64 `json:"score,omitempty"`
	EntityType string  `json:"entity_type,omitempty"`
}

type analyzerResponse struct {
	Results []analyzeResult `json:"results"`
}

type analyzeResult struct {
	Type       string  `json:"entity_type"`
	Score      float64 `json:"score"`
	StartIndex int     `json:"start"`
	EndIndex   int     `json:"end"`
	Text       string  `json:"text"`
}

type anonymizeRequest struct {
	Text        string                        `json:"text"`
	Anonymizers map[string]PresidioAnonymizer `json:"anonymizers"`
	Analyzer    []analyzeResult               `json:"analyzer_results"`
}

func (p *Plugin) Process(data map[string]any, _ map[string][]string) (processed map[string]any, skipped bool) {
	// Convert data to JSON string for batch processing
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Error marshaling data: %v\n", err)
		return data, false
	}

	// Convert rules to analyzer templates
	analyzeTemplates := make([]analyzeTemplate, 0)
	for _, rule := range p.cfg.AnonymizerRules {
		template := analyzeTemplate{
			PII: rule.Type,
		}
		analyzeTemplates = append(analyzeTemplates, template)
	}

	// Call Analyzer API first
	analyzerReq := analyzerRequest{
		Text:             string(jsonData),
		AnalyzeTemplates: analyzeTemplates,
		Language:         p.cfg.Language,
	}

	analyzerBody, err := json.Marshal(analyzerReq)
	if err != nil {
		logrus.Errorf("Error marshaling analyzer request: %v\n", err)
		return data, false
	}

	analyzerResp, err := http.Post(p.cfg.AnalyzerURL, "application/json", bytes.NewBuffer(analyzerBody))
	if err != nil {
		logrus.Errorf("Error calling Presidio Analyzer API: %v\n", err)
		return data, false
	}
	defer analyzerResp.Body.Close()

	if analyzerResp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(analyzerResp.Body)
		logrus.Errorf("Presidio Analyzer API returned status code: %d\n, %s\n", analyzerResp.StatusCode, string(raw))
		return data, false
	}

	var analyzerResults analyzerResponse
	if err := json.NewDecoder(analyzerResp.Body).Decode(&analyzerResults.Results); err != nil {
		logrus.Errorf("Error decoding analyzer response: %v\n", err)
		return data, false
	}

	// If no PII found, return original data
	if len(analyzerResults.Results) == 0 {
		return data, false
	}

	// Convert rules to Presidio format
	anonymizers := make(map[string]PresidioAnonymizer)
	for _, rule := range p.cfg.AnonymizerRules {
		anonymizer := PresidioAnonymizer{}

		switch rule.Operator {
		case "mask":
			anonymizer.Type = "mask"
			anonymizer.MaskingChar = rule.MaskingChar
			anonymizer.CharsToMask = rule.CharsToMask
			anonymizer.FromEnd = false
		case "replace":
			anonymizer.Type = "replace"
			anonymizer.NewValue = rule.NewValue
		case "hash":
			anonymizer.Type = "hash"
			anonymizer.HashType = p.cfg.HashType
		case "encrypt":
			anonymizer.Type = "encrypt"
			anonymizer.CryptoKey = p.cfg.EncryptKey
		default:
			anonymizer.Type = "replace"
		}

		anonymizers[rule.Type] = anonymizer
	}

	// Call Anonymizer API with analyzer results
	anonymizeReq := anonymizeRequest{
		Text:        string(jsonData),
		Anonymizers: anonymizers,
		Analyzer:    analyzerResults.Results,
	}

	anonymizeBody, err := json.Marshal(anonymizeReq)
	if err != nil {
		logrus.Errorf("Error marshaling anonymize request: %v\n", err)
		return data, false
	}

	anonymizeResp, err := http.Post(p.cfg.AnonymizeURL, "application/json", bytes.NewBuffer(anonymizeBody))
	if err != nil {
		logrus.Errorf("Error calling Presidio Anonymizer API: %v\n", err)
		return data, false
	}
	defer anonymizeResp.Body.Close()

	if anonymizeResp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(anonymizeResp.Body)
		logrus.Errorf("Presidio Anonymizer API returned status code: %d\n%s\n", anonymizeResp.StatusCode, string(raw))
		return data, false
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(anonymizeResp.Body).Decode(&result); err != nil {
		logrus.Errorf("Error decoding anonymizer response: %v\n", err)
		return data, false
	}

	// Parse the anonymized JSON back into the data map
	var anonymizedData map[string]any
	if err := json.Unmarshal([]byte(result.Text), &anonymizedData); err != nil {
		logrus.Errorf("Error unmarshaling anonymized data: %v\n", err)
		return data, false
	}

	// Update all fields with anonymized data
	for field, val := range anonymizedData {
		data[field] = val
	}

	return data, false
}

func New(config Config) (plugins.Interceptor, error) {
	if config.AnonymizeURL == "" {
		return nil, fmt.Errorf("presidio_url is required")
	}
	if config.AnalyzerURL == "" {
		return nil, fmt.Errorf("analyzer_url is required")
	}
	if config.Language == "" {
		config.Language = "en"
	}

	return &Plugin{
		cfg: config,
	}, nil
}
