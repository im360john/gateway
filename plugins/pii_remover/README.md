# PII Remover Plugin

Removes or masks Personally Identifiable Information (PII) from query results.

## Type
- Interceptor

## Description
Scans and removes/masks PII data from query results based on field patterns and custom detection rules.

## Configuration

```yaml
pii_remover:
  fields:                    # Fields to check for PII
    - "*.email"
    - "users.phone"
    - "*.credit_card"
  replacement: "[REDACTED]"  # Replacement text for PII values
  detection_rules:          # Custom regex patterns for PII detection
    credit_card: "\\d{4}-\\d{4}-\\d{4}-\\d{4}"
    phone: "\\+?\\d{10,12}"
``` 