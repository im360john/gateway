---
title: PII Remover Plugin
---

Removes or masks Personally Identifiable Information (PII) from query results.

## Type
- Interceptor

## Description
Scans and removes/masks PII data from query results based on field patterns and custom detection rules.

## Configuration

Below is an example of different regex rules, you dont need to use them all, just pick that you need or add a new one by yourself.

```yaml
pii_remover:
  fields:                    # Fields to check for PII
    - "*.email"
    - "users.phone"          #users is column that could have json or child property: phone
    - "*.credit_card"        
  replacement: "[REDACTED]"  # Replacement text for PII values
  detection_rules:          # Custom regex patterns for PII detection
    detection_rules:          # Custom regex patterns for PII detection
      credit_card: |
          "\d{4}-\d{4}-\d{4}-\d{4}"
      phone: |
          "\+?\d{10,12}"
      ssn: |
          "^\d{3}-\d{2}-\d{4}$"
      us_address: |
          "^(?:\d{1,5})\s[A-Za-z0-9\s\.,]{5,}(?:Avenue|Ave|Street|St|Road|Rd|Boulevard|Blvd|Lane|Ln|Drive|Dr|Way|Court|Ct|Circle|Cir|Trail|Trl)[\s,]*(?:[A-Za-z\s]{2,})?[,\s]+(?:A[KLRZ]|C[AOT]|D[CE]|FL|GA|HI|I[ADLN]|K[SY]|LA|M[ADEINOST]|N[CDEHJMVY]|O[HKR]|P[AR]|RI|S[CD]|T[NX]|UT|V[AIT]|W[AIVY])[,\s]+\d{5}(?:-\d{4})?$"
      email: |
          "\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b"
      ipv4: |
          "\b(?:\d{1,3}\.){3}\d{1,3}\b"
      ipv6: |
          "\b(?:[A-Fa-f0-9]{1,4}:){7}[A-Fa-f0-9]{1,4}\b"
      iban: |
          "\b[A-Z]{2}\d{2}[A-Z0-9]{11,30}\b"
      swift: |
          "\b[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?\b"
``` 