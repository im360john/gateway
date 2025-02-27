# Lua RLS Plugin

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