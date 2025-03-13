---
title: 'Neon'
---

This connector allows you to automatically generate a secured API layer over your Neon PostgreSQL database tables using CentralMind Gateway. API methods will be automatically created based on your tables and prompt using LLM (OpenAI model).

## Prerequisites

1. Install Gateway binary following the [installation guide](https://docs.centralmind.ai/getting-started/installation)
2. Set up a Neon project with configured database tables
3. Prepare OpenAI API key or compatible OpenAI API endpoint

## Configuration

You can find your connection string in the Neon Console. Navigate to your project dashboard and click on "Connection Details" to get the connection string.

![img](../assets/neon-connection.png)

Create `connection.yaml` with your Neon database credentials:
```bash
echo "
type: postgres
connection_string: \"postgresql://user:password@ep-example-123456.us-east-2.aws.neon.tech/dbname?sslmode=require\"" > connection.yaml
```

## Discovery

1. Run the discovery process to analyze your database structure, it will take a few minutes to finish:

```bash
./gateway discover --connection connection.yaml
```

2. Verify the generated configuration in `gateway.yaml`. You can:
   - Review the API structure and SQL Queries
   - Limit tables using `--tables` parameter
   - Add plugins if needed (e.g., PII data cleaning)
   - Configure additional settings

## Run API Server

Start the REST API server:

```bash
./gateway start --config gateway.yaml rest
```

[Optional] Start the MCP SSE API server:

```bash
./gateway start --config gateway.yaml mcp
```

[Optional] Start the MCP stdio API server:

```bash
./gateway start --config gateway.yaml mcp-stdio
```

## Validate that API is working
By default the Swagger UI for API will available locally on:

```
http://localhost:9090/swagger/
```

You should see Swagger UI with ability to execute methods
![img](../assets/neon-swagger.jpg)

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

## Limitations

- SSL connection is required (enabled by default in connection string)
- Some Neon-specific features like branching are not exposed through this connector
- Connection pooling settings should be configured based on your project tier 
