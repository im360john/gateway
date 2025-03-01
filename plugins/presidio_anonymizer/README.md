---
title: Presidio Anonymizer Plugin 
---

# Overview

This plugin integrates with <a href="https://microsoft.github.io/presidio/anonymizer/">Microsoft's Presidio</a> Anonymizer API to anonymize sensitive data in your fields. The Presidio anonymizer is module for anonymizing detected PII text entities with desired values.

![Microsoft's Presidio demo](https://microsoft.github.io/presidio/assets/detection_flow.gif)

- **Predefined or custom PII recognizers** leveraging Named Entity Recognition, regular expressions, rule-based logic, and checksum with relevant context in multiple languages.  


## Configuration

```yaml
presidio_anonymizer:
    presidio_url: http://localhost:8080/api/v1/projects/1/anonymize
    anonymizer_rules:
      - type: EMAIL_ADDRESS
        operator: mask
        masking_char: "*"
        chars_to_mask: 4
      - type: PERSON
        operator: replace
        new_value: "[REDACTED NAME]"
      - type: PHONE_NUMBER
        operator: mask
        masking_char: "#"
        chars_to_mask: 6
```

### Configuration Parameters

- `presidio_url`: Required. The URL of your Presidio Anonymizer API endpoint.
- `anonymizer_rules`: List of anonymization rules that will be applied to all fields.

Each rule contains:
- `type`: The type of PII to detect (e.g., "PERSON", "EMAIL_ADDRESS", "PHONE_NUMBER", etc.)
- `operator`: The anonymization operation ("replace" or "mask")
- `new_value`: Used with "replace" operator - the value to replace the detected PII with
- `masking_char`: Used with "mask" operator - the character to use for masking
- `chars_to_mask`: Used with "mask" operator - number of characters to mask

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
  "email": "john****@example.com",
  "name": "[REDACTED NAME]",
  "phone": "+1-555-######67",
  "description": "Contact [REDACTED NAME] at john****@example.com or +1-555-######67"
}
``` 
