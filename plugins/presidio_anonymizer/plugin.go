package presidioanonymizer

import (
	_ "embed"
	"encoding/json"
	"fmt"

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
	cfg    Config
	client *PresidioClient
}

func (p *Plugin) Doc() string {
	return docString
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

	// Call Analyzer API
	analyzerResults, err := p.client.Analyze(string(jsonData), analyzeTemplates, p.cfg.Language)
	if err != nil {
		logrus.Errorf("Error analyzing data: %v\n", err)
		return data, false
	}

	// If no PII found, return original data
	if len(analyzerResults) == 0 {
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

	// Call Anonymizer API
	anonymizedText, err := p.client.Anonymize(string(jsonData), anonymizers, analyzerResults)
	if err != nil {
		logrus.Errorf("Error anonymizing data: %v\n", err)
		return data, false
	}

	// Parse the anonymized JSON back into the data map
	var anonymizedData map[string]any
	if err := json.Unmarshal([]byte(anonymizedText), &anonymizedData); err != nil {
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

	client := NewPresidioClient(config.AnalyzerURL, config.AnonymizeURL)

	return &Plugin{
		cfg:    config,
		client: client,
	}, nil
}
