---
title: 'Snowflake'
---

Snowflake connector allows querying Snowflake data warehouse.

## Config Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| account | string | yes | Snowflake account identifier |
| database | string | yes | Database name |
| user | string | yes | Username |
| password | string | yes | Password |
| warehouse | string | yes | Warehouse to use for queries |
| schema | string | yes | Schema name |
| role | string | yes | Role to assume |

## Config example:

```yaml
account: myaccount    # Your Snowflake account identifier
database: MYDB
user: myuser
password: secret
warehouse: COMPUTE_WH # Warehouse to use for queries
schema: PUBLIC        # Schema name
role: ACCOUNTADMIN    # Role to assume 
```

## Config Schema:

```yaml
type: object
properties:
  account:
    type: string
    required: true
    description: Snowflake account identifier
  database:
    type: string
    required: true
    description: Database name
  user:
    type: string
    required: true
    description: Username
  password:
    type: string
    required: true
    description: Password
  warehouse:
    type: string
    required: true
    description: Warehouse to use for queries
  schema:
    type: string
    required: true
    description: Schema name
  role:
    type: string
    required: true
    description: Role to assume
```
