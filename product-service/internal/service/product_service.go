package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"product-service/internal/domain"
	"time"

	"go.uber.org/zap"
)

// ProductService contains the business logic for product operations
// This is the service layer - it orchestrates between repositories
// Following Clean Architecture: business logic is independent of infrastructure
type ProductService struct {
	productRepo    domain.ProductRepository
	searchRepo     domain.ProductSearchRepository
	cacheRepo      CacheRepository
	categoryRepo   domain.CategoryRepository
	eventPublisher domain.EventPublisher
	logger         *zap.Logger
}

// CacheRepository defines cache operations (abstraction for Redis)
// This interface allows us to swap Redis for other caching solutions if needed
type CacheRepository interface {
	SetProduct(ctx context.Context, product *domain.Product, ttl time.Duration) error
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	AcquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, lockKey string) error
}

// NewProductService creates a new product service with all dependencies
// Dependency injection: we inject all repositories and external services
func NewProductService(
	productRepo domain.ProductRepository,
	searchRepo domain.ProductSearchRepository,
	cacheRepo CacheRepository,
	categoryRepo domain.CategoryRepository,
	eventPublisher domain.EventPublisher,
	logger *zap.Logger,
) *ProductService {
	return &ProductService{
		productRepo:    productRepo,
		searchRepo:     searchRepo,
		cacheRepo:      cacheRepo,
		categoryRepo:   categoryRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateProduct creates a new product with full integration
// This demonstrates the orchestration pattern:
// 1. Save to PostgreSQL (source of truth)
// 2. Update Redis cache (fast reads)
// 3. Index to Elasticsearch (search capability)
// 4. Publish event to Kafka (event-driven architecture)
func (s *ProductService) CreateProduct(ctx context.Context, product *domain.Product) error {
	// Business logic validation
	if product.Name == "" {
		return errors.New("name is required")
	}
	if product.BasePrice < 0 {
		return errors.New("base price cannot be negative")
	}

	// 1. Save to PostgreSQL (source of truth)
	fmt.Fprintf(os.Stderr, "🟢🟢🟢 Service: About to create product in DB - Name: %s\n", product.Name)
	log.Printf("🟢 Service: About to create product in DB - Name: %s", product.Name)
	if err := s.productRepo.Create(product); err != nil {
		fmt.Fprintf(os.Stderr, "❌❌❌ Service: Failed to create product in DB: %v\n", err)
		log.Printf("❌ Service: Failed to create product in DB: %v", err)
		s.logger.Error("failed to create product in database", zap.Error(err))
		return fmt.Errorf("failed to create product: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✅✅✅ Service: Product created in DB - ID: %d, Name: %s\n", product.ID, product.Name)
	log.Printf("✅ Service: Product created in DB - ID: %d, Name: %s", product.ID, product.Name)
	s.logger.Info("product created in database", zap.Uint("product_id", product.ID))
	_ = s.logger.Sync()

	// 2. Update Redis cache (async - don't block on cache)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cacheRepo.SetProduct(cacheCtx, product, 1*time.Hour); err != nil {
			s.logger.Warn("failed to cache product", zap.Error(err))
		}
	}()

	// 3. Index to Elasticsearch (async - search is eventually consistent)
	go func() {
		if err := s.searchRepo.IndexProduct(product); err != nil {
			s.logger.Warn("failed to index product in elasticsearch", zap.Error(err))
		} else {
			s.logger.Info("product indexed in elasticsearch", zap.Uint("product_id", product.ID))
		}
	}()

	// 4. Publish event to Kafka (async - event-driven communication)
	// CRITICAL: Log BEFORE starting goroutine to confirm we reach this point
	s.logger.Info("🔵🔵🔵 ABOUT TO START EVENT PUBLISHING GOROUTINE",
		zap.Uint("product_id", product.ID),
		zap.String("product_name", product.Name),
		zap.Bool("eventPublisher_nil", s.eventPublisher == nil),
	)
	_ = s.logger.Sync()

	go func() {
		// CRITICAL: Use Zap logger with Sync to ensure logs are flushed immediately
		s.logger.Info("🚀🚀🚀 EVENT PUBLISHING GOROUTINE CALLED!",
			zap.Uint("product_id", product.ID),
			zap.String("product_name", product.Name),
		)
		_ = s.logger.Sync() // Force flush logs immediately

		// Check if eventPublisher is nil
		if s.eventPublisher == nil {
			s.logger.Error("❌❌❌ Event publisher is nil - cannot publish event",
				zap.Uint("product_id", product.ID),
				zap.String("product_name", product.Name),
			)
			_ = s.logger.Sync()
			return
		}

		event := &domain.ProductEvent{
			EventType:   "product_created",
			ProductID:   product.ID,
			ProductData: product,
			Timestamp:   time.Now(),
		}

		s.logger.Info("📤 Publishing product event to Kafka",
			zap.Uint("product_id", product.ID),
			zap.String("event_type", event.EventType),
			zap.String("product_name", product.Name),
		)
		_ = s.logger.Sync()

		if err := s.eventPublisher.PublishProductEvent(event); err != nil {
			s.logger.Error("❌❌❌ Failed to publish product event to Kafka",
				zap.Uint("product_id", event.ProductID),
				zap.String("event_type", event.EventType),
				zap.Error(err),
			)
			_ = s.logger.Sync()
		} else {
			s.logger.Info("✅✅✅ Product event published to Kafka successfully",
				zap.Uint("product_id", event.ProductID),
				zap.String("event_type", event.EventType),
				zap.String("product_name", product.Name),
			)
			_ = s.logger.Sync()
		}
	}()

	return nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	// Validate product exists
	existing, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return errors.New("product not found")
	}

	// Business logic: preserve created_at
	product.CreatedAt = existing.CreatedAt

	// 1. Update in PostgreSQL
	if err := s.productRepo.Update(product); err != nil {
		s.logger.Error("failed to update product in database", zap.Error(err))
		return fmt.Errorf("failed to update product: %w", err)
	}

	s.logger.Info("product updated in database", zap.Uint("product_id", product.ID))

	// 2. Update cache
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cacheRepo.SetProduct(cacheCtx, product, 1*time.Hour); err != nil {
			s.logger.Warn("failed to update product cache", zap.Error(err))
		}
	}()

	// 3. Update Elasticsearch index
	go func() {
		if err := s.searchRepo.IndexProduct(product); err != nil {
			s.logger.Warn("failed to update product in elasticsearch", zap.Error(err))
		}
	}()

	// 4. Publish update event
	go func() {
		event := &domain.ProductEvent{
			EventType:   "product_updated",
			ProductID:   product.ID,
			ProductData: product,
			Timestamp:   time.Now(),
		}

		if err := s.eventPublisher.PublishProductEvent(event); err != nil {
			s.logger.Warn("failed to publish product update event", zap.Error(err))
		}
	}()

	return nil
}

