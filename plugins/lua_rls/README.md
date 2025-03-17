---
title: Lua RLS Plugin
---

Row-Level Security implementation using Lua scripts.

## Type
- Interceptor

## Description
Allows defining custom row-level security logic using Lua scripts, which are executed for each row in the result set.

## Configuration

```yaml
lua_rls:
  script: |
    function filter_rows(row, context)
      if context.user_role == "admin" then
        return true
      end
      return row.tenant_id == context.tenant_id
    end
  variables:           # Global variables available to Lua script
    max_rows: 1000
    debug: true
  cache_size: 100     # Size of the Lua VM cache
``` 

## Context Object Properties

The `context` parameter of the `filter_rows` function contains the following properties:

- **Authentication Claims**: All JWT claims from the authenticated user are available directly as properties (e.g., `context.user_role`, `context.tenant_id`, `context.email`, etc.).
- **Request Headers**: All HTTP request headers are available as properties (with the same case as they appear in the request).
- **Custom Variables**: Any custom variables defined in the `variables` configuration section are available as global variables.

### Example Context Properties

```lua
-- Authentication claims from JWT or OAuth providers
context.user_id       -- User's unique identifier
context.user_role     -- User's role (e.g., "admin", "user")
context.tenant_id     -- Tenant/organization identifier
context.email         -- User's email address
context.groups        -- User's groups or permissions
context.org_id        -- Organization identifier

-- Request headers (same case as in HTTP request)
context.Authorization -- Authorization header
context.X-Tenant-ID   -- Custom tenant header
context.User-Agent    -- User agent header
