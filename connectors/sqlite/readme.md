---
title: 'SQLite'
---
This connector allows you to connect to SQLite databases. SQLite is a self-contained, serverless, zero-configuration, transactional SQL database engine.

## Configuration

The connector supports the following configuration options:

```yaml
# Simple connection string format
sqlite: "path/to/database.db"

# Or detailed configuration
sqlite:
  hosts: ["/path/to/directory"]  # Directory containing the database file
  database: "database.db"        # Database file name
  read_only: false              # Whether to open database in read-only mode
  memory: false                 # Whether to create an in-memory database
```

### Configuration Options

- `hosts`: List of database file paths. If specified with `database`, the first host is used as the directory path.
- `database`: Database file name. If specified with `hosts`, it's combined with the first host path.
- `read_only`: Whether to open the database in read-only mode. Defaults to false.
- `memory`: Whether to create an in-memory database. If true, ignores hosts and database settings.
- `conn_string`: Direct connection string. If provided, all other options are ignored.

### Examples

1. Using a file-based database:
```yaml
sqlite:
  hosts: ["/data"]
  database: "mydb.db"
```

2. Using an in-memory database:
```yaml
sqlite:
  memory: true
```

3. Using a direct connection string:
```yaml
sqlite: "/absolute/path/to/database.db"
```

## Features

- Supports both file-based and in-memory SQLite databases
- Read-only mode support
- Named parameter support in queries (using `:param` syntax)
- Automatic table discovery
- Column type inference
- Transaction support 