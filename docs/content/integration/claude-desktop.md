---
title: 'Claude Desktop Integration'
---

This guide explains how to integrate your Gateway API with Claude Desktop by adding Gateway as a custom tool.

## Prerequisites

Before integrating with Claude Desktop, ensure you have:

1. Claude Desktop application installed on your computer (https://claude.ai/download)
2. Gateway installed and configured with a valid gateway.yaml file

## Adding Gateway as a Tool in Claude Desktop

Claude Desktop allows you to extend its capabilities by adding custom tools through the configuration file. Follow these steps to add your Gateway API as a tool:

### Step 1: Locate the Claude Desktop Configuration File

The Claude Desktop configuration file is typically located at:

- Windows: `C:\Users\%user%\AppData\Roaming\Claude\claude_desktop_config.json`
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

### Step 2: Modify the Configuration File

Edit the `claude_desktop_config.json` file to add Gateway as a tool. Here's an example configuration:

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

Replace the placeholders with your actual paths:

- `PATH_TO_GATEWAY_BINARY`: The full path to your Gateway executable
  - Example: `"/usr/local/bin/gateway"` or `"C:\\Program Files\\Gateway\\gateway.exe"`
- `PATH_TO_GATEWAY_YAML_CONFIG`: The full path to your gateway.yaml configuration file
  - Example: `"/path/to/your/gateway.yaml"` or `"C:\\Path\\To\\gateway.yaml"`

### Step 3: Configuration Explanation

Let's break down what each part of the configuration means:

- `mcpServers`: The main section for defining tool servers that Claude can interact with
- `gateway`: A unique identifier for your Gateway tool (you can name this anything)
- `command`: The command to execute (path to your Gateway binary)
- `args`: Command line arguments for Gateway
  - `start`: Starts the Gateway service
  - `--config`: Specifies the configuration file path
  - `PATH_TO_GATEWAY_YAML_CONFIG`: Path to your gateway.yaml file
  - `mcp-stdio`: Special mode that allows Gateway to communicate with Claude Desktop

### Step 4: Save and Restart Claude Desktop

1. Save the configuration file
2. Close Claude Desktop completely
3. Restart Claude Desktop for the changes to take effect

## Using Gateway in Claude Desktop

Once configured, you can interact with your Gateway API directly in Claude Desktop:

1. Start a new conversation in Claude Desktop
2. Ask Claude to perform tasks that require accessing your API
3. Claude will automatically use the Gateway tool to access your data when needed


For more details on Claude Desktop tool configuration, refer to the [Claude Desktop Documentation](https://docs.anthropic.com/claude/docs/claude-desktop-tools).
