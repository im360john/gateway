---
title: 'Complex E2E Example with Gateway API'
description: 'This example demonstrates a complete end-to-end setup of the Gateway API with various features including data access control, caching, and PII data handling.'
---

This example demonstrates a complete end-to-end setup of the Gateway API with various features including data access control, caching, and PII data handling.

## Components

- **Gateway API**: REST API service with customer order endpoints
- **PostgreSQL**: Main database storing order and customer data
- **Jaeger**: Distributed tracing system for monitoring
- **Presidio**: Service for PII data anonymization

## Features Demonstrated

- REST API endpoints for customer orders
- LRU caching with TTL
- API key authentication
- Row-level security using Lua scripts
- PII data removal
- OpenTelemetry integration with Jaeger
- Star schema data model for order analytics

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Start all services:
   ```bash
   docker-compose up -d
   ```

2. Services will be available at:
   - Gateway API: http://localhost:8182
   - Jaeger UI: http://localhost:16686
   - PostgreSQL: localhost:5432
   - Presidio: http://localhost:5001

## API Endpoints

### 1. Search Customer Orders
```
GET /customer/{customer_key}/orders
```
Query parameters:
- `start_date` (optional): Filter orders from this date
- `end_date` (optional): Filter orders until this date
- `min_total` (optional): Minimum order total
- `max_total` (optional): Maximum order total
- `limit` (optional, default: 50): Number of results
- `offset` (optional, default: 0): Pagination offset

### 2. Get Order Details
```
GET /customer/{customer_key}/order/{payment_key}
```

## Security Features

1. API Key Authentication:
   - Header: `x-api-key`
   - Available keys:
     - `all_methods`: Access to all endpoints
     - `only_orders`: Limited to order-related endpoints

2. Row-Level Security:
   - Validates user access using `X-User-ID` header
   - Ensures users can only access their own data

3. PII Protection:
   - Automatically removes sensitive information from `address` fields

## Database Schema

The database follows a star schema design with:
- Fact table: `fact_table` (order transactions)
- Dimension tables:
  - `payment_dim`: Payment information
  - `customer_dim`: Customer details
  - `item_dim`: Product information
  - `store_dim`: Store locations
  - `time_dim`: Time-based dimensions

## Configuration Files

- `docker-compose.yml`: Service orchestration
- `gateway.yaml`: Gateway API configuration
- `connection.yaml`: Database connection settings
- `init.sql`: Database initialization script

## Data Loading

Sample data is automatically loaded during initialization from CSV files mounted in the PostgreSQL container. 