// GetProduct retrieves a product by ID with cache-first strategy
// This demonstrates the cache-aside pattern
func (s *ProductService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	// 1. Try cache first (fast path)
	product, err := s.cacheRepo.GetProduct(ctx, id)
	if err == nil && product != nil {
		s.logger.Debug("product retrieved from cache", zap.Uint("product_id", id))
		return product, nil
	}

	// 2. Cache miss - get from database (slow path)
	product, err = s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// 3. Populate cache for next time (async)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.cacheRepo.SetProduct(cacheCtx, product, 1*time.Hour); err != nil {
			s.logger.Warn("failed to cache product", zap.Error(err))
		}
	}()

	return product, nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.productRepo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all products", zap.Error(err))
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return products, nil
}

// ListProducts retrieves products with pagination and filters
func (s *ProductService) ListProducts(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*domain.Product, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	products, total, err := s.productRepo.ListProducts(filters, page, limit)
	if err != nil {
		s.logger.Error("failed to list products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, total, nil
}

// GetProductsByCategory retrieves products by category ID with pagination
// If category is a parent (has children), it will fetch products from all child categories too
func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID uint, page, limit int) ([]*domain.Product, int64, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	// Build category IDs array (include category and its children recursively)
	categoryIDs := []uint{categoryID}

	// Recursive helper to get all descendants
	var getAllDescendants func(parentID uint)
	getAllDescendants = func(parentID uint) {
		children, err := s.categoryRepo.GetChildren(parentID)
		if err == nil && len(children) > 0 {
			s.logger.Debug("found children for category",
				zap.Uint("parent_id", parentID),
				zap.Int("children_count", len(children)))
			for _, child := range children {
				categoryIDs = append(categoryIDs, child.ID)
				// Recursively get grandchildren
				getAllDescendants(child.ID)
			}
		}
	}

	// Get all descendants of this category
	getAllDescendants(categoryID)

	s.logger.Info("fetching products for category tree",
		zap.Uint("root_category_id", categoryID),
		zap.Int("total_categories", len(categoryIDs)),
		zap.Uints("category_ids", categoryIDs))

	products, total, err := s.productRepo.GetProductsByCategoryIDs(categoryIDs, page, limit)
	if err != nil {
		s.logger.Error("failed to get products by category", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get products by category: %w", err)
	}

	return products, total, nil
}

// SearchProducts searches products using Elasticsearch
func (s *ProductService) SearchProducts(ctx context.Context, query string, filters map[string]interface{}) ([]*domain.Product, error) {
	products, err := s.searchRepo.SearchProducts(query, filters)
	if err != nil {
		s.logger.Error("failed to search products", zap.Error(err))
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	return products, nil
}
