---
title: 'ChatGPT Integration'
---

# Overview

This guide explains how to integrate your Gateway-generated API with ChatGPT by creating a custom GPT with actions.

## Prerequisites

Before integrating with ChatGPT, ensure you have:

1. A running Gateway API with a publicly accessible endpoint
2. OpenAPI/Swagger documentation available at your API endpoint (e.g., `https://your-api.com/swagger/`) 
3. A ChatGPT Plus subscription to access the GPT Builder
4. Necessary authentication information if your API requires it

## Creating a Custom GPT with Gateway API Actions

### Step 1: Access GPT Builder

1. Log in to your ChatGPT Plus account
2. Navigate to GPT Builder at https://chatgpt.com/gpts/mine

![img](/assets/mygpts.jpg)

### Step 2: Set Up Your Custom GPT

1. Click "Create new GPT"
2. Fill in the basic information:
   - Name: Choose a name related to your API's purpose
   - Description: Explain what your API does and how users can interact with it
   - Instructions: Provide detailed guidance on how the GPT should use your API

[IMAGE PLACEHOLDER: Screenshot of the "Create new GPT" form with fields filled in]

### Step 3: Add Actions to Connect to Your Gateway API

1. In the GPT Builder interface, navigate to the "Actions" tab
2. Click "Create new action"

[IMAGE PLACEHOLDER: Screenshot of the Actions tab with "Create new action" button]

3. Enter a name for your action that describes what it does
4. For "Authentication," select the appropriate method:
   - "None" if your API is publicly accessible
   - "API Key" if your API uses key-based authentication
   - "OAuth" for OAuth authentication

[IMAGE PLACEHOLDER: Screenshot of action creation form with authentication options]

### Step 4: Configure OpenAPI Specification

1. Select "OpenAPI schema" for Schema type
2. Enter your Gateway API's OpenAPI URL (e.g., `https://your-gateway-api.com/openapi.json` or `https://your-gateway-api.com/swagger/`)
   
   Alternatively, you can:
   - Click "Import from URL" and enter your OpenAPI URL
   - Copy and paste the OpenAPI JSON directly into the schema field

[IMAGE PLACEHOLDER: Screenshot of entering the OpenAPI URL]

### Step 5: Customize Privacy and Settings

1. Review the endpoints that will be accessible to the GPT
2. Configure which operations your GPT can perform
3. If required, set up authentication details:
   - For API Key authentication, specify the header name (e.g., `X-API-KEY`)
   - Enter placeholder values that users will need to provide

[IMAGE PLACEHOLDER: Screenshot of privacy and authentication settings]

### Step 6: Save and Test Your Custom GPT

1. Click "Save" to create your action
2. Test your API connection with example queries
3. Make adjustments to the instructions as needed to optimize the interaction

[IMAGE PLACEHOLDER: Screenshot of testing interface with example query]

## Example Conversation Flow

Here's an example of how a conversation with your Gateway API-enabled GPT might look:

**User**: "Can you show me data about product sales for Q2?"

**ChatGPT**: "I'll retrieve that information for you. Let me access the sales data API..."
[ChatGPT uses your Gateway API to fetch the data]
"Here are the Q2 product sales:
- Product A: $125,000
- Product B: $93,500
- Product C: $78,200

Would you like to see any specific product's details or compare with previous quarters?"

## Troubleshooting

If you encounter issues with your ChatGPT integration:

1. **API not accessible**: Ensure your API is publicly accessible from the internet
2. **Authentication errors**: Verify your authentication details are correctly configured
3. **Schema errors**: Check that your OpenAPI specification is valid and properly formatted
4. **Rate limiting**: Consider implementing rate limiting on your API to prevent abuse

## Publishing Your Custom GPT (Optional)

If you want to share your GPT with others:

1. Go to the "Configure" tab and select visibility settings
2. Add appropriate usage instructions and examples
3. Click "Publish" to make your GPT available according to your selected visibility

[IMAGE PLACEHOLDER: Screenshot of the publishing interface with visibility options]

For more details on creating custom GPTs, refer to the [OpenAI GPT Builder Documentation](https://platform.openai.com/docs/guides/gpt-builder).