---
title: 'PostgreSQL'
---

PostgreSQL connector allows querying PostgreSQL databases.

## Config Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| hosts | string[] | yes | List of database hosts |
| database | string | yes | Database name |
| user | string | yes | Username |
| password | string | yes | Password |
| port | integer | yes | TCP port (default 5432) |
| tlsFile | string | no | PEM-encoded certificate for TLS |
| enableTLS | boolean | no | Enable TLS/SSL connection |

## Config example:

```yaml
hosts: 
  - localhost
database: mydb
user: postgres
password: secret
port: 5432
tlsFile: ""        # Optional PEM-encoded certificate
enableTLS: false   # Enable TLS/SSL connection 
```
