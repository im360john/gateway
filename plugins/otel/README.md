---
title: OpenTelemetry Plugin
---

Integrates OpenTelemetry tracing.

## Type
- Wrapper

## Description
Adds OpenTelemetry tracing support for database operations.

## Configuration

```yaml
otel:
  exporter_type: "oltp"              # "oltp" or "stdout"
  service_name: "my-gateway"
  service_version: "1.0.0"
  environment: "production"
  endpoint: "localhost:4317"
  tls_mode: "insecure"
  span_max_queue_size: 5000
  span_max_export_batch: 512
  batch_timeout: "1s"
  resource_attributes:               # Additional trace attributes
    team: "backend"
    region: "us-east-1"
``` 