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
	Text         string              `json:"text"`
	Anonymizers  []presidioAnonymizer `json:"anonymizers"`
}

type presidioAnonymizer struct {
	Type     string `json:"type"`
	NewValue string `json:"new_value,omitempty"`
	Masking  *struct {
		Char         string `json:"char,omitempty"`
		CharsToMask int    `json:"chars_to_mask,omitempty"`
	} `json:"masking,omitempty"`
}

func (p *Plugin) Process(data map[string]any, context map[string][]string) (processed map[string]any, skipped bool) {
	for field, value := range data {
		if strVal, ok := value.(string); ok {
			if rules, exists := p.cfg.AnonymizerRules[field]; exists {
				anonymized, err := p.anonymizeText(strVal, rules)
				if err != nil {
					// Log error but continue processing
					fmt.Printf("Error anonymizing field %s: %v\n", field, err)
					continue
				}
				data[field] = anonymized
			}
		}
	}
	return data, false
}

func (p *Plugin) anonymizeText(text string, rules []AnonymizerRule) (string, error) {
	reqBody := presidioRequest{
		Text:        text,
		Anonymizers: make([]presidioAnonymizer, len(rules)),
	}

	for i, rule := range rules {
		anonymizer := presidioAnonymizer{
			Type:     rule.Type,
			NewValue: rule.NewValue,
		}
		
		if rule.Operator == "mask" {
			anonymizer.Masking = &struct {
				Char         string `json:"char,omitempty"`
				CharsToMask int    `json:"chars_to_mask,omitempty"`
			}{
				Char:         rule.MaskingChar,
				CharsToMask: rule.CharsToMask,
			}
		}
		
		reqBody.Anonymizers[i] = anonymizer
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(p.cfg.PresidioURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to call Presidio API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Presidio API returned status code: %d", resp.StatusCode)
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Text, nil
}

func New(config Config) (plugins.Interceptor, error) {
	if config.PresidioURL == "" {
		return nil, fmt.Errorf("presidio_url is required")
	}

	return &Plugin{
		cfg: config,
	}, nil
} 