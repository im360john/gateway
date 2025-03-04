---
title: 'CentralMind'
description: 'Build a data platform for LLMs in one day. Securely connect any data source and let AI handle the rest.'
---

<h1 align="center">CentralMind Gateway: AI-First Data Gateway</h1>

<div align="center">

## 🛸 Introduction

</div>

AI agents and LLM-powered applications need fast, secure access to data, but traditional APIs and databases aren't built for this. We're building an API layer that automatically generates secure, LLM-optimized APIs on top of your structured data.

- Filters out PII and sensitive data to ensure compliance with GDPR, CPRA, SOC 2, and other regulations.
- Adds traceability and auditing, so AI applications aren't black boxes and security teams can control.
- Optimizes for AI workloads, supports Model Context Protocol (MCP) with extra meta information to help AI agents understand APIs, caching and security.

Our first users are companies deploying AI agents for customer support and analytics, where they need models to access the right data without security risks or compliance headaches.

![demo](/assets/demo.gif)

## Features

- ⚡ **Automatic API Generation** – Creates APIs using LLM based on table schema and sampled data.
- 🗄️ **Structured Database Support** – Works with <a href="https://docs.centralmind.ai/connectors/postgres/">PostgreSQL</a>, <a href="https://docs.centralmind.ai/connectors/mysql/">MySQL</a>, <a href="https://docs.centralmind.ai/connectors/clickhouse/">ClickHouse</a>, and <a href="https://docs.centralmind.ai/connectors/snowflake/">Snowflake</a> connectors.
- 🌍 **Run APIs as Rest or MCP Server** – Easily expose APIs in multiple protocols.
- 📜 **Swagger & OpenAPI 3.1.0 Documentation** – Automatically generated API documentation and OpenAPI spec.
- 🔒 **PII Cleanup** – Uses <a href="https://docs.centralmind.ai/plugins/pii_remover/">regex plugin</a> or <a href="https://docs.centralmind.ai/plugins/presidio_anonymizer/">Microsoft Presidio plugin</a> for reducting PII or sensetive data.
- ⚡ **Configurable via YAML & Plugin System** – Extend API functionality effortlessly.
- 🐳 **Run as Binary or Docker** – Comes with a ready-to-use <a href="https://docs.centralmind.ai/helm/gateway/">Helm chart</a> and docker container.
- 📦 **Local & On-Prem Usage** – Allow usage with self-hosted LLMs, just specify --ai-endpoint and --ai-model parameters.
- 🔑 **Row-Level Security (RLS)** – Restrict data access using <a href="https://docs.centralmind.ai/plugins/lua_rls/">Lua scripts</a>.
- 🔐 **Authentication** – Supports <a href="https://docs.centralmind.ai/plugins/api_keys/">API keys</a> and <a href="https://docs.centralmind.ai/plugins/oauth/">OAuth</a>.
- 👀 **Observability & Audit Trail** – Uses <a href="https://docs.centralmind.ai/plugins/otel/">OpenTelemetry (OTel) plugin</a> for tracking requests including remote endpoints.
- 🏎️ **Caching** – Supports time-based and <a href="https://docs.centralmind.ai/plugins/lru_cache/">LRU caching</a> for efficiency.

## How it Works

<div align="center">

![img.png](assets/diagram.png)

</div>

### Connect & Discover

Gateway connects to your structured databases like PostgreSQL. Automatically analyzes the schema and samples data to generate an optimized API structure based on your prompt. Ensures security by detecting PII. On this stage the tool is using AI service to generate API configuration. You can use OpenAI or any OpenAI compatible API providers.

### Deploy

Runs as a standalone binary, Docker container, or Helm chart for Kubernetes. Configuration is managed via YAML and a plugin system, allowing customization without modifying the core code. Supports row-level security (RLS) with Lua scripts, caching strategies like LRU and time-based expiration, and observability through OpenTelemetry. Cleaning PII data using regex rules.

### Use & Integrate

Exposes APIs through REST, and MCP with built-in authentication via API keys and OAuth. Designed for seamless integration with AI models, including OpenAI, Anthropic Claude, Google Gemini, and DeepSeek. Automatically provides OpenAPI 3.1.0 documentation for easy adoption and supports flexible query execution with structured access control.

## Documentation

- Getting Started
  - <a href="https://docs.centralmind.ai/docs/content/getting-started/quickstart/">Quickstart</a>
  - <a href="https://docs.centralmind.ai/docs/content/getting-started/installation/">Installation</a>
  - <a href="https://docs.centralmind.ai/docs/content/getting-started/generating-api/">Generating an API</a>
  - <a href="https://docs.centralmind.ai/docs/content/getting-started/launching-api/">Launching an API</a>
