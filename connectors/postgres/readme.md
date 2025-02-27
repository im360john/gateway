PostgreSQL connector allows querying PostgreSQL databases.

Config example:
    hosts: 
      - localhost
    database: mydb
    user: postgres
    password: secret
    port: 5432
    tlsFile: ""        # Optional PEM-encoded certificate
    enableTLS: false   # Enable TLS/SSL connection 