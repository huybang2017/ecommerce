package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"api-gateway/config"
	"api-gateway/internal/domain"
	"api-gateway/internal/handler"
	"api-gateway/internal/repository"
	"api-gateway/internal/router"
	"api-gateway/internal/service"
	"api-gateway/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting API Gateway...")

	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)

	// Initialize service registry
	serviceRegistry := repository.NewServiceRegistry()

	// Register microservices from configuration
	// Product Service
	productServiceConfig, exists := cfg.Services["product_service"]
	if !exists {
		appLogger.Fatal("Product service configuration not found")
	}

	// Get base URL from config or environment variable
	baseURL := productServiceConfig.BaseURL
	if baseURL == "" {
		// Fallback to environment variable or default
		baseURL = "http://product-service:8080"
		appLogger.Warn("Using default base URL for product service", zap.String("url", baseURL))
	}

	productService := &domain.Service{
		Name:            "product_service",
		BaseURL:         baseURL,
		HealthCheckPath: productServiceConfig.HealthCheckPath,
		Routes: []domain.Route{
			{Path: "/api/v1/products", Methods: []string{"GET", "POST"}, RequireAuth: false},
			{Path: "/api/v1/products/:id", Methods: []string{"GET"}, RequireAuth: false},
			{Path: "/api/v1/products/search", Methods: []string{"GET"}, RequireAuth: false},
			{Path: "/api/v1/products/:id", Methods: []string{"PUT", "PATCH", "DELETE"}, RequireAuth: true},
			{Path: "/api/v1/products/:id/inventory", Methods: []string{"PATCH"}, RequireAuth: true},
		},
	}

	if err := serviceRegistry.RegisterService(productService); err != nil {
		appLogger.Fatal("Failed to register product service", zap.Error(err))
	}

	// Register Identity Service
	identityServiceConfig, exists := cfg.Services["identity_service"]
	if exists {
		identityBaseURL := identityServiceConfig.BaseURL
		if identityBaseURL == "" {
			identityBaseURL = "http://localhost:8081"
			appLogger.Warn("Using default base URL for identity service", zap.String("url", identityBaseURL))
		}

		identityService := &domain.Service{
			Name:            "identity_service",
			BaseURL:         identityBaseURL,
			HealthCheckPath: identityServiceConfig.HealthCheckPath,
			Routes: []domain.Route{
				{Path: "/api/v1/auth/register", Methods: []string{"POST"}, RequireAuth: false},
				{Path: "/api/v1/auth/login", Methods: []string{"POST"}, RequireAuth: false},
				{Path: "/api/v1/users/profile", Methods: []string{"GET", "PUT"}, RequireAuth: true},
				{Path: "/api/v1/users/password", Methods: []string{"PUT"}, RequireAuth: true},
				{Path: "/api/v1/addresses", Methods: []string{"GET", "POST"}, RequireAuth: true},
				{Path: "/api/v1/addresses/:id", Methods: []string{"GET", "PUT", "DELETE"}, RequireAuth: true},
				{Path: "/api/v1/addresses/:id/default", Methods: []string{"PUT"}, RequireAuth: true},
			},
		}

		if err := serviceRegistry.RegisterService(identityService); err != nil {
			appLogger.Fatal("Failed to register identity service", zap.Error(err))
		}
		appLogger.Info("Identity service registered", zap.String("base_url", identityBaseURL))
	}

	// Initialize proxy client (use max timeout from all services)
	maxTimeout := productServiceConfig.Timeout
	if exists && identityServiceConfig.Timeout > maxTimeout {
		maxTimeout = identityServiceConfig.Timeout
	}
	proxyClient := repository.NewProxyClient(maxTimeout)

	// Initialize gateway service
	gatewayService := service.NewGatewayService(serviceRegistry, proxyClient, appLogger)

	// Initialize handlers
	gatewayHandler := handler.NewGatewayHandler(gatewayService, appLogger)

	// Setup router
	router := router.SetupRouter(gatewayHandler, cfg, appLogger)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("API Gateway starting", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down API Gateway...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("API Gateway exited gracefully")
}

