<h1 align="center">CentralMind Gateway: AI-First Data Gateway</h1>

<div align="center">

## 🛸 Introduction

</div>

AI agents and LLM-powered applications need fast, secure access to data, but traditional APIs and databases aren’t built for this. We’re building an API layer that automatically generates secure, LLM-optimized APIs on top of your structured data.

- Filters out PII and sensitive data to ensure compliance with GDPR, CPRA, SOC 2, and other regulations.
- Adds traceability and auditing, so AI applications aren’t black boxes and security teams can control.
- Optimizes for AI workloads, supports Model Context Protocol (MCP) with extra meta information to help AI agents understand APIs, caching and security.

Our first users are companies deploying AI agents for customer support and analytics, where they need models to access the right data without security risks or compliance headaches.


<div align="center">

## Design

</div>

![img.png](assets/diagram.png)

## How to generate

Gateway is LLM-model first, i.e. it's designed to be generated via LLM-models.
To generate your gateway config simply run discover command with your connection info:

1. Connection info
   ```yaml
   hosts:
   - localhost
   user: postgres
   password: password
   database: mydb
   port: 5432
   ```
2. Discovery command
   ```shell
   gateway start  \
      --config PATH_TO_CONFIG \
      discover \
      --db-type postgres \
      --tables table_name_1 --tables table_name_2 \ 
      --open-ai-key $TOKEN \
      --prompt "Generate for me awesome readonly api"
   ```
3. Wait for completion
   ```shell
   INFO[0000] Step 1: Read configs                         
   INFO[0000] Step 2: Discover data                        
   INFO[0000] Step 2: Found: 8 tables                      
   INFO[0000] Step 3: Prepare prompt to AI                 
   INFO[0000] Step 3 done. Prompt: prompt.txt              
   INFO[0000] Step 4: Do AI Magic                          
   INFO[0074] Step 4: open-ai usage: {1813 8665 10478 0x140000ce910 0x140000ce920}
   INFO[0074] ✅ API schema saved в gateway.yaml            
   INFO[0074] Done: in 1m14.140552125s
   ```
4. Explore results, the result would be saved in output file:
   ```yaml
   api:
       name: Awesome Readonly API
       description: ""
       version: "1.0"
   database:
       type: YOUR_DB_TYPE
       connection: YOUR_CONNECTION_INFO
       tables:
           - name: table_name_1
             columns:
               ... // Columns for this table
             endpoints:
               - http_method: GET
                 http_path: /some_path
                 mcp_method: some_method
                 summary: Some readable summary.
                 description: 'Some description'
                 query: SQL Query with params
                 params:
                   ... // List of params for query
   ```


## How to run

```shell
go build .
./gateway start --config ./example/gatиeway.yaml
```

### Docker compose

```shell
docker compose up ./example/docker-compose.yml
```

### MCP Protocol

Gateway implement MCP protocol, for easy access to your data right from claude, to use it

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
            "args": [
                "start", 
                "--config",
                "PATH_TO_GATEWAY_YAML_CONFIG", 
                "mcp-stdio"
            ]
        }
    }
   }
   ```
3. Ask something regards your data:
   ![claude_integration.png](./assets/claude_integration.png)

