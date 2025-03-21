# LlamaIndex Integration

This guide demonstrates how to integrate CentralMind Gateway with LlamaIndex to create an AI agent that can interact with your database via API endpoints using MCP.

## Prerequisites

- Python 3.8+
- LlamaIndex
- OpenAI API key or any other provider compatible with OpenAI SDK
- AutoAPI running locally

## Installation

```bash
pip install llama-index llama-index-agent-openai llama-index-llms-openai llama-index-tools-mcp
```

## Start Gateway
You can start the gateway using your database or try our demo database that is available in read-only mode. The command below will start both MCP and REST endpoints simultaneously and add default methods to the API that will help the LLM understand the data structure and make read-only queries. This is extremely powerful for analytical scenarios.
```bash
./gateway start --connection-string "postgres://postgres.erjbgpchxpyteqwhxauj:!LHdWKdju8j*MLL@aws-0-eu-central-1.pooler.supabase.com:6543/postgres"
```

## Example Usage

Here's an example of how to use AutoAPI with LlamaIndex to create an AI agent:

```python
# Setup OpenAI Agent
import os
from llama_index.agent.openai import OpenAIAgent
from llama_index.llms.openai import OpenAI
from llama_index.tools.mcp import BasicMCPClient, McpToolSpec

# Set your OpenAI API key
os.environ["OPENAI_API_KEY"] = "YOUR_API_KEY"

# Initialize MCP client with your AutoAPI server
mcp_client = BasicMCPClient("http://localhost:9090/sse")

# Create MCP tool specification
mcp_tool_spec = McpToolSpec(
    client=mcp_client,    
    # You can filter specific tools by name if needed
    # allowed_tools=["tool1", "tool2"]
)

# Convert tool specification to a list of tools
tools = mcp_tool_spec.to_tool_list()

# Initialize OpenAI LLM
llm = OpenAI(model="gpt-4")

# Create the agent with tools
agent = OpenAIAgent.from_tools(tools, llm=llm, verbose=True)

# Example queries
response1 = agent.chat("what is the base url for the server")
print(response1)

response2 = agent.chat("Show me data from Staff table")
print(response2)
```

## Features

- Natural language interaction with your database (aka "chat with your database")
- AI Agent works in agentic mode, retrying different approaches and fixing mistakes automatically
- You can add plugins like PII reduction, caching, etc. Or even set specific SQL queries as an API methods by providing `gateway.yaml`

## Configuration Options

The `McpToolSpec` class accepts the following parameters:

- `client`: The MCP client instance
- `allowed_tools`: (Optional) List of tool names to filter

## Error Handling

The agent will handle errors gracefully and provide meaningful responses when:
- The API server is not accessible
- Invalid queries are provided
- Authentication issues occur

## Best Practices

1. Always secure your API keys and sensitive information
2. Use environment variables for configuration
3. Test the agent with various queries to ensure proper functionality
4. Monitor the verbose output for debugging purposes
