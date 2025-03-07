This connector allows you to automatically generate an API layer over your Supabase database tables using CentralMind Gateway.

## Prerequisites

1. Install Gateway binary following the [installation guide](https://docs.centralmind.ai/getting-started/installation)
2. Set up a Supabase project with configured database tables
3. Prepare OpenAI API key or compatible OpenAI API endpoint

## Configuration

Create `connection.yaml` with your Supabase database credentials. You can find credentials on your project's home page if you will click `connect` button in the menu on top. Pick `Transaction pooler` connection string cause its have ipv4 on free tier or any other that is more suitable for your case.

![img](/assets/supabase-connection.jpg)

```bash
echo "hosts:
  - xxx.pooler.supabase.com
user: \"postgres.xxxxx\"
password: \"password\"
database: \"postgres\"
port: 6543" > connection.yaml
```

## Discovery

1.Run the discovery process to analyze your database structure, it will take a few minutes to finish:

```bash
./gateway discover --connection connection.yaml
```

2.Verify the generated configuration in `gateway.yaml`. You can:
   - Review the API structure and SQL Queries
   - Add plugins if needed (e.g., PII data cleaning)
   - Configure additional settings

## Run API Server

Start the REST API server:

```bash
./gateway start --config gateway.yaml rest
```

![img](/../../assets/supabase-swagger.jpg)

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
http://localhost:9090/
```


## Available Plugins

You can enhance your API with various plugins:
- PII Data Cleaning
- OAuth Authentication
- API Keys Management
- LRU Cache
- OpenTelemetry Integration
- And more...


## Additional Configuration

To modify the generated API:
- Edit the `gateway.yaml` configuration
- Add custom plugins if needed
- Modify the discovery process configuration via custom prompt and restart discovery process
- Restart the server to apply changes
