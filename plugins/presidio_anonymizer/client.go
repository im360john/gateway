package presidioanonymizer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

// PresidioClient handles communication with Presidio API endpoints
type PresidioClient struct {
	analyzerURL   string
	anonymizerURL string
	httpClient    *http.Client
}

// NewPresidioClient creates a new instance of PresidioClient
func NewPresidioClient(analyzerURL, anonymizerURL string) *PresidioClient {
	return &PresidioClient{
		analyzerURL:   analyzerURL,
		anonymizerURL: anonymizerURL,
		httpClient:    &http.Client{},
	}
}

// Analyze sends request to Presidio Analyzer API
func (c *PresidioClient) Analyze(text string, templates []analyzeTemplate, language string) ([]analyzeResult, error) {
	req := analyzerRequest{
		Text:             text,
		AnalyzeTemplates: templates,
		Language:         language,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, xerrors.Errorf("error marshaling analyzer request: %w", err)
	}

	resp, err := c.httpClient.Post(c.analyzerURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, xerrors.Errorf("error calling Presidio Analyzer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, xerrors.Errorf("Presidio Analyzer API returned status code: %d, body: %s", resp.StatusCode, string(raw))
	}

	var results []analyzeResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, xerrors.Errorf("error decoding analyzer response: %w", err)
	}

	return results, nil
}

// Anonymize sends request to Presidio Anonymizer API
func (c *PresidioClient) Anonymize(text string, anonymizers map[string]PresidioAnonymizer, analyzerResults []analyzeResult) (string, error) {
	req := anonymizeRequest{
		Text:        text,
		Anonymizers: anonymizers,
		Analyzer:    analyzerResults,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", xerrors.Errorf("error marshaling anonymize request: %w", err)
	}

	resp, err := c.httpClient.Post(c.anonymizerURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", xerrors.Errorf("error calling Presidio Anonymizer API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", xerrors.Errorf("Presidio Anonymizer API returned status code: %d, body: %s", resp.StatusCode, string(raw))
	}

	var result anonymizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", xerrors.Errorf("error decoding anonymizer response: %w", err)
	}

	return result.Text, nil
}
