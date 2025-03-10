package providers

import (
	"strings"
)

func ExtractJSON(text string) string {
	// Try to extract JSON objects (starting with { and ending with })
	firstBrace := strings.Index(text, "{")
	if firstBrace >= 0 {
		lastBrace := strings.LastIndex(text, "}")
		if lastBrace > firstBrace {
			return strings.TrimSpace(text[firstBrace : lastBrace+1])
		}
	}

	// If no JSON object found, try JSON arrays (starting with [ and ending with ])
	firstBracket := strings.Index(text, "[")
	if firstBracket >= 0 {
		lastBracket := strings.LastIndex(text, "]")
		if lastBracket > firstBracket {
			return strings.TrimSpace(text[firstBracket : lastBracket+1])
		}
	}

	// If no JSON found, return original string
	return text
}
