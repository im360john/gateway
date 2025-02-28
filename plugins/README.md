---
title: Plugins
---

Gateway supports a plugin system that allows extending functionality through various types of plugins:

## Plugin Types

- **Interceptor** - Processes and modifies data before it reaches the connector
- **Wrapper** - Wraps and enhances connector functionality
- **Swaggerer** - Modifies OpenAPI documentation

## Available Plugins

| Plugin | Type | Description |
|--------|------|-------------|
| api_keys | Wrapper, Swaggerer | API key authentication |
| lru_cache | Wrapper | LRU-based response caching |
| lua_rls | Interceptor | Row-level security using Lua scripts |
| otel | Wrapper | OpenTelemetry integration |
| pii_remover | Interceptor | PII data removal/masking |

## Plugin Configuration

Plugins are configured in the gateway configuration file under the `plugins` section:

```yaml
plugins:
  plugin_name:
    # plugin specific configuration
    option1: value1
    option2: value2
```

See individual plugin documentation for specific configuration options. 