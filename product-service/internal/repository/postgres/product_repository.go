package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// productRepository implements the ProductRepository interface
// This is the infrastructure layer - it knows HOW to interact with PostgreSQL
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new PostgreSQL product repository
// Dependency injection: we inject the database connection
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database
func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

// Update updates an existing product
func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

// GetByID retrieves a product by its ID
func (r *productRepository) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by its SKU
func (r *productRepository) GetBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products
func (r *productRepository) GetAll() ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// ListProducts retrieves products with pagination and filters
func (r *productRepository) ListProducts(filters map[string]interface{}, page, limit int) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	// Build query with filters
	query := r.db.Model(&domain.Product{})

	// Apply filters
	if categoryID, ok := filters["category_id"]; ok {
		query = query.Where("category_id = ?", categoryID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if minPrice, ok := filters["min_price"]; ok {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice, ok := filters["max_price"]; ok {
		query = query.Where("price <= ?", maxPrice)
	}
	if search, ok := filters["search"]; ok {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search.(string)+"%", "%"+search.(string)+"%")
	}

	// Count total (before pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetProductsByCategory retrieves products by category ID with pagination
func (r *productRepository) GetProductsByCategory(categoryID uint, page, limit int) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	// Count total
	if err := r.db.Model(&domain.Product{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get products with pagination
	offset := (page - 1) * limit
	if err := r.db.Where("category_id = ?", categoryID).Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Delete soft deletes a product (or hard delete based on your business logic)
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

// GetProductsByShopID retrieves products by shop ID with pagination
func (r *productRepository) GetProductsByShopID(shopID uint, page, limit int) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&domain.Product{}).Where("shop_id = ?", shopID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloaded Category
	if err := r.db.Preload("Category").Where("shop_id = ?", shopID).
		Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

