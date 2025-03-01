package presidioanonymizer

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

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

type presidioRequest struct {
	Text        string               `json:"text"`
	Anonymizers []PresidioAnonymizer `json:"anonymizers"`
}

type PresidioAnonymizer struct {
	Type     string       `json:"type"`
	NewValue string       `json:"new_value,omitempty"`
	Masking  *MaskingRule `json:"masking,omitempty"`
}

type MaskingRule struct {
	Char        string `json:"char,omitempty"`
	CharsToMask int    `json:"chars_to_mask,omitempty"`
}

func (p *Plugin) Process(data map[string]any, context map[string][]string) (processed map[string]any, skipped bool) {
	// Convert data to JSON string for batch processing
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return data, false
	}

	// Convert rules to Presidio format
	var allAnonymizers []PresidioAnonymizer
	for _, rule := range p.cfg.AnonymizerRules {
		anonymizer := PresidioAnonymizer{
			Type:     rule.Type,
			NewValue: rule.NewValue,
		}

		if rule.Operator == "mask" {
			anonymizer.Masking = &MaskingRule{
				Char:        rule.MaskingChar,
				CharsToMask: rule.CharsToMask,
			}
		}

		allAnonymizers = append(allAnonymizers, anonymizer)
	}

	// If no rules defined, return original data
	if len(allAnonymizers) == 0 {
		return data, false
	}

	// Make a single API call with all rules
	reqBody := presidioRequest{
		Text:        string(jsonData),
		Anonymizers: allAnonymizers,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return data, false
	}

	resp, err := http.Post(p.cfg.PresidioURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error calling Presidio API: %v\n", err)
		return data, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Presidio API returned status code: %d\n", resp.StatusCode)
		return data, false
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return data, false
	}

	// Parse the anonymized JSON back into the data map
	var anonymizedData map[string]any
	if err := json.Unmarshal([]byte(result.Text), &anonymizedData); err != nil {
		fmt.Printf("Error unmarshaling anonymized data: %v\n", err)
		return data, false
	}

	// Update all fields with anonymized data
	for field, val := range anonymizedData {
		data[field] = val
	}

	return data, false
}

func New(config Config) (plugins.Interceptor, error) {
	if config.PresidioURL == "" {
		return nil, fmt.Errorf("presidio_url is required")
	}

	return &Plugin{
		cfg: config,
	}, nil
}
