# Product Service - E-commerce Microservice Boilerplate

A production-ready Go microservice boilerplate for handling product operations in an e-commerce system. This service demonstrates Clean Architecture principles with full integration of PostgreSQL, Redis, Kafka, and Elasticsearch.

## Architecture

This service follows **Clean Architecture** (Hexagonal Architecture) with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│         Transport Layer (Gin)           │
│         (HTTP Handlers)                 │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│      Business Logic Layer (Service)      │
│      (Domain Logic & Orchestration)      │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│    Infrastructure Layer (Repository)    │
│    (PostgreSQL, Redis, Kafka, ES)       │
└─────────────────────────────────────────┘
```

### Key Principles

1. **Dependency Inversion**: High-level modules (service) don't depend on low-level modules (repositories). Both depend on abstractions (interfaces).
2. **Separation of Concerns**: Each layer has a single responsibility.
3. **Testability**: Interfaces allow easy mocking for unit tests.
4. **Independence**: Business logic is independent of frameworks and databases.

## Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin Gonic
- **Database**: PostgreSQL (GORM)
- **Cache**: Redis (go-redis/v9)
- **Message Queue**: Kafka (segmentio/kafka-go)
- **Search Engine**: Elasticsearch (official go-client)
- **Configuration**: Viper
- **Logging**: Uber Zap

## Project Structure

```
product-service/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   ├── config.yaml            # Configuration file
│   └── config.go              # Config structs and loader
├── internal/
│   ├── domain/                # Domain entities and interfaces
│   │   ├── product.go
│   │   └── event.go
│   ├── repository/            # Infrastructure implementations
│   │   ├── postgres/
│   │   ├── redis/
│   │   ├── kafka/
│   │   └── elasticsearch/
│   ├── service/               # Business logic
│   │   └── product_service.go
│   ├── handler/               # HTTP handlers (Gin)
│   │   └── product_handler.go
│   └── router/                # Route definitions
│       └── router.go
├── pkg/                       # Shared packages
│   ├── database/              # DB connection singleton
│   ├── redis/                 # Redis client singleton
│   ├── elasticsearch/          # ES client singleton
│   └── logger/                # Logger setup
├── deploy/                    # Deployment files
├── docker-compose.yml         # Local development stack
├── Dockerfile                 # Multi-stage build
└── go.mod                     # Go dependencies
```

## Features

### 1. Product Management
- Create products with validation
- Update products
- Get product by ID (with cache-first strategy)
- Search products using Elasticsearch

### 2. Caching Strategy
- **Cache-Aside Pattern**: Check cache first, fallback to database
- Automatic cache population on cache miss
- Cache invalidation on updates

### 3. Distributed Locking
- Redis-based distributed locks for inventory updates
- Prevents race conditions in concurrent scenarios

### 4. Event-Driven Architecture
- Kafka integration for publishing product events
- Events: `product_created`, `product_updated`
- Enables microservice communication

### 5. Full-Text Search
- Elasticsearch integration for product search
- Supports filtering by category, price range, etc.

### 6. Graceful Shutdown
- Proper cleanup of all connections (DB, Redis, Kafka)
- Context-based timeout for shutdown operations

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Make (optional, for convenience)

### Local Development

1. **Start infrastructure services**:
   ```bash
   docker-compose up -d
   ```

   This starts:
   - PostgreSQL on port 5432
   - Redis on port 6379
   - Zookeeper on port 2181
   - Kafka on port 9092
   - Elasticsearch on port 9200

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Run the service**:
   ```bash
   go run cmd/main.go
   ```

   The service will start on `http://localhost:8080`

### Configuration

Configuration is managed via `config/config.yaml` and environment variables. Environment variables take precedence.

Example environment variables:
```bash
export SERVER_PORT=8080
export DATABASE_HOST=localhost
export DATABASE_PASSWORD=postgres
export REDIS_HOST=localhost
export KAFKA_BROKERS=localhost:9092
export ELASTICSEARCH_ADDRESSES=http://localhost:9200
```

## API Endpoints

### Health Check
```
GET /health
```

### Create Product
```
POST /api/v1/products
Content-Type: application/json

{
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "sku": "LAP-001",
  "category": "Electronics",
  "stock": 10,
  "is_active": true
}
```

### Get Product
```
GET /api/v1/products/:id
```

### Update Product
```
PUT /api/v1/products/:id
Content-Type: application/json

{
  "name": "Updated Laptop",
  "price": 899.99
}
```

### Search Products
```
GET /api/v1/products/search?q=laptop&category=Electronics
```

### Update Inventory
```
PATCH /api/v1/products/:id/inventory
Content-Type: application/json

{
  "quantity": -2
}
```

## Docker Build

Build the Docker image:
```bash
docker build -t product-service:latest .
```

Run the container:
```bash
docker run -p 8080:8080 \
  -e DATABASE_HOST=host.docker.internal \
  -e REDIS_HOST=host.docker.internal \
  product-service:latest
```

## Architecture Decisions

### Why Clean Architecture?

1. **Testability**: Business logic can be tested without databases or HTTP frameworks
2. **Flexibility**: Easy to swap implementations (e.g., PostgreSQL → MongoDB)
3. **Maintainability**: Clear boundaries make code easier to understand and modify
4. **Independence**: Business rules don't depend on external frameworks

### Why Singleton Pattern for Connections?

- **Resource Efficiency**: Single connection pool per service instance
- **Consistency**: All parts of the application use the same connection
- **Thread Safety**: Go's `sync.Once` ensures safe concurrent access

### Why Async Operations for Cache/ES/Kafka?

- **Performance**: Don't block HTTP response on non-critical operations
- **Resilience**: Service remains available even if cache/ES/Kafka is down
- **Eventual Consistency**: Acceptable for search and caching use cases

### Why Manual Dependency Injection?

- **Clarity**: Explicit dependencies make code easier to understand
- **No Magic**: No code generation or reflection overhead
- **Control**: Full control over initialization order and error handling

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/service/...
```

## Production Considerations

1. **Environment Variables**: Use secrets management (e.g., Kubernetes Secrets, AWS Secrets Manager)
2. **Monitoring**: Add Prometheus metrics and distributed tracing
3. **Rate Limiting**: Implement rate limiting middleware
4. **Authentication**: Add JWT authentication middleware
5. **Database Migrations**: Use a migration tool (e.g., golang-migrate)
6. **Error Handling**: Implement structured error responses
7. **Circuit Breakers**: Add resilience patterns for external services
8. **Health Checks**: Enhance health endpoint with dependency checks

## License

MIT

