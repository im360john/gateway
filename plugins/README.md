---
title: Plugins
---

Gateway supports a plugin system that allows extending functionality through various types of plugins:

## Plugin Types

- **Interceptor** - Processes and modifies data before it reaches the connector
- **Wrapper** - Wraps and enhances connector functionality
- **Swaggerer** - Modifies OpenAPI documentation
- **HTTPServer** - Adds HTTP endpoints to the gateway

## Available Plugins

| Plugin | Type | Description |
|--------|------|-------------|
| api_keys | Wrapper, Swaggerer | API key authentication |
| lru_cache | Wrapper | LRU-based response caching |
| lua_rls | Interceptor | Row-level security using Lua scripts |
| oauth | Wrapper, Swaggerer, HTTPServer | OAuth 2.0 authentication with support for multiple providers (Google, GitHub, Auth0, Keycloak, Okta) |
| otel | Wrapper | OpenTelemetry integration |
| pii_remover | Interceptor | PII data removal/masking |
| presidio_anonymizer | Interceptor | Microsoft Presidio-based PII detection and anonymization |

## Plugin Configuration

Plugins are configured in the gateway configuration file under the `plugins` section:

```yaml
plugins:
  plugin_name:
    # plugin specific configuration
    option1: value1
    option2: value2
```

Each plugin has its own specific configuration options. Below are some examples:

### OAuth Plugin
```yaml
oauth:
  provider: "github"           # OAuth provider (google, github, auth0, keycloak, okta)
  client_id: "xxx"            # OAuth Client ID
  client_secret: "xxx"        # OAuth Client Secret
  redirect_url: "http://localhost:8080/oauth/callback"
  scopes:                     # Required access scopes
    - "profile"
    - "email"
```

### Presidio Anonymizer Plugin
```yaml
presidio_anonymizer:
  presidio_url: "http://localhost:8080/api/v1/projects/1/anonymize"
  anonymizer_rules:
    email:
      - type: EMAIL_ADDRESS
        operator: mask
        masking_char: "*"
        chars_to_mask: 4
    name:
      - type: PERSON
        operator: replace
        new_value: "[REDACTED]"
```

See individual plugin documentation for complete configuration options. 