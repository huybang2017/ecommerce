# Multi-stage Dockerfile for Go application
# This reduces the final image size by excluding build tools

# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies (cached layer if go.mod/go.sum unchanged)
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary
# -ldflags="-w -s" reduces binary size by stripping debug info
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o order-service ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install CA certificates for HTTPS requests and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/order-service .

# Copy config file
COPY --from=builder /build/config/config.yaml ./config/

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8083

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
  CMD wget -O - --no-verbose --tries=1 http://localhost:8083/health || exit 1

# Run the application
CMD ["./order-service"]


