---
title: 'DuckDB'
---

DuckDB connector allows querying DuckDB databases, which is an embedded analytical database similar to SQLite but optimized for OLAP workloads.

## Config Schema

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | string | yes | constant: `duckdb` |
| hosts | string[] | no* | List of paths (only first path is used) |
| database | string | no* | Database file name, will be opened in readonly mode |
| init_sql | string | no | SQL commands to execute on connection initialization (e.g. installing extensions, attaching databases) |
| memory | boolean | no | If true, uses an in-memory database |
| conn_string | string | no | Direct connection string, overrides other parameters |


## Config Examples

1. Using directory path in hosts with initialization SQL:
```yaml
connection:
  type: duckdb
  hosts:
    - ./data    # relative path to directory
  database: analytics.duckdb
  init_sql: |
    FORCE INSTALL aws FROM core_nightly;
    FORCE INSTALL httpfs FROM core_nightly;
    FORCE INSTALL iceberg FROM core_nightly;
    CREATE TABLE weather AS
            SELECT * FROM read_csv_auto('https://raw.githubusercontent.com/duckdb/duckdb-web/main/data/weather.csv');
```

2. Using full file path in hosts:
```yaml
connection:
  type: duckdb
  hosts:
    - /absolute/path/to/analytics.duckdb    # Unix-style path
    # or
    - C:/Users/MyUser/data/analytics.duckdb  # Windows-style path
```

3. Using relative file path:
```yaml
connection:
  type: duckdb
  hosts:
    - ./data/analytics.duckdb
```

4. Using current directory:
```yaml
connection:
  type: duckdb
  hosts:
    - .
  database: analytics.duckdb
```

5. Using in-memory mode (recommended format):
```yaml
connection:
  type: duckdb
  memory: true
```

6. Using in-memory mode with direct connection string:
```yaml
connection:
  type: duckdb
  conn_string: ":memory:"
```

7. Using empty connection section (defaults to in-memory):
```yaml
connection:
```

## Running Discovery and API
You can also pass connection string as parameter:

### File-based connection strings
Using absolute path on Linux:
```
./gateway discover --ai-provider gemin --connection-string "duckdb:///absolute/path/to/duckdb-demo.duckdb"
```
or on Windows
```
.\gateway discover --ai-provider gemini --connection-string "duckdb://C:/path/duckdb-demo.duckdb"
```

### In-memory connection strings
```
./gateway discover --ai-provider openai --connection-string "duckdb://:memory:"
```

Start server, it will use `gateway.yaml` generated from prev step:
```
./gateway start
```

## Path Resolution

The final database path is determined as follows:
1. If `conn_string` is provided: uses it directly
2. If `memory` is true: uses in-memory database (`:memory:`)
3. If `hosts[0]` and `database` are provided: `hosts[0]/database`
4. If only `hosts[0]` is provided: uses it as the complete path
5. If only `database` is provided: uses it as a local path

## Safety Features

For security reasons, the connector automatically adds the following safety guard rails to connection strings:

1. For all file-based databases (non-memory), the `access_mode=READ_ONLY` parameter is applied to prevent write operations
2. For all database connections, `allow_community_extensions=false` is added to prevent loading potentially unsafe extensions
3. These parameters are automatically added as query parameters (after `?` or `&` as appropriate) to the connection string

Memory databases (:memory: or memory=true) do not have the READ_ONLY restriction, but still have community extensions disabled.

## Notes

- DuckDB is an embedded database, so no server setup is required
- Only the first path in `hosts` is used (others are ignored)
- Both forward slashes `/` and backslashes `\` are supported for Windows paths
- Relative paths are resolved relative to the current working directory
- File-based databases are opened in read-only mode by default
- For in-memory databases, use the `memory: true` flag or `conn_string: ":memory:"`
- In-memory databases still create temporary files for persistence, which is normal behavior
- The `init_sql` field allows executing multiple SQL commands on connection initialization, separated by semicolons