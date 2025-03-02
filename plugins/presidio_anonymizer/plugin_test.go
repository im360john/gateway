package presidioanonymizer

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPluginWithContainer(t *testing.T) {
	ctx := context.Background()

	// Define Analyzer container
	analyzerReq := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/presidio-analyzer:latest",
		ExposedPorts: []string{"3000/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("3000/tcp"),
	}

	// Start Analyzer container
	analyzerContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: analyzerReq,
		Started:          true,
	})
	require.NoError(t, err)

	// Get the Analyzer container's host and port
	analyzerHost, err := analyzerContainer.Host(ctx)
	require.NoError(t, err)
	analyzerPort, err := analyzerContainer.MappedPort(ctx, "3000")
	require.NoError(t, err)

	// Define Anonymizer container
	anonymizerReq := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/presidio-anonymizer:latest",
		ExposedPorts: []string{"3000/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("3000/tcp"),
	}

	// Start Anonymizer container
	anonymizerContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: anonymizerReq,
		Started:          true,
	})
	require.NoError(t, err)

	// Get the Anonymizer container's host and port
	anonymizerHost, err := anonymizerContainer.Host(ctx)
	require.NoError(t, err)
	anonymizerPort, err := anonymizerContainer.MappedPort(ctx, "3000")
	require.NoError(t, err)

	// Create plugin instance
	plugin, err := New(Config{
		AnonymizeURL: fmt.Sprintf("http://%s:%s/anonymize", anonymizerHost, anonymizerPort.Port()),
		AnalyzerURL:  fmt.Sprintf("http://%s:%s/analyze", analyzerHost, analyzerPort.Port()),
		Language:     "en",
		AnonymizerRules: []AnonymizerRule{
			{
				Type:     "PERSON",
				Operator: "redact",
			},
			{
				Type:        "EMAIL_ADDRESS",
				Operator:    "mask",
				MaskingChar: "*",
				CharsToMask: 4,
			},
			{
				Type:     "PHONE_NUMBER",
				Operator: "replace",
				NewValue: "[PHONE REMOVED]",
			},
		},
	})
	require.NoError(t, err)

	// Test cases
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "anonymize multiple fields",
			input: map[string]any{
				"email":       "john.doe@example.com",
				"description": "Contact John Doe at john.doe@example.com or +1-555-123-4567",
			},
			expected: map[string]interface{}{
				"email":       "****.doe@example.com",
				"description": "Contact <PERSON> at ****.doe@example.com or +<IN_PAN>4567",
			},
		},
		{
			name: "no sensitive data",
			input: map[string]any{
				"text":   "Hello World",
				"number": 42,
			},
			expected: map[string]any{
				"text":   "Hello World",
				"number": 42,
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, skipped := plugin.Process(tt.input, nil)
			assert.False(t, skipped)
			assert.Equal(t, tt.expected, result)
		})
	}
}
