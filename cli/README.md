---
title: 'Gateway CLI'
---

This document provides information about the available CLI commands and their parameters for the Gateway application.

## Available Commands

### `gateway connectors`

List all available database connectors

**Usage:**

```
gateway connectors
```




### `gateway discover`

Discover generates gateway config

**Usage:**

```
gateway discover [flags]
```

**Flags:**

- `ai-api-key` - AI API token
- `ai-endpoint` - Custom OpenAI-compatible API endpoint URL
- `ai-max-tokens` - Maximum tokens to use (default: "0")
- `ai-model` - AI model to use
- `ai-provider` - AI provider to use (default: "openai")
- `ai-reasoning` - Enable reasoning (default: "true")
- `ai-temperature` - AI temperature (default: "-1")
- `bedrock-region` - Bedrock region
- `config` - Path to connection yaml file (default: "connection.yaml")
- `llm-log` - Path to save the raw LLM response (default: "/Users/tserakhau/go/src/github.com/gateway/binaries/.gateway/llm_raw_response.log")
- `output` - Resulted YAML path (default: "gateway.yaml")
- `prompt` - Custom input to generate APIs (default: "generate reasonable set of APIs for this data")
- `prompt-file` - Path to save the generated prompt (default: "/Users/tserakhau/go/src/github.com/gateway/binaries/.gateway/prompt_default.txt")
- `tables` - Comma-separated list of tables to include (e.g. 'table1,table2,table3')
- `vertexai-project` - Vertex AI project
- `vertexai-region` - Vertex AI region




### `gateway generate-docs`

Generate CLI documentation

**Usage:**

```
gateway generate-docs [flags]
```

**Flags:**

- `output` - Path to output README.md file (default: "cli/README.md")




### `gateway plugins`

List all available plugins

**Usage:**

```
gateway plugins
```




### `gateway start`

Start gateway

**Usage:**

```
gateway start [flags]
```

**Flags:**

- `connection-string` - Database connection string (DSN)
- `disable-swagger` - disable Swagger UI (default: "false")
- `prefix` - prefix for protocol path
- `raw` - enable as raw protocol (default: "true")
- `type` - type of database to use (default: "postgres")




### `gateway start stdio`

MCP gateway via std-io

**Usage:**

```
gateway start stdio [flags]
```

**Flags:**

- `log-file` - path to log file (default: "/Users/tserakhau/go/src/github.com/gateway/binaries/.gateway/mcp.log")
- `raw` - enable as raw protocol (default: "false")




### `gateway verify`

Verify connection config

**Usage:**

```
gateway verify [flags]
```

**Flags:**

- `config` - Path to connection yaml file (default: "connection.yaml")
- `llm-log` - Path to save the raw LLM response (default: "/Users/tserakhau/go/src/github.com/gateway/binaries/.gateway/sample.yaml")
- `tables` - Comma-separated list of tables to include (e.g. 'table1,table2,table3')






## Configuration File

The gateway.yaml configuration file defines:

- API endpoints
- Database connections
- Security settings
- Plugin configurations

Example configuration:

```yaml
# Example gateway.yaml
api:
  # API configuration
database:
  # Database connection settings
plugins:
  # Plugin configurations
```

For detailed configuration options, please refer to the main documentation.
