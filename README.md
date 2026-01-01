# 🛒 E-Commerce Microservices Platform

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-15.1-black?logo=next.js)](https://nextjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?logo=postgresql)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)](https://redis.io/)

A production-ready, scalable e-commerce platform built with microservices architecture, inspired by Shopee's system design. This project demonstrates modern software engineering practices including Clean Architecture, Domain-Driven Design, and Event-Driven Architecture.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [Services Overview](#-services-overview)
- [Development](#-development)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### Current Features (Gate 1 Completed)

- ✅ **Authentication & Authorization**
  - Session-based authentication with Redis
  - JWT tokens (Access + Refresh)
  - HttpOnly cookies for security
  - Role-based access control (RBAC)
  - Device tracking and multi-session management

### Planned Features

- 🔄 Product catalog with variants and SKU management
- 🔍 Full-text search with Elasticsearch
- 🛒 Shopping cart and wishlist
- 💳 Order processing and checkout
- 📦 Inventory management with reservation system
- 💰 Flash sale and promotion engine
- 💳 Payment gateway integration (Mock)
- 📧 Event-driven notifications
- 📊 Admin dashboard for management
- 📈 Analytics and reporting

## 🏗️ Architecture

This project follows **Microservices Architecture** with clear separation of concerns:

```
┌─────────────┐     ┌──────────────┐
│   Client    │────▶│  API Gateway │
│  (Next.js)  │     │   Port 8000  │
└─────────────┘     └───────┬──────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Identity    │    │   Product    │    │    Order     │
│  Service     │    │   Service    │    │   Service    │
│  Port 8081   │    │  Port 8082   │    │  Port 8083   │
└──────┬───────┘    └──────┬───────┘    └──────┬───────┘
       │                   │                   │
       └───────────────────┴───────────────────┘
                           │
       ┌───────────────────┼───────────────────┐
       ▼                   ▼                   ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  PostgreSQL  │    │    Redis     │    │ Elasticsearch│
│  Port 5432   │    │  Port 6379   │    │  Port 9200   │
└──────────────┘    └──────────────┘    └──────────────┘
```

### Architectural Principles

- **Clean Architecture**: Each service follows Clean Architecture with clear layer separation
- **Domain-Driven Design**: Business logic organized around domain models
- **CQRS Pattern**: Separation of read and write operations where needed
- **Event-Driven**: Asynchronous communication via Kafka for decoupling
- **API Gateway Pattern**: Single entry point for all client requests

## 🛠️ Tech Stack

### Backend

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP), gRPC (Service-to-Service)
- **ORM**: GORM
- **Validation**: go-playground/validator
- **Documentation**: Swagger/OpenAPI

### Frontend

- **Framework**: Next.js 15.1
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Query, Context API
- **HTTP Client**: Axios

### Databases

- **Primary DB**: PostgreSQL 16
- **Cache**: Redis 7
- **Search**: Elasticsearch 8

### Infrastructure

- **Message Queue**: Apache Kafka
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes (planned)
- **CI/CD**: GitHub Actions (planned)

### Monitoring & Observability

- **Logging**: Zap (structured logging)
- **Metrics**: Prometheus (planned)
- **Tracing**: Jaeger (planned)

## 📁 Project Structure

```
ecommerce/
├── api-gateway/           # API Gateway Service
│   ├── cmd/              # Application entry point
│   ├── config/           # Configuration management
│   ├── internal/         # Private application code
│   │   ├── handler/     # HTTP handlers
│   │   ├── middleware/  # Custom middleware (CORS, Auth)
│   │   └── router/      # Route definitions
│   └── docs/            # Swagger documentation
│
├── identity-service/      # Identity & Auth Service
│   ├── cmd/
│   ├── config/
│   ├── internal/
│   │   ├── domain/      # Domain models & interfaces
│   │   ├── handler/     # HTTP handlers
│   │   ├── repository/  # Data access layer
│   │   │   ├── postgres/
│   │   │   └── redis/
│   │   ├── service/     # Business logic
│   │   └── middleware/  # Service middleware
│   └── pkg/             # Public libraries
│       ├── database/
│       ├── logger/
│       └── redis/
│
├── product-service/       # Product Catalog Service
├── search-service/        # Search Service (Elasticsearch)
├── order-service/         # Order Management Service
├── inventory-service/     # Inventory & Stock Service
├── promotion-service/     # Promotion & Flash Sale Service
├── payment-service/       # Payment Processing Service
├── notification-service/  # Notification Service
│
├── client/               # Customer Frontend (Next.js)
│   ├── src/
│   │   ├── app/         # App Router (Next.js 13+)
│   │   ├── components/  # React components
│   │   ├── contexts/    # React contexts
│   │   ├── hooks/       # Custom hooks
│   │   └── lib/         # Utilities & API client
│   └── public/
│
├── admin/                # Admin Dashboard (Next.js)
│
├── docs/                 # Project documentation
│   ├── BLUEPRINT.md     # System architecture
│   ├── CONTRACT.md      # API contracts
│   ├── GATES.md         # Development phases
│   └── INTAKE.md        # Project requirements
│
├── scripts/              # Utility scripts
│   └── init-databases.sql
│
├── docker-compose.yml    # Local development setup
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18+ and npm/yarn
- **Docker** and Docker Compose
- **Make** (optional, for convenience)
- **Git**

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/huybang2017/ecommerce.git
   cd ecommerce
   ```

2. **Start infrastructure services**

   ```bash
   docker-compose up -d postgres redis elasticsearch kafka zookeeper
   ```

3. **Initialize databases**

   ```bash
   # Wait for PostgreSQL to be ready
   docker-compose exec postgres psql -U postgres -f /scripts/init-databases.sql
   ```

4. **Start backend services**

   ```bash
   # Identity Service
   cd identity-service
   go mod download
   go run cmd/main.go

   # API Gateway (in new terminal)
   cd api-gateway
   go mod download
   go run cmd/main.go
   ```

5. **Start frontend**

   ```bash
   cd client
   npm install
   npm run dev
   ```

6. **Access the application**
   - Frontend: http://localhost:3000
   - API Gateway: http://localhost:8000
   - Identity Service: http://localhost:8081
   - Swagger UI: http://localhost:8000/swagger/index.html

## 🔧 Services Overview

| Service              | Port | Status     | Description                                |
| -------------------- | ---- | ---------- | ------------------------------------------ |
| API Gateway          | 8000 | ✅ Live    | Single entry point, routing, CORS handling |
| Identity Service     | 8081 | ✅ Live    | Authentication, user management, sessions  |
| Product Service      | 8082 | 🔄 WIP     | Product catalog, categories, variants      |
| Search Service       | 8083 | 📋 Planned | Full-text search, filters, suggestions     |
| Order Service        | 8084 | 📋 Planned | Cart, checkout, order processing           |
| Inventory Service    | 8085 | 📋 Planned | Stock management, reservations             |
| Promotion Service    | 8086 | 📋 Planned | Discounts, flash sales, coupons            |
| Payment Service      | 8087 | 📋 Planned | Payment processing (mock)                  |
| Notification Service | 8088 | 📋 Planned | Email, SMS, push notifications             |

## 💻 Development

### Environment Variables

Each service uses environment variables for configuration. Create `.env` files:

**Identity Service** (`identity-service/.env`):

```env
SERVER_PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_USER=identity_user
DB_PASSWORD=identity_pass
DB_NAME=identity_db
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your-secret-key-change-in-production
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
cd identity-service
go test ./internal/service/...
```

### Building Services

```bash
# Build identity service
cd identity-service
go build -o bin/identity-service cmd/main.go

# Build API Gateway
cd api-gateway
go build -o bin/api-gateway cmd/main.go
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check for vulnerabilities
go list -json -m all | nancy sleuth
```

## 📚 API Documentation

### Swagger Documentation

Each service exposes Swagger documentation:

- API Gateway: http://localhost:8000/swagger/index.html
- Identity Service: http://localhost:8081/swagger/index.html

### Authentication Flow

```bash
# 1. Register
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "full_name": "Test User"
  }'

# 2. Login (receives session_id cookie)
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }' \
  --cookie-jar cookies.txt

# 3. Access protected route
curl -X GET http://localhost:8000/api/v1/users/profile \
  -H "Authorization: Bearer <access_token>"

# 4. Refresh token
curl -X POST http://localhost:8000/api/v1/auth/refresh \
  --cookie cookies.txt

# 5. Logout
curl -X POST http://localhost:8000/api/v1/auth/logout \
  --cookie cookies.txt
```

## 🧪 Testing

### Manual Testing

See [SESSION-AUTH-IMPLEMENTATION.md](./SESSION-AUTH-IMPLEMENTATION.md) for detailed testing guide.

### Load Testing (Planned)

- Target: 10,000 concurrent users
- Tools: k6, Artillery
- Scenarios: Login flow, product browsing, checkout

## 🚢 Deployment

### Docker Deployment

```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f identity-service
```

### Kubernetes (Planned)

Kubernetes manifests will be added for production deployment.

## 🗺️ Roadmap

### ✅ Gate 1: Identity Service (Completed)

- User registration and authentication
- Session-based auth with Redis
- JWT token management
- User profile management

### 🔄 Gate 2: Product Service (In Progress)

- Product CRUD operations
- Category management
- Variant and SKU system
- Image upload

### 📋 Gate 3-12 (Planned)

See [GATES.md](./docs/GATES.md) for detailed roadmap.

## 📖 Documentation

- [System Architecture](./docs/BLUEPRINT.md) - Overall system design
- [API Contracts](./docs/CONTRACT.md) - Service contracts and interfaces
- [Development Gates](./docs/GATES.md) - Phased development plan
- [Gate 1 Report](./GATE1-REPORT.md) - Identity Service completion report
- [Session Auth Implementation](./SESSION-AUTH-IMPLEMENTATION.md) - Session-based authentication details

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Huy Bang** - _Initial work_ - [huybang2017](https://github.com/huybang2017)

## 🙏 Acknowledgments

- Inspired by Shopee's system architecture
- Built for educational purposes and portfolio demonstration
- Thanks to the open-source community for amazing tools and libraries

---

**⭐ If you find this project useful, please consider giving it a star!**
