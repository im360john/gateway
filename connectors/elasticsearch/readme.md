---
title: 'Elasticsearch'
---

Elasticsearch connector allows querying Elasticsearch clusters using their native REST API.

## Config Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | string | yes | constant: `elasticsearch` |
| hosts | string[] | yes | List of Elasticsearch nodes (e.g., ["http://localhost:9200"]) |
| user | string | no | Username for authentication |
| password | string | no | Password for authentication |
| api_key | string | no | API key for authentication (alternative to username/password) |
| cloud_id | string | no | Cloud ID for Elastic Cloud deployments |
| index | string | no | Default index to query |
| tls_verify | boolean | no | Verify TLS certificates (default: true) |
| conn_string | string | no | Direct connection string |

## Config example:

```yaml
connection:
    type: elasticsearch
    hosts:
    - http://localhost:9200
    - http://es-node2:9200
    user: elastic
    password: secret
    index: my_index
    tls_verify: true
```

Or as alternative with direct connection string:

```yaml
connection:
    type: elasticsearch
    conn_string: https://elastic:secret@localhost:9200
```

## Authentication Methods

The connector supports multiple authentication methods:

1. Basic Authentication (username/password)
2. API Key Authentication
3. No Authentication (for development)

For Elastic Cloud deployments, you can use the cloud_id parameter instead of hosts:

```yaml
connection:
    type: elasticsearch
    cloud_id: deployment:dXMtZWFzdC0xLmF3cy5mb3VuZC5pbyQ0ZmE...
    api_key: your_api_key
```

## Notes

- The connector uses the official Elasticsearch Go client
- For high availability, specify multiple hosts - the client will automatically handle failover
- TLS is required for production deployments
- When using API keys, both username and password should be omitted
- The index parameter is optional but recommended to set a default index for queries
- Supports Elasticsearch version 7.x and above
