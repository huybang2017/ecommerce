package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"order-service/config"
	"order-service/internal/domain"
	"order-service/internal/handler"
	"order-service/internal/repository/kafka"
	"order-service/internal/repository/postgres"
	"order-service/internal/repository/redis"
	"order-service/internal/router"
	"order-service/internal/service"
	"order-service/pkg/database"
	"order-service/pkg/logger"
	redisClient "order-service/pkg/redis"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Order Service API
// @version 1.0
// @description Order Service API for e-commerce platform - Cart and Order management endpoints
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8083
// @BasePath /api/v1
// @schemes http https

func main() {
	log.Println("🚀 Starting Order Service...")

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

	appLogger.Info("Starting Order Service...")

	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection (Singleton)
	db, err := database.GetDB(&cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.CloseDB()

	// Run database migrations
	if err := db.AutoMigrate(&domain.Order{}, &domain.OrderItem{}); err != nil {
		appLogger.Fatal("Failed to run migrations", zap.Error(err))
	}
	appLogger.Info("Database migrations completed")

	// Initialize Redis client (Singleton)
	redisClientInstance, err := redisClient.GetClient(&cfg.Redis)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.CloseClient()

	// Initialize Kafka event publisher
	appLogger.Info("Initializing Kafka event publisher",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("topic", cfg.Kafka.TopicOrderCreated),
	)
	eventPublisher := kafka.NewEventPublisher(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicOrderCreated,
		cfg.Kafka.WriteTimeout,
		cfg.Kafka.RequiredAcks,
	)
	if eventPublisher == nil {
		appLogger.Fatal("Failed to create Kafka event publisher")
	}
	defer eventPublisher.Close()
	appLogger.Info("Kafka event publisher initialized successfully")

	// Initialize repositories
	cartRepo := redis.NewCartRepository(redisClientInstance)
	orderRepo := postgres.NewOrderRepository(db)

	// Initialize services
	cartService := service.NewCartService(cartRepo, appLogger)
	orderService := service.NewOrderService(orderRepo, cartRepo, eventPublisher, appLogger)

	// Initialize handlers
	cartHandler := handler.NewCartHandler(cartService, appLogger)
	orderHandler := handler.NewOrderHandler(orderService, appLogger)

	// Setup router
	router := router.SetupRouter(cartHandler, orderHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited gracefully")
}


