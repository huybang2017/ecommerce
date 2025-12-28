package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"search-service/config"
	"search-service/internal/handler"
	"search-service/internal/repository/elasticsearch"
	"search-service/internal/repository/kafka"
	"search-service/internal/router"
	"search-service/internal/service"
	esClient "search-service/pkg/elasticsearch"
	"search-service/pkg/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Wrap main in recover to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in main: %v\n", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	log.Println("=== Search Service Starting ===")

	// Load configuration
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("✅ Config loaded")

	// Debug: Print config values
	log.Printf("Config loaded - ES Index: %s, Kafka Topic: %s, Brokers: %v",
		cfg.Elasticsearch.IndexName,
		cfg.Kafka.TopicProductUpdated,
		cfg.Kafka.Brokers,
	)

	// Initialize logger
	log.Println("Initializing logger...")
	appLogger, err := logger.NewLogger(&cfg.Logging)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()
	log.Println("✅ Logger initialized")

	appLogger.Info("Starting Search Service...",
		zap.String("elasticsearch_index", cfg.Elasticsearch.IndexName),
		zap.String("kafka_topic", cfg.Kafka.TopicProductUpdated),
		zap.Strings("kafka_brokers", cfg.Kafka.Brokers),
	)

	// Set Gin mode based on config
	gin.SetMode(cfg.Server.Mode)

	// Validate config
	if cfg.Elasticsearch.IndexName == "" {
		appLogger.Fatal("Elasticsearch index name is empty")
	}
	if cfg.Kafka.TopicProductUpdated == "" {
		appLogger.Fatal("Kafka topic is empty")
	}

	// Initialize Elasticsearch client
	appLogger.Info("Initializing Elasticsearch client...")
	esClientInstance, err := esClient.GetClient(&cfg.Elasticsearch)
	if err != nil {
		appLogger.Fatal("Failed to initialize Elasticsearch client", zap.Error(err))
	}
	appLogger.Info("✅ Elasticsearch connection established")

	// Ensure Elasticsearch index exists
	appLogger.Info("Ensuring Elasticsearch index exists...")
	if err := esClient.EnsureIndex(esClientInstance, cfg.Elasticsearch.IndexName); err != nil {
		appLogger.Warn("Failed to ensure Elasticsearch index", zap.Error(err))
	} else {
		appLogger.Info("✅ Elasticsearch index ready", zap.String("index", cfg.Elasticsearch.IndexName))
	}

	// Initialize repositories (Infrastructure Layer)
	log.Println("Initializing repositories...")
	appLogger.Info("Initializing repositories...")
	searchRepo := elasticsearch.NewSearchRepository(esClientInstance, cfg.Elasticsearch.IndexName)
	log.Println("✅ Search repository initialized")
	appLogger.Info("✅ Search repository initialized")

	// Initialize service (Business Logic Layer)
	log.Println("Initializing services...")
	appLogger.Info("Initializing services...")
	searchService := service.NewSearchService(
		searchRepo,
		appLogger,
	)
	log.Println("✅ Search service initialized")
	appLogger.Info("✅ Search service initialized")

	// Initialize handlers (Transport Layer)
	log.Println("Initializing handlers...")
	appLogger.Info("Initializing handlers...")
	searchHandler := handler.NewSearchHandler(searchService, appLogger)
	log.Println("✅ Search handler initialized")
	appLogger.Info("✅ Search handler initialized")

	// Setup router
	log.Println("Setting up router...")
	appLogger.Info("Setting up router...")
	router := router.SetupRouter(searchHandler)
	log.Println("✅ Router setup complete")
	appLogger.Info("✅ Router setup complete")

	// Initialize Kafka consumer
	log.Println("Initializing Kafka consumer...")
	appLogger.Info("Initializing Kafka consumer...",
		zap.String("topic", cfg.Kafka.TopicProductUpdated),
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.String("consumer_group", cfg.Kafka.ConsumerGroup),
	)

	var eventConsumer *kafka.EventConsumer
	var ctx context.Context
	var cancel context.CancelFunc

	func() {
		defer func() {
			if r := recover(); r != nil {
				appLogger.Error("Panic during Kafka consumer initialization", zap.Any("panic", r))
				debug.PrintStack()
			}
		}()

		log.Println("Creating Kafka event consumer...")
		appLogger.Info("Creating Kafka event consumer...")
		eventConsumer = kafka.NewEventConsumer(
			cfg.Kafka.Brokers,
			cfg.Kafka.TopicProductUpdated,
			cfg.Kafka.ConsumerGroup,
			cfg.Kafka.ReadTimeout,
			cfg.Kafka.MinBytes,
			cfg.Kafka.MaxBytes,
			searchRepo,
			appLogger,
		)
		log.Println("✅ Kafka event consumer created")
		appLogger.Info("✅ Kafka event consumer created")

		// Start Kafka consumer in background
		log.Println("Setting up Kafka consumer context...")
		appLogger.Info("Setting up Kafka consumer context...")
		ctx, cancel = context.WithCancel(context.Background())
		log.Println("✅ Context created")
		appLogger.Info("✅ Context created")

		log.Println("Starting Kafka consumer goroutine...")
		appLogger.Info("Starting Kafka consumer goroutine...")
		log.Println("🚀 About to start Kafka consumer goroutine...")
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ Kafka consumer goroutine panicked: %v\n", r)
					appLogger.Error("Kafka consumer goroutine panicked", zap.Any("panic", r))
					debug.PrintStack()
				}
			}()
			log.Println("✅ Kafka consumer goroutine started, calling Start()...")
			appLogger.Info("Kafka consumer goroutine started, calling Start()...")
			if err := eventConsumer.Start(ctx); err != nil {
				log.Printf("❌ Kafka consumer stopped with error: %v\n", err)
				appLogger.Error("Kafka consumer stopped", zap.Error(err))
			} else {
				log.Println("ℹ️ Kafka consumer Start() returned without error")
				appLogger.Info("Kafka consumer Start() returned")
			}
		}()

		// Give Kafka consumer a moment to start
		log.Println("Waiting for Kafka consumer to initialize...")
		appLogger.Info("Waiting for Kafka consumer to initialize...")
		time.Sleep(2 * time.Second)
		log.Println("✅ Kafka consumer started in background")
		appLogger.Info("✅ Kafka consumer started in background")
	}()

	// Setup cleanup
	defer func() {
		appLogger.Info("Cleaning up Kafka consumer...")
		if cancel != nil {
			cancel()
		}
		if eventConsumer != nil {
			eventConsumer.Close()
		}
		appLogger.Info("✅ Kafka consumer cleaned up")
	}()

	// Create HTTP server with timeouts
	log.Println("Creating HTTP server...")
	appLogger.Info("Creating HTTP server...")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	log.Println("✅ HTTP server created")
	appLogger.Info("✅ HTTP server created", zap.Int("port", cfg.Server.Port))

	// Start server in a goroutine
	log.Println("Starting HTTP server goroutine...")
	serverErr := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("HTTP server goroutine panicked: %v\n", r)
				appLogger.Error("HTTP server goroutine panicked", zap.Any("panic", r))
				debug.PrintStack()
				serverErr <- fmt.Errorf("panic: %v", r)
			}
		}()
		log.Println("HTTP server goroutine started, calling ListenAndServe...")
		appLogger.Info("Starting HTTP server...", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v\n", err)
			appLogger.Error("HTTP server error", zap.Error(err))
			serverErr <- err
		}
	}()

	// Give server a moment to start
	log.Println("Waiting for HTTP server to start...")
	appLogger.Info("Waiting for HTTP server to start...")
	time.Sleep(1 * time.Second)
	
	// Test if server is actually listening
	log.Println("Testing HTTP server health endpoint...")
	testCtx, testCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer testCancel()
	
	testReq, _ := http.NewRequestWithContext(testCtx, "GET", fmt.Sprintf("http://localhost:%d/health", cfg.Server.Port), nil)
	resp, err := http.DefaultClient.Do(testReq)
	if err != nil {
		log.Printf("Server health check failed: %v\n", err)
		appLogger.Warn("Server health check failed (may be starting)", zap.Error(err))
	} else {
		resp.Body.Close()
		log.Println("✅ HTTP server is responding")
		appLogger.Info("✅ HTTP server is responding")
	}

	log.Println("✅✅✅ Search Service is ready ✅✅✅")
	appLogger.Info("✅✅✅ Search Service is ready ✅✅✅",
		zap.Int("port", cfg.Server.Port),
		zap.String("kafka_topic", cfg.Kafka.TopicProductUpdated),
		zap.Strings("elasticsearch_addresses", cfg.Elasticsearch.Addresses),
	)

	// Wait for interrupt signal or server error
	log.Println("Waiting for interrupt signal or server error...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	appLogger.Info("Waiting for interrupt signal or server error...")
	
	log.Println("Entering select statement...")
	select {
	case sig := <-quit:
		log.Printf("Received interrupt signal: %v\n", sig)
		appLogger.Info("Received interrupt signal", zap.String("signal", sig.String()))
	case err := <-serverErr:
		log.Printf("HTTP server error received: %v\n", err)
		appLogger.Error("HTTP server error received", zap.Error(err))
		appLogger.Fatal("Server failed", zap.Error(err))
	}

	appLogger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Cancel Kafka consumer context
	cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	appLogger.Info("Server exited")
}

