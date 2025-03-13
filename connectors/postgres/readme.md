---
title: 'PostgreSQL'
---

PostgreSQL connector allows querying PostgreSQL databases.

## Config Schema

| Field       | Type     | Required | Description                            |
|-------------|----------|----------|----------------------------------------|
| type        | string   | yes | constant: `postgres`                   |
| hosts       | string[] | yes | List of database hosts                 |
| database    | string   | yes | Database name                          |
| user        | string   | yes | Username                               |
| password    | string   | yes | Password                               |
| port        | integer  | yes | TCP port (default 5432)                |
| schema      | string   | no | instead of "public" it could be custom |
| tls_file    | string   | no | PEM-encoded certificate for TLS        |
| enable_tls  | boolean  | no | Enable TLS/SSL connection              |
| conn_string | string   | no | DSN-like connection string             |

## Config example:

```yaml
type: postgres
hosts: 
  - localhost
database: mydb
user: postgres
password: secret
port: 5432
schema: sales
tls_file: ""        # Optional PEM-encoded certificate
enable_tls: false   # Enable TLS/SSL connection 
```

Or as alternative with plan DSN:

```yaml
type: postgres
conn_string: postgresql://my_user:my_pass@localhost:5432/mydb
```

## Deploy PostgreSQL Sample Databases 

 - You can deploy sample database using a docker and example of DVD store data
<a href="/example/postgresql-dvdstore-sample" /> DVDStore Sample Database</a>
