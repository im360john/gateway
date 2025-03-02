---
title: 'Launch API'
---

# Overview

This guide explains how to launch the API you created using the Gateway discovery process.

## Prerequisites

Before launching your API, ensure you have:

1. Gateway installed using one of the [installation methods](./installation.md)
2. A generated `gateway.yaml` configuration file from the [API generation process](./generate-api.md)

## Starting the API Server

Once you have your configuration file, you can start the API server with a simple command:

```bash
# Launch the API using your gateway.yaml configuration
./gateway start --config gateway.yaml rest
```

### Parameter Descriptions:

- `start`: Initiates the Gateway service
- `--config gateway.yaml`: Path to your generated API configuration file
- `rest`: Specifies that you want to start a REST API server

## Accessing Your API

After running the command, Gateway will launch a web server with the following default settings:

- **Server Address**: http://localhost:9090
- **Swagger Documentation**: http://localhost:9090/swagger/

The Swagger documentation provides a complete interactive reference for all endpoints in your API, allowing you to:

- Explore available endpoints
- Test API calls directly from the browser
- View request and response schemas
- Understand parameter requirements

## Customizing Server Settings

If you need to customize the server address or port, you can modify the `gateway.yaml` file or provide command-line overrides:

```bash
# Launch on a different port
./gateway start --config gateway.yaml rest --addr :7000

```


For more detailed information on API configuration options, refer to the [Configuration Guide](./configuration.md). 