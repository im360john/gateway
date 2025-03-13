---
title: 'Gateway CLI'
---

This document provides information about the available CLI commands and their parameters for the Gateway application.

## Available Commands

### `start`

Starts the Gateway with the specified configuration.

**Usage:**

```
gateway start [flags]
```

**Flags:**

- `--config` - Path to YAML file with gateway configuration (default: "./gateway.yaml")
- `--addr` - Address for gateway server (default: ":9090")
- `--servers` - Comma-separated list of additional server URLs for Swagger UI (e.g., "https://dev1.example.com,https://dev2.example.com")

### `start rest`

Starts the REST Gateway service.

**Usage:**

```
gateway start rest [flags]
```

**Flags:**

- `--disable-swagger` - Disable Swagger UI (default: false)

### `start mcp`

Starts the MCP (Message-Coupling Protocol) Gateway service.

**Usage:**

```
gateway start mcp [flags]
```

This mode is particularly useful for:
- Testing and debugging MCP communication
- Integration with systems that have well known queries to execute
- Script-based automation and pipeline processing


### `start mcp-raw`

For scenarios where you need direct access to raw SQL queries through MCP protocol:

**Usage:**

```
gateway start mcp-raw [flags]
```

This mode is particularly useful for:
- Direct database access through MCP protocol
- Advanced data querying and exploration
- Development and debugging of database queries


### `start mcp-stdio`

Starts the MCP Gateway service using stdin/stdout for communication.

**Usage:**

```
gateway start mcp-stdio [flags]
```

**Flags:**

- `--log-file` - Path to log file (default: "mcp.log")

### MCP StdInOut Parameters:

- `start`: Initiates the Gateway service
- `--config gateway.yaml`: Path to your generated API configuration file
- `mcp-stdio`: Specifies that you want to use MCP with standard input/output

This mode is particularly useful for:
- Testing and debugging MCP communication
- Integration with systems that require direct stdin/stdout communication and local launching applications
- Script-based automation and pipeline processing

### `discover`

Discovers and generates a gateway configuration based on database schema.

**Usage:**

```
gateway discover [flags]
```

**Flags:**

- `--config` - Path to connection yaml file. Default is "connection.yaml".
- `--tables` - Comma-separated list of tables to include (e.g. 'table1,table2,table3')
- `--ai-provider` - AI provider to use. Default is "openai".
- `--ai-endpoint` - Custom OpenAI-compatible API endpoint URL
- `--ai-api-key` - AI API token
- `--bedrock-region` - Bedrock region
- `--vertexai-region` - Vertex AI region
- `--vertexai-project` - Vertex AI project
- `--ai-model` - AI model to use
- `--ai-max-tokens` - Maximum tokens to use. Default is 0.
- `--ai-temperature` - AI temperature. Default is -1.0.
- `--ai-reasoning` - Enable reasoning. Default is true.
- `--output` - Resulted YAML path. Default is "gateway.yaml".
- `--prompt` - Custom input to generate APIs. Default is "generate reasonable set of APIs for this data".
- `--prompt-file` - Path to save the generated prompt.
- `--llm-log` - Path to save the raw LLM response.

### `connectors`

Lists all available database connectors.

**Usage:**

```
gateway connectors [connector_name]
```

If a connector name is provided, displays detailed documentation for that connector.

### `plugins`

Lists all available plugins.

**Usage:**

```
gateway plugins [plugin_name]
```

If a plugin name is provided, displays detailed documentation for that plugin.

### `verify`

Verifies database connection configuration and displays table schemas and sample data.

**Usage:**

```
gateway verify [flags]
```

**Flags:**

- `--config` - Path to connection YAML file (default: "connection.yaml")
- `--tables` - Comma-separated list of tables to include (e.g. 'table1,table2,table3')
- `--llm-log` - Path to save the raw LLM response (default: logs directory)

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
