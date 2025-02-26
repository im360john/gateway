package piiremover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPIIRemover(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "exact field match",
			config: Config{
				Fields:      []string{"email", "phone"},
				Replacement: "[HIDDEN]",
			},
			input: map[string]any{
				"email":    "test@example.com",
				"phone":    "+79001234567",
				"username": "john_doe",
			},
			expected: map[string]any{
				"email":    "[HIDDEN]",
				"phone":    "[HIDDEN]",
				"username": "john_doe",
			},
		},
		{
			name: "wildcard patterns",
			config: Config{
				Fields:      []string{"*.email", "users.*.phone"},
				Replacement: "***",
			},
			input: map[string]any{
				"user.email":       "test@example.com",
				"users.1.phone":    "+79001234567",
				"users.2.phone":    "+79009876543",
				"users.name":       "John",
				"customer.address": "123 Street",
				"customer.email":   "customer@example.com",
				"something.else":   "value",
			},
			expected: map[string]any{
				"user.email":       "***",
				"users.1.phone":    "***",
				"users.2.phone":    "***",
				"users.name":       "John",
				"customer.address": "123 Street",
				"customer.email":   "***",
				"something.else":   "value",
			},
		},
		{
			name: "pii detection rules",
			config: Config{
				Replacement: "[REDACTED]",
				DetectionRules: map[string]string{
					"payment": `\d{4}-\d{4}-\d{4}-\d{4}`,
					"contact": `\+?\d{10,12}`,
					"phone":   `\+?\d{10,12}`,
				},
			},
			input: map[string]any{
				"payment":   "1234-5678-9012-3456",
				"contact":   "+79001234567",
				"other":     "some text",
				"card_num":  "1234567890123456", // doesn't match format
				"phone_alt": "+7900123456",      // too short
			},
			expected: map[string]any{
				"payment":   "[REDACTED]",
				"contact":   "[REDACTED]",
				"other":     "some text",
				"card_num":  "1234567890123456",
				"phone_alt": "+7900123456",
			},
		},
		{
			name: "combined rules",
			config: Config{
				Fields: []string{"*.email", "personal.phone"},
				DetectionRules: map[string]string{
					"credit_card": `\d{4}-\d{4}-\d{4}-\d{4}`,
				},
				Replacement: "XXX",
			},
			input: map[string]any{
				"user.email":      "test@example.com",
				"personal.phone":  "+79001234567",
				"credit_card":     "1234-5678-9012-3456",
				"shipping.email":  "shipping@example.com",
				"some.other.data": "regular text",
			},
			expected: map[string]any{
				"user.email":      "XXX",
				"personal.phone":  "XXX",
				"credit_card":     "XXX",
				"shipping.email":  "XXX",
				"some.other.data": "regular text",
			},
		},
		{
			name: "default replacement value",
			config: Config{
				Fields: []string{"secret"},
			},
			input: map[string]any{
				"secret": "sensitive data",
				"public": "public data",
			},
			expected: map[string]any{
				"secret": "[REDACTED]",
				"public": "public data",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin, err := New(tt.config)
			assert.NoError(t, err)

			result, skipped := plugin.Process(tt.input, nil)
			assert.False(t, skipped)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInvalidRegex(t *testing.T) {
	config := Config{
		DetectionRules: map[string]string{
			"invalid": "[", // invalid regex pattern
		},
	}

	_, err := New(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid detection rule pattern")
}
