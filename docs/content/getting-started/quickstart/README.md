
This guide will help you get started with Gateway using Docker, discover your database and launch API on top of it.

## Prerequisites

- <a href="https://docs.docker.com/get-started/get-docker/">Docker</a> installed on your system
- <a href="https://platform.openai.com/api-keys">OpenAI API key</a> for AI-powered API generation
- Your PostgreSQL Database or any other db that gateway supports. You can also take our example databases:
    - <a href="/example/postgresql-dvdstore-sample/">PostgreSQL DVD Store Sample Database</a>. 
    - <a href="/example/postgresql-ecommerce-sample/">PostgreSQL Ecommerce Sample Database</a>. 

## Installation Steps

1. Pull the latest Gateway Docker image:
```bash
docker pull ghcr.io/centralmind/gateway:latest
```

2. Create a `connection.yaml` configuration file:
```yaml
hosts:
  - localhost
user: "your-database-user"
password: "your-database-password"
database: "your-database-name"
port: 5432
```

3. Run the discovery process with AI-powered API generation:
```bash
docker run -v $(pwd)/connection.yaml:/app/connection.yaml \
  ghcr.io/centralmind/gateway:latest discover \
  --config connection.yaml \
  --db-type postgres \  
  --ai-api-key $OPENAI_KEY \
  --prompt "Develop an API that enables a chatbot to retrieve information about data. \
Try to place yourself as analyst and think what kind of data you will require, \
based on that come up with useful API methods for that"
```

4. Start the REST server:
```bash
docker run -p 9090:9090 -v \
  -v $(pwd)/gateway.yaml:/app/gateway.yaml \
  ghcr.io/centralmind/gateway:latest start --config gateway.yaml rest
```

## Verification

After starting the REST server, you can verify the installation by accessing the Swagger UI:
```
http://localhost:9090/swagger/
```

