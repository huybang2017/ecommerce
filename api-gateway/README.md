# API Gateway - E-commerce Microservices Gateway

A production-ready API Gateway for routing requests to backend microservices. This gateway implements authentication, rate limiting, request logging, and service discovery.

## Architecture

The API Gateway follows **Clean Architecture** principles and acts as the single entry point for all client requests:

```
┌─────────────────────────────────────────┐
│         Client Applications             │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│         API Gateway (Port 8000)         │
│  - Authentication (JWT)                  │
│  - Rate Limiting                        │
│  - Request Routing                      │
│  - Request Logging                      │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
┌───────▼──────┐   ┌────────▼────────┐
│ Product      │   │ Other Services   │
│ Service      │   │ (Future)         │
│ (Port 8080)  │   │                  │
└──────────────┘   └──────────────────┘
```

## Features

### 1. Request Routing
- Routes requests to appropriate microservices based on path patterns
- Configurable service endpoints
- Health check monitoring for all services

### 2. Authentication & Authorization
- JWT token validation
- Bearer token authentication
- Protected and public routes
- Optional authentication for certain endpoints

### 3. Rate Limiting
- Per-IP rate limiting
- Configurable requests per minute
- Burst capacity support
- Prevents API abuse

### 4. CORS Support
- Configurable allowed origins
- Preflight request handling
- Credential support

### 5. Request Logging
- Structured logging with Zap
- Request/response logging
- Error tracking
- Performance metrics (latency)

### 6. Health Checks
- Gateway health endpoint
- Service health aggregation
- Service discovery support

## Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin Gonic
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Rate Limiting**: golang.org/x/time
- **Configuration**: Viper
- **Logging**: Uber Zap

## Project Structure

```
api-gateway/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   ├── config.yaml            # Configuration file
│   └── config.go              # Config structs and loader
├── internal/
│   ├── domain/                # Domain entities and interfaces
│   │   └── service.go         # Service registry interfaces
│   ├── repository/            # Infrastructure implementations
│   │   ├── service_registry.go
│   │   └── proxy_client.go
│   ├── service/               # Business logic
│   │   └── gateway_service.go
│   ├── handler/               # HTTP handlers
│   │   └── gateway_handler.go
│   ├── middleware/            # Middleware
│   │   ├── auth.go            # JWT authentication
│   │   ├── rate_limit.go      # Rate limiting
│   │   └── logging.go        # Request logging
│   └── router/                # Route definitions
│       └── router.go
├── pkg/
│   └── logger/                # Logger setup
├── docker-compose.yml         # Local development
├── Dockerfile                  # Multi-stage build
└── go.mod                      # Go dependencies
```

## Configuration

Configuration is managed via `config/config.yaml` and environment variables:

```yaml
server:
  port: 8000
  mode: "debug"

jwt:
  secret: "your-secret-key-change-in-production"
  expiration: 24h

rate_limit:
  enabled: true
  requests_per_minute: 100
  burst: 20

services:
  product_service:
    base_url: "http://localhost:8080"
    timeout: 30s
```

## API Endpoints

### Gateway Endpoints

- `GET /health` - Gateway health check
- `GET /api/gateway/health` - Gateway health with service status

### Proxied Endpoints (Product Service)

All requests to `/api/v1/products/*` are proxied to the Product Service:

- `GET /api/v1/products` - List/search products (public)
- `GET /api/v1/products/:id` - Get product by ID (public)
- `GET /api/v1/products/search` - Search products (public)
- `POST /api/v1/products` - Create product (public)
- `PUT /api/v1/products/:id` - Update product (protected)
- `PATCH /api/v1/products/:id` - Partial update (protected)
- `PATCH /api/v1/products/:id/inventory` - Update inventory (protected)
- `DELETE /api/v1/products/:id` - Delete product (protected)

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Product Service running (or use docker-compose)

### Local Development

1. **Start infrastructure and services**:
   ```bash
   # From project root
   docker-compose up -d
   ```

2. **Install dependencies**:
   ```bash
   cd api-gateway
   go mod download
   ```

3. **Run the gateway**:
   ```bash
   go run cmd/main.go
   ```

   The gateway will start on `http://localhost:8000`

### Testing the Gateway

```bash
# Health check
curl http://localhost:8000/health

# Get products (proxied to product-service)
curl http://localhost:8000/api/v1/products

# Create product (proxied to product-service)
curl -X POST http://localhost:8000/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 999.99,
    "sku": "LAP-001",
    "category": "Electronics",
    "stock": 10
  }'
```

## Authentication

Protected routes require a JWT token in the Authorization header:

```bash
curl -X PUT http://localhost:8000/api/v1/products/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Product"}'
```

## Rate Limiting

Rate limiting is enabled by default:
- 100 requests per minute per IP
- Burst capacity: 20 requests

When rate limit is exceeded, the gateway returns:
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```
Status code: `429 Too Many Requests`

## Docker Build

```bash
# Build the gateway
docker build -t api-gateway:latest .

# Run the gateway
docker run -p 8000:8000 \
  -e SERVICES_PRODUCT_SERVICE_BASE_URL=http://product-service:8080 \
  api-gateway:latest
```

## Production Considerations

1. **JWT Secret**: Change the default JWT secret in production
2. **Service Discovery**: Implement Consul, Eureka, or Kubernetes service discovery
3. **Load Balancing**: Add load balancing for multiple service instances
4. **Circuit Breaker**: Implement circuit breaker pattern for resilience
5. **Metrics**: Add Prometheus metrics and Grafana dashboards
6. **TLS/HTTPS**: Enable HTTPS with TLS certificates
7. **Request ID**: Add request ID tracking for distributed tracing
8. **API Versioning**: Implement proper API versioning strategy

## Architecture Decisions

### Why In-Memory Service Registry?

- Simple and fast for small to medium deployments
- Easy to extend with external service discovery
- No external dependencies required

### Why Per-IP Rate Limiting?

- Prevents abuse from individual clients
- Simple to implement and understand
- Can be extended with user-based rate limiting

### Why JWT Authentication?

- Stateless authentication
- Scalable across multiple gateway instances
- Industry standard

## License

MIT

