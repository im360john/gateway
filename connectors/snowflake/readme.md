Snowflake connector allows querying Snowflake data warehouse.

Config example:
    account: myaccount    # Your Snowflake account identifier
    database: MYDB
    user: myuser
    password: secret
    warehouse: COMPUTE_WH # Warehouse to use for queries
    schema: PUBLIC        # Schema name
    role: ACCOUNTADMIN    # Role to assume 