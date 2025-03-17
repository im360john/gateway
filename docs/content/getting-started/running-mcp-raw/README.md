---
title: 'Launching an MCP Raw'
---

This guide explains how to use Gateway in MCP Raw mode, which provides direct database access through a structured interface.

## Prerequisites

Before launching MCP Raw mode, ensure you have:

1. Gateway installed using one of the [installation methods](/docs/content/getting-started/installation)
2. A generated `gateway.yaml` configuration file from the [API generation process](/docs/content/getting-started/generating-api)

## Starting the MCP Raw Server

To start Gateway in MCP Raw mode, use the following command:

```bash
# Launch Gateway in MCP Raw mode
./gateway start --config gateway.yaml --raw
```

### Parameter Descriptions:

- `start`: Initiates the Gateway service
- `--config gateway.yaml`: Path to your generated API configuration file
- `--raw`: Specifies that you want to start an MCP Raw server

## Understanding MCP Raw Tools

MCP Raw mode provides four main tools for data interaction:

### 1. List Tables Tool

Returns a list of available tables in the database:

```sql
-- Will return all available tables with their structures
```

### 2. Discover Data Tool

Examines data structure and provides samples:

```sql
-- Accepts comma-separated table names
-- Example: table1,table2,table3
```

### 3. Prepare Query Tool

Validates SQL queries before execution:

```sql
-- Verifies query structure and returns expected output schema
SELECT * FROM table WHERE id = 123
```

### 4. Query Tool

Executes the SQL query and returns results:

```sql
-- Executes the query with provided parameters
SELECT * FROM table WHERE id = 123
```

## Usage Flow

The typical workflow follows these steps:

1. Use `list_tables` to view available tables
2. Use `discover_data` to examine table structures and sample data
3. Use `prepare_query` to validate your SQL query
4. Use `query` to execute the query and get results

## Configuration

Like other Gateway modes, MCP Raw supports environment variables in configuration:

```yaml
database:
  type: db-type
  connection:
    host: ${DB_HOST}
    user: ${DB_USER}
    password: ${DB_PASSWORD}
    database: ${DB_NAME}
plugins: {}
```

## Available Plugins

You can enhance your API with various plugins:
- <a href="../../plugins/pii_remover/"> PII Data reduction using Regex </a>
- <a href="../../plugins/presidio_anonymizer/"> PII Data reduction using Microsoft Presidio</a>
- <a href="../../plugins/oauth/"> OAuth Authentication</a>
- <a href="../../plugins/api_keys/"> API Keys Management</a>
- <a href="../../plugins/lru_cache/"> LRU Cache</a>
- <a href="../../plugins/otel/"> OpenTelemetry Integration</a>
- And more...


## Additional Configuration

To modify the generated API:
- Edit the `gateway.yaml` configuration
- Add custom plugins if needed
- Modify the discovery process configuration via custom prompt and restart discovery process
- Restart the server to apply changes

## Important Notes

1. Queries must be compatible with the database type specified in your configuration
2. All queries pass through plugin and interceptor systems
3. Parameter binding is mandatory for security
4. Sensitive information should be managed through environment variables

## Error Handling

The server provides detailed error messages for:
- Invalid SQL syntax
- Schema validation failures
- Connection issues
- Permission problems

