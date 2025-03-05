---
title: 'Launching an API'
---


This guide explains how to launch the API you created using the Gateway discovery process.

## Prerequisites

Before launching your API, ensure you have:

1. Gateway installed using one of the [installation methods](/docs/content/getting-started/installation)
2. A generated `gateway.yaml` configuration file from the [API generation process](/docs/content/getting-started/generating-api)

## Starting the REST API Server

Once you have your configuration file, you can start the API server with a simple command:

```bash
# Launch the API using your gateway.yaml configuration
./gateway start --config gateway.yaml rest
```

### Parameter Descriptions:

- `start`: Initiates the Gateway service
- `--config gateway.yaml`: Path to your generated API configuration file
- `rest`: Specifies that you want to start a REST API server

### Accessing Your API

After running the command, Gateway will launch a web server with the following default settings:

- **Server Address**: http://localhost:9090
- **Swagger Documentation**: http://localhost:9090/swagger/

The Swagger documentation provides a complete interactive reference for all endpoints in your API, allowing you to:

- Explore available endpoints
- Test API calls directly from the browser
- View request and response schemas
- Understand parameter requirements

### Customizing Server Settings

If you need to customize the server address or port, you can modify the `gateway.yaml` file or provide command-line overrides:

```bash
# Launch on a different port
./gateway start --config gateway.yaml rest --addr :7000

```

## Launching MCP SSE Server Mode

To start Gateway in MCP (Message Communication Protocol) SSE server mode, use the following command:

```bash
# Launch the API in MCP server mode
./gateway start --config gateway.yaml mcp
```

By default it will be available on address:

```
http://localhost:9090/sse  
```

### MCP Server Parameters:

- `start`: Initiates the Gateway service
- `--addr`: Address for gateway server (default: “:9090”)
- `--config gateway.yaml`: Path to your generated API configuration file
- `mcp`: Specifies that you want to start an MCP server

The MCP server mode allows for efficient message-based communication between services.

## Using MCP StdInOut Mode

For scenarios where you need direct input/output communication, you can use the MCP StdInOut mode for example when working with local applications like Claude Desktop or Cursor:

```bash
# Launch the API in MCP StdInOut mode
./gateway start --config gateway.yaml mcp-stdio
```

### MCP StdInOut Parameters:

- `start`: Initiates the Gateway service
- `--config gateway.yaml`: Path to your generated API configuration file
- `mcp-stdio`: Specifies that you want to use MCP with standard input/output

This mode is particularly useful for:
- Testing and debugging MCP communication
- Integration with systems that require direct stdin/stdout communication and local launching applications
- Script-based automation and pipeline processing
