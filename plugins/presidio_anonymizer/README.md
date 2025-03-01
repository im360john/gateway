---
title: Presidio Anonymizer Plugin
---

# Presidio Anonymizer Plugin

This plugin integrates with Microsoft's Presidio Anonymizer API to anonymize sensitive data in your fields.

## Configuration

```yaml
presidio_url: http://localhost:8080/api/v1/projects/1/anonymize
anonymizer_rules:
  email:
    - type: EMAIL_ADDRESS
      operator: mask
      masking_char: "*"
      chars_to_mask: 4
  name:
    - type: PERSON
      operator: replace
      new_value: "[REDACTED NAME]"
  phone:
    - type: PHONE_NUMBER
      operator: mask
      masking_char: "#"
      chars_to_mask: 6
```

### Configuration Parameters

- `presidio_url`: Required. The URL of your Presidio Anonymizer API endpoint.
- `anonymizer_rules`: Map of field names to their anonymization rules.

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
  "phone": "+1-555-123-4567"
}
```

Output:
```json
{
  "email": "john****@example.com",
  "name": "[REDACTED NAME]",
  "phone": "+1-555-######67"
}
``` 