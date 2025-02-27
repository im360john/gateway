ClickHouse connector allows querying ClickHouse databases.

Config example:
    host: localhost      # Single host address
    hosts:              # Or multiple hosts for cluster setup
      - host1.example.com
      - host2.example.com
    database: mydb
    user: default
    password: secret
    port: 8123
    secure: false       # Use HTTPS instead of HTTP 