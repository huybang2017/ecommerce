# Quick Start Guide

This guide will help you get the Product Service up and running in minutes.

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Make (optional but recommended)

## Step 1: Start Infrastructure Services

Start all required services (PostgreSQL, Redis, Kafka, Elasticsearch) using Docker Compose:

```bash
docker-compose up -d
```

Wait for all services to be healthy (about 30 seconds). Check status:

```bash
docker-compose ps
```

## Step 2: Install Dependencies

```bash
go mod download
```

## Step 3: Configure (Optional)

Copy the example environment file and modify if needed:

```bash
cp .env.example .env
```

Edit `.env` if you need to change any default values. The service will use `config/config.yaml` by default, and environment variables will override those values.

## Step 4: Run the Service

### Option A: Using Make

```bash
make run
```

### Option B: Direct Go Command

```bash
go run cmd/main.go
```

The service will start on `http://localhost:8080`

## Step 5: Test the API

### Health Check

```bash
curl http://localhost:8080/health
```

### Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "sku": "LAP-001",
    "category": "Electronics",
    "stock": 10,
    "is_active": true
  }'
```

### Get Product

```bash
curl http://localhost:8080/api/v1/products/1
```

### Search Products

```bash
curl "http://localhost:8080/api/v1/products/search?q=laptop&category=Electronics"
```

## Verify Integration

### Check PostgreSQL

```bash
docker exec -it product-service-postgres psql -U postgres -d product_service -c "SELECT * FROM products;"
```

### Check Redis Cache

```bash
docker exec -it product-service-redis redis-cli KEYS "product:*"
```

### Check Elasticsearch Index

```bash
curl http://localhost:9200/products/_search?pretty
```

### Check Kafka Topic

```bash
docker exec -it product-service-kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic product_updated \
  --from-beginning
```

## Stopping Services

```bash
docker-compose down
```

## Troubleshooting

### Services not starting

Check logs:
```bash
docker-compose logs
```

### Port conflicts

If ports are already in use, modify `docker-compose.yml` to use different ports.

### Connection errors

Ensure all services are healthy:
```bash
docker-compose ps
```

All services should show "healthy" status.

## Next Steps

- Read the [README.md](README.md) for architecture details
- Review the code structure in `internal/` directory
- Customize the configuration in `config/config.yaml`
- Add authentication and authorization
- Implement additional endpoints
- Add unit and integration tests

