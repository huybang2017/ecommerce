package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-service/config"
	"product-service/internal/domain"
	"product-service/internal/handler"
	"product-service/internal/repository/elasticsearch"
	"product-service/internal/repository/kafka"
	"product-service/internal/repository/postgres"
	"product-service/internal/repository/redis"
	"product-service/internal/router"
	"product-service/internal/service"
	"product-service/pkg/database"
	esClient "product-service/pkg/elasticsearch"
	"product-service/pkg/logger"
	redisClient "product-service/pkg/redis"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	fmt.Fprintf(os.Stderr, "🚀🚀🚀 PRODUCT SERVICE MAIN() STARTED! 🚀🚀🚀\n")
	log.Printf("🚀 PRODUCT SERVICE MAIN() STARTED!")
	
	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Fprintf(os.Stderr, "✅ Config loaded - Topic: %s, Brokers: %v\n", cfg.Kafka.TopicProductUpdated, cfg.Kafka.Brokers)

	// Initialize logger
	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	appLogger.Info("Starting Product Service...")

	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection (Singleton)
	db, err := database.GetDB(&cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.CloseDB()

	// Run database migrations
	if err := db.AutoMigrate(&domain.Product{}, &domain.Category{}); err != nil {
		appLogger.Fatal("Failed to run migrations", zap.Error(err))
	}
	appLogger.Info("Database migrations completed")

	// Initialize Redis client (Singleton)
	redisClientInstance, err := redisClient.GetClient(&cfg.Redis)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.CloseClient()

	// Initialize Elasticsearch client (Singleton)
	esClientInstance, err := esClient.GetClient(&cfg.Elasticsearch)
	if err != nil {
		appLogger.Fatal("Failed to connect to Elasticsearch", zap.Error(err))
	}

	// Ensure Elasticsearch index exists
	if err := esClient.EnsureIndex(esClientInstance, cfg.Elasticsearch.IndexName); err != nil {
		appLogger.Warn("Failed to ensure Elasticsearch index", zap.Error(err))
	}

	// Initialize Kafka event publisher
	log.Printf("🔧 Initializing Kafka event publisher - brokers: %v, topic: %s", cfg.Kafka.Brokers, cfg.Kafka.TopicProductUpdated)
	appLogger.Info("Initializing Kafka event publisher",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("topic", cfg.Kafka.TopicProductUpdated),
	)
	eventPublisher := kafka.NewEventPublisher(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicProductUpdated,
		cfg.Kafka.WriteTimeout,
		cfg.Kafka.RequiredAcks,
	)
	if eventPublisher == nil {
		log.Printf("❌❌❌ Failed to create Kafka event publisher - eventPublisher is nil")
		appLogger.Fatal("Failed to create Kafka event publisher")
	}
	log.Printf("✅✅✅ Kafka event publisher initialized successfully")
	appLogger.Info("✅ Kafka event publisher initialized successfully")
	defer eventPublisher.Close()

	// Initialize repositories (Infrastructure Layer)
	productRepo := postgres.NewProductRepository(db)
	categoryRepo := postgres.NewCategoryRepository(db)
	searchRepo := elasticsearch.NewProductSearchRepository(esClientInstance, cfg.Elasticsearch.IndexName)
	cacheRepo := redis.NewCacheRepository(redisClientInstance)

	// Initialize service (Business Logic Layer)
	fmt.Fprintf(os.Stderr, "🔧 Creating ProductService with eventPublisher: %p\n", eventPublisher)
	productService := service.NewProductService(
		productRepo,
		searchRepo,
		cacheRepo,
		eventPublisher,
		appLogger,
	)
	fmt.Fprintf(os.Stderr, "✅ ProductService created - eventPublisher injected: %p\n", eventPublisher)
	categoryService := service.NewCategoryService(
		categoryRepo,
		appLogger,
	)

	// Initialize handlers (Transport Layer)
	fmt.Fprintf(os.Stderr, "🔧 Creating handlers...\n")
	productHandler := handler.NewProductHandler(productService, appLogger)
	categoryHandler := handler.NewCategoryHandler(categoryService, appLogger)
	fmt.Fprintf(os.Stderr, "✅ Handlers created - ProductHandler: %p, eventPublisher in service: %p\n", productHandler, productService)

	// Setup router
	router := router.SetupRouter(productHandler, categoryHandler)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				appLogger.Error("Server goroutine panicked", zap.Any("panic", r))
				log.Printf("Server goroutine panicked: %v", r)
			}
		}()
		appLogger.Info("Server starting", zap.Int("port", cfg.Server.Port))
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server error", zap.Error(err))
			log.Printf("Server error: %v", err)
			// Don't use Fatal here - it will exit the entire program
			// Instead, log the error and let the main goroutine handle shutdown
		}
	}()

	// Give server a moment to start
	time.Sleep(2 * time.Second)
	
	// Verify server is running
	testCtx, testCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer testCancel()
	testReq, _ := http.NewRequestWithContext(testCtx, "GET", fmt.Sprintf("http://localhost:%d/health", cfg.Server.Port), nil)
	resp, err := http.DefaultClient.Do(testReq)
	if err != nil {
		appLogger.Warn("Server health check failed (may be starting)", zap.Error(err))
		log.Printf("Server health check failed: %v", err)
	} else {
		resp.Body.Close()
		appLogger.Info("Server is responding", zap.Int("port", cfg.Server.Port))
		log.Printf("Server is responding on port %d", cfg.Server.Port)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	appLogger.Info("Product Service is ready and waiting for requests", zap.Int("port", cfg.Server.Port))
	log.Printf("Product Service is ready and waiting for requests on port %d", cfg.Server.Port)
	<-quit

	appLogger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close all connections
	// Note: Kafka publisher and Redis/ES clients are closed via defer
	appLogger.Info("Server exited gracefully")
}

