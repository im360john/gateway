# Gateway CLI

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
gateway start mcp
```

### `start mcp-stdio`

Starts the MCP Gateway service using stdin/stdout for communication.

**Usage:**

```
gateway start mcp-stdio [flags]
```

**Flags:**

- `--log-file` - Path to log file (default: "mcp.log")

### `discover`

Discovers and generates a gateway configuration based on database schema.

**Usage:**

```
gateway discover [flags]
```

**Flags:**

- `--database-type` - Type of database to connect to
- `--tables` - Specific tables to include in the discovery
- `--ai-api-key` - API key for AI service, for instance OpenAI key
- `--ai-endpoint` - Endpoint for AI service, compatible with OpenAI api schema
- `--ai-model` - Model name for AI service, eg o3-mini
- `--output` - Output file path for generated configuration
- `--extra-prompt` - Additional instructions for the AI model

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
