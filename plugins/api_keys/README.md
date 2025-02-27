# API Keys Plugin

Provides API key authentication for endpoints.

## Type
- Wrapper
- Swaggerer

## Description
Implements API key authentication by validating keys from headers or query parameters. Supports method-level permissions per key.

## Configuration

```yaml
api_keys:
  name: "X-API-Key"          # Header or query parameter name
  location: "header"         # Where to look for the key: "header" or "query"
  keys:                      # List of valid API keys
    - key: "secret-key-1"
      allowed_methods:       # Allowed methods for this key
        - "get_users"
        - "create_user"
    - key: "admin-key"      # Key with all methods allowed
      allowed_methods: []    
  keys_file: "/path/to/keys.yaml"  # Optional external keys file
``` 