- <a href="https://docs.centralmind.ai/docs/content/integration/chatgpt/"> Integration - ChatGPT</a>
- <a href="https://docs.centralmind.ai/connectors/"> Database Connectors</a>
- <a href="https://docs.centralmind.ai/plugins/"> Plugins</a>

## How to build

```shell
   git clone https://github.com/centralmind/gateway.git

   cd gateway

   # Install dependencies
   go mod download

   # Building
   go build .

```

## How to generate API

Gateway is LLM-model first, i.e. it's designed to be generated via LLM-models.
To generate your gateway config simply run discover command with your connection info:

1. Create a connection configuration file (e.g., `connection.yaml`) with your database credentials:
   ```shell
   echo 'hosts:
         - localhost
         user: "your-database-user"
         password: "your-database-password"
         database: "your-database-name"
         port: 5432' > connection.yaml
   ```
2. Discovery command
   ```shell
   gateway discover  \
      --config connection.yaml \
      --db-type postgres \
      --tables "table_name_1,table_name_2" \
      --ai-api-key $TOKEN \
      --prompt "Generate for me awesome readonly api"
   ```
3. Wait for completion

   ```shell
      INFO 🚀 API Discovery Process
      INFO Step 1: Read configs
      INFO ✅ Step 1 completed. Done.

      INFO Step 2: Discover data
      INFO Discovered Tables:
      INFO   - payment_dim: 3 columns
      INFO   - fact_table: 9 columns
      ...
      INFO ✅ Step 2 completed. Done.

      INFO Step 3: Sample data from tables
      INFO Data Sampling Results:
      INFO   - payment_dim: 5 rows sampled
      INFO   - fact_table: 5 rows sampled
      ...
      INFO ✅ Step 3 completed. Done.

      INFO Step 4: Prepare prompt to AI
      INFO Prompt saved locally to prompt_default.txt
      INFO ✅ Step 4 completed. Done.

      INFO Step 5: Using AI to design API
      Waiting for OpenAI response... Done!
      INFO OpenAI usage:  Input tokens=3187 Output tokens=14872 Total tokens=18059
      INFO API Functions Created:
      INFO   - GET /payment_dim/{payment_key} - Retrieve a payment detail by its payment key
      .....
      .....
      INFO API schema saved to: gateway.yaml

      INFO ✅ Step 5: API Specification Generation Completed!

      INFO ✅ All steps completed. Done.

      INFO --- Execution Statistics ---
      INFO Total time taken: 2m12s
      INFO Tokens used: 18059 (Estimated cost: $0.0689)
      INFO Tables processed: 6
      INFO API methods created: 18
      INFO Total number of columns with PII data: 2
   ```

4. Explore results, the result would be saved in output file `gateway.yaml`:
   ```yaml
   api:
     name: Awesome Readonly API
     description: ''
     version: '1.0'
   database:
     type: YOUR_DB_TYPE
     connection: YOUR_CONNECTION_INFO
     tables:
       - name: table_name_1
         columns: ... // Columns for this table
         endpoints:
           - http_method: GET
             http_path: /some_path
             mcp_method: some_method
             summary: Some readable summary.
             description: 'Some description'
             query: SQL Query with params
             params: ... // List of params for query
   ```

## How to start API

```shell
gateway start --config gateway.yaml rest
```

### Docker compose

```shell
docker compose up ./example/simple/docker-compose.yml
```

### MCP Protocol

Gateway implement MCP protocol, for easy access to your data right from claude, to use it. Full instruction how to setup MCP and integrate with <a href="https://docs.centralmind.ai/docs/content/integration/claude-desktop/">Anthropic Claude is here</a>.

1. Build binary
   ```shell
   go build .
   ```
2. Add gateway to claude integrations config:
   ```json
   {
     "mcpServers": {
       "gateway": {
         "command": "PATH_TO_GATEWAY_BINARY",
         "args": ["start", "--config", "PATH_TO_GATEWAY_YAML_CONFIG", "mcp-stdio"]
       }
     }
   }
   ```
3. Ask something regards your data:
   ![claude_integration.png](./assets/claude_integration.png)

## Roadmap

- 🗄️ **Expand Database Support** – Add support for Redshift, S3, Oracle, MS SQL, Elasticsearch.
- 🔍 **Complex filters and Aggregations** - Support API methods with advanced filtering and aggregation syntax.
- 🔐 **MCP with Authentication** – Secure Model Context Protocol with API keys and OAuth.
- 🤖 **More LLM Providers** – Integrate Anthropic Claude, Google Gemini, DeepSeek.- 🏠
- 📦 **Schema Evolution & Versioning** – Track changes and auto-migrate APIs.
- 🚦 **Traffic Control & Rate Limiting** – Intelligent throttling for high-scale environments.
