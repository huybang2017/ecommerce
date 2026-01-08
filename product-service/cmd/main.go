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
	"gorm.io/gorm"
)

// migrateProductsTable handles the migration of products table with special handling for shop_id
// This is needed because we cannot add a NOT NULL column to a table with existing data
func migrateProductsTable(db *gorm.DB, logger *zap.Logger) error {
	// Check if shop_id column already exists
	var count int64
	err := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = CURRENT_SCHEMA() 
		AND table_name = 'products' 
		AND column_name = 'shop_id'
	`).Scan(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check shop_id column: %w", err)
	}

	if count == 0 {
		// Step 1: Add shop_id column as nullable first
		logger.Info("Adding shop_id column as nullable...")
		if err := db.Exec(`ALTER TABLE products ADD COLUMN shop_id bigint`).Error; err != nil {
			return fmt.Errorf("failed to add shop_id column: %w", err)
		}

		// Step 2: Update all existing products with default shop_id = 1
		// NOTE: This assumes shop with id=1 exists (should be created by Identity Service)
		logger.Info("Updating existing products with default shop_id = 1...")
		if err := db.Exec(`UPDATE products SET shop_id = 1 WHERE shop_id IS NULL`).Error; err != nil {
			logger.Warn("Failed to update products with shop_id, will set to 1 anyway", zap.Error(err))
			// Continue even if update fails (might be no products yet)
		}

		// Step 3: Set NOT NULL constraint
		logger.Info("Setting shop_id NOT NULL constraint...")
		if err := db.Exec(`ALTER TABLE products ALTER COLUMN shop_id SET NOT NULL`).Error; err != nil {
			return fmt.Errorf("failed to set shop_id NOT NULL: %w", err)
		}

		// Step 4: Add index
		logger.Info("Adding index on shop_id...")
		if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_products_shop_id ON products(shop_id)`).Error; err != nil {
			logger.Warn("Failed to create index on shop_id", zap.Error(err))
			// Continue even if index creation fails (might already exist)
		}
	}

	// Step 5: Handle base_price column (similar to shop_id)
	var basePriceCount int64
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = CURRENT_SCHEMA() 
		AND table_name = 'products' 
		AND column_name = 'base_price'
	`).Scan(&basePriceCount).Error
	if err != nil {
		return fmt.Errorf("failed to check base_price column: %w", err)
	}

	if basePriceCount == 0 {
		// Step 5.1: Add base_price column as nullable first
		logger.Info("Adding base_price column as nullable...")
		if err := db.Exec(`ALTER TABLE products ADD COLUMN base_price decimal(15,2)`).Error; err != nil {
			return fmt.Errorf("failed to add base_price column: %w", err)
		}

		// Step 5.2: Update all existing products with base_price = price (copy from existing price)
		logger.Info("Updating existing products with base_price = price...")
		if err := db.Exec(`UPDATE products SET base_price = price WHERE base_price IS NULL`).Error; err != nil {
			logger.Warn("Failed to update products with base_price", zap.Error(err))
			// Continue even if update fails (might be no products yet)
		}

		// Step 5.3: Set NOT NULL constraint
		logger.Info("Setting base_price NOT NULL constraint...")
		if err := db.Exec(`ALTER TABLE products ALTER COLUMN base_price SET NOT NULL`).Error; err != nil {
			return fmt.Errorf("failed to set base_price NOT NULL: %w", err)
		}
	}

	// Step 6: Handle sold_count column (default to 0)
	var soldCountCount int64
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.columns 
		WHERE table_schema = CURRENT_SCHEMA() 
		AND table_name = 'products' 
		AND column_name = 'sold_count'
	`).Scan(&soldCountCount).Error
	if err != nil {
		return fmt.Errorf("failed to check sold_count column: %w", err)
	}

	if soldCountCount == 0 {
		// Step 6.1: Add sold_count column with default 0
		logger.Info("Adding sold_count column with default 0...")
		if err := db.Exec(`ALTER TABLE products ADD COLUMN sold_count integer DEFAULT 0 NOT NULL`).Error; err != nil {
			return fmt.Errorf("failed to add sold_count column: %w", err)
		}
	}

	// Now run AutoMigrate for Product (will handle other fields)
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		return fmt.Errorf("failed to auto-migrate Product: %w", err)
	}

	logger.Info("Products table migration completed")
	return nil
}

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
	// NOTE: Special handling for shop_id column - must add nullable first, then update data, then set NOT NULL
	if err := migrateProductsTable(db, appLogger); err != nil {
		appLogger.Fatal("Failed to migrate products table", zap.Error(err))
	}

	// AutoMigrate other tables
	if err := db.AutoMigrate(
		&domain.Category{},
		&domain.Variation{},
		&domain.VariationOption{},
		&domain.ProductItem{},
		&domain.SKUConfiguration{},
		&domain.CategoryAttribute{},
		&domain.ProductAttributeValue{},
	); err != nil {
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
	variationRepo := postgres.NewVariationRepository(db)
	variationOptRepo := postgres.NewVariationOptionRepository(db)
	productItemRepo := postgres.NewProductItemRepository(db)
	skuConfigRepo := postgres.NewSKUConfigurationRepository(db)
	categoryAttrRepo := postgres.NewCategoryAttributeRepository(db)
	productAttrRepo := postgres.NewProductAttributeValueRepository(db)
	searchRepo := elasticsearch.NewProductSearchRepository(esClientInstance, cfg.Elasticsearch.IndexName)
	cacheRepo := redis.NewCacheRepository(redisClientInstance)

	// Initialize services (Business Logic Layer)
	fmt.Fprintf(os.Stderr, "🔧 Creating ProductService with eventPublisher: %p\n", eventPublisher)
	productService := service.NewProductService(
		productRepo,
		searchRepo,
		cacheRepo,
		categoryRepo,
		eventPublisher,
		appLogger,
	)
	fmt.Fprintf(os.Stderr, "✅ ProductService created - eventPublisher injected: %p\n", eventPublisher)
	categoryService := service.NewCategoryService(
		categoryRepo,
		appLogger,
	)
	productItemService := service.NewProductItemService(
		productItemRepo,
		variationRepo,
		variationOptRepo,
		skuConfigRepo,
		productRepo,
		appLogger,
	)
	attributeService := service.NewAttributeService(
		categoryAttrRepo,
		productAttrRepo,
		categoryRepo,
		productRepo,
		appLogger,
	)
	stockService := service.NewStockService(
		productItemRepo,
		redisClientInstance,
		appLogger,
	)

	// Initialize handlers (Transport Layer)
	fmt.Fprintf(os.Stderr, "🔧 Creating handlers...\n")
	productHandler := handler.NewProductHandler(productService, appLogger)
	categoryHandler := handler.NewCategoryHandler(categoryService, appLogger)
	skuHandler := handler.NewSKUHandler(productItemService, appLogger)
	attrHandler := handler.NewAttributeHandler(attributeService, appLogger)
	stockHandler := handler.NewStockHandler(stockService, appLogger)
	fmt.Fprintf(os.Stderr, "✅ Handlers created - ProductHandler: %p, eventPublisher in service: %p\n", productHandler, productService)

	// Setup router
	router := router.SetupRouter(productHandler, categoryHandler, skuHandler, attrHandler, stockHandler)

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
