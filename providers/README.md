---
title: 'AI Providers'
description: 'Overview of supported AI providers and configuration options'
---

This guide provides an overview of the AI providers supported by CentralMind Gateway, along with configuration options and examples.

## Supported Providers

We support the following AI providers:

- [**OpenAI**](/providers/openai) and all OpenAI-compatible providers
- [**Anthropic**](/providers/anthropic)
- [**Amazon Bedrock**](/providers/bedrock)
- [**Google Vertex AI (Anthropic)**](/providers/anthropic-vertexai)

We've tested with `OpenAI o3-mini` and `Anthropic Claude 3.7`, which we recommend for optimal performance.

## Recommended Models

For best performance, we recommend using:

- **OpenAI**: o3-mini
- **Anthropic**: Claude 3.7

These models provide a good balance of performance, speed, and cost for most use cases.

## Configuration Schema

Below is the configuration schema for all supported AI providers:

| Field              | Type    | Required | Description                                                                                               |
| ------------------ | ------- | -------- | --------------------------------------------------------------------------------------------------------- |
| `ai-provider`      | string  | No       | AI provider to use. Options: `openai`, `anthropic`, `bedrock`, `anthropic-vertexai`. Defaults to `openai` |
| `ai-endpoint`      | string  | No       | Custom OpenAI-compatible API endpoint URL                                                                 |
| `ai-api-key`       | string  | No       | AI API token for authentication                                                                           |
| `bedrock-region`   | string  | No       | AWS region for Amazon Bedrock                                                                             |
| `vertexai-region`  | string  | No       | Google Cloud region for Vertex AI                                                                         |
| `vertexai-project` | string  | No       | Google Cloud project ID for Vertex AI                                                                     |
| `ai-model`         | string  | No       | AI model to use (provider-specific)                                                                       |
| `ai-max-tokens`    | integer | No       | Maximum tokens to use in the response (0 = provider default)                                              |
| `ai-temperature`   | float   | No       | Temperature for AI responses (-1.0 = provider default)                                                    |
| `ai-reasoning`     | boolean | No       | Enable reasoning mode for supported models (default: true)                                                |

## Example

First specify `OPENAI_API_KEY` in the [environment](https://help.openai.com/en/articles/5112595-best-practices-for-api-key-safety). You can get OpenAI API Key on [OpenAI Platform](https://platform.openai.com/api-keys).

```bash
export OPENAI_API_KEY='yourkey'
```

```bash
./gateway discover \
  --ai-provider openai \
  --config connection.yaml
```

## Additional Configuration Options

You can further customize the AI behavior with these optional parameters:

```bash
./gateway discover \
  --ai-provider openai \
  --ai-api-key your-openai-api-key \
  --ai-model o3-mini \
  --ai-max-tokens 8192 \
  --ai-temperature 1.0 \
  --ai-reasoning true \
  --config connection.yaml
```
