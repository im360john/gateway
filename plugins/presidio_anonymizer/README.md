---
title: Presidio Anonymizer Plugin 
---

# Overview

This plugin integrates with [Microsoft's Presidio](https://microsoft.github.io/presidio/) to analyze and anonymize sensitive data in your fields. The plugin uses two Presidio services:
- Analyzer API for detecting PII entities
- Anonymizer API for anonymizing detected entities

This plugin integrates with <a href="https://microsoft.github.io/presidio/anonymizer/">Microsoft's 
Presidio</a> Anonymizer API to anonymize sensitive data in your fields. The Presidio anonymizer is 
module for anonymizing detected PII text entities with desired values.

![Microsoft's Presidio demo](https://microsoft.github.io/presidio/assets/detection_flow.gif)

## Configuration

```yaml
presidio_anonymizer:
    anonymize_url: http://localhost:8080/anonymize
    analyzer_url: http://localhost:8080/analyze
    language: en
    hash_type: md5    # Optional, used for hash operator
    encrypt_key: ""    # Optional, used for encrypt operator
    anonymizer_rules:
      - type: EMAIL_ADDRESS
        operator: mask
        masking_char: "*"
        chars_to_mask: 4
      - type: PERSON
        operator: replace
        new_value: "[REDACTED]"
      - type: PHONE_NUMBER
        operator: hash
      - type: CREDIT_CARD
        operator: encrypt
```

### Configuration Parameters

- `anonymize_url`: Required. The URL of your Presidio Anonymizer API endpoint.
- `analyzer_url`: Required. The URL of your Presidio Analyzer API endpoint.
- `language`: Optional. Language for the analyzer (default: "en").
- `hash_type`: Optional. Hash algorithm for "hash" operator (e.g., "md5", "sha256").
- `encrypt_key`: Optional. Encryption key for "encrypt" operator.
- `anonymizer_rules`: List of anonymization rules that will be applied to detected entities.

Each rule contains:
- `type`: The type of PII to detect (e.g., "PERSON", "EMAIL_ADDRESS", "PHONE_NUMBER", etc.)
- `operator`: The anonymization operation. Supported values:
  - `mask`: Mask the value with a character
  - `replace`: Replace with a new value
  - `hash`: Hash the value using specified algorithm
  - `encrypt`: Encrypt the value using provided key
- `masking_char`: Used with "mask" operator - the character to use for masking
- `chars_to_mask`: Used with "mask" operator - number of characters to mask
- `new_value`: Used with "replace" operator - the value to replace the detected PII with

## Example

Input:
```json
{
  "email": "john.doe@example.com",
  "name": "John Doe",
  "phone": "+1-555-123-4567",
  "description": "Contact John Doe at john.doe@example.com or +1-555-123-4567"
}
```

Output:
```json
{
  "email": "****.doe@example.com",
  "name": "John Doe",
  "phone": "+1-555-123-4567",
  "description": "Contact <PERSON> at ****.doe@example.com or +<IN_PAN>4567"
}
```

## Notes

1. The plugin first uses Presidio Analyzer to detect PII entities in the text
2. Then it applies the configured anonymization rules to the detected entities
3. If no PII is detected, the original data is returned unchanged
4. Each anonymization operator requires specific parameters:
   - `mask`: requires `masking_char` and `chars_to_mask`
   - `replace`: requires `new_value`
   - `hash`: uses global `hash_type` configuration
   - `encrypt`: uses global `encrypt_key` configuration
5. The anonymization is applied to all detected entities of the specified type in the text
