package domain

import (
	"time"

	"gorm.io/datatypes"
)

// Product represents the core domain entity
// This is the business object that exists independently of infrastructure
// Following Clean Architecture: domain layer has no external dependencies
// NOTE: Following db-diagram.db schema (SOURCE OF TRUTH)
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ShopID      uint           `gorm:"index;not null" json:"shop_id"` // Product thuộc shop (theo db-diagram.db)
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	BasePrice   float64        `gorm:"column:base_price;type:decimal(15,2);not null" json:"base_price"` // Giá gốc (theo db-diagram.db)
	Price       float64        `gorm:"not null" json:"price"`                                           // GIỮ LẠI để backward compatibility (sẽ sync với BasePrice)
	SKU         string         `gorm:"uniqueIndex;not null" json:"sku"`                                 // GIỮ LẠI (sẽ deprecated sau khi có product_item)
	CategoryID  *uint          `gorm:"index" json:"category_id,omitempty"`                              // Foreign key to categories
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Status      string         `gorm:"default:'ACTIVE'" json:"status"`                // ACTIVE, INACTIVE
	Images      datatypes.JSON `gorm:"type:jsonb" json:"images"`                      // JSON array of image URLs
	Stock       int            `gorm:"default:0" json:"stock"`                        // GIỮ LẠI (sẽ deprecated sau khi có product_item)
	IsActive    bool           `gorm:"default:true" json:"is_active"`                 // Boolean theo db-diagram.db
	SoldCount   int            `gorm:"column:sold_count;default:0" json:"sold_count"` // Số lượng đã bán (theo db-diagram.db)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Product) TableName() string {
	return "products"
}

// ProductRepository defines the interface for product data access
// This is part of the domain layer - it defines WHAT we need, not HOW
// The implementation will be in the repository layer (infrastructure)
type ProductRepository interface {
	Create(product *Product) error
	Update(product *Product) error
	GetByID(id uint) (*Product, error)
	GetBySKU(sku string) (*Product, error)
	GetAll() ([]*Product, error)
	ListProducts(filters map[string]interface{}, page, limit int) ([]*Product, int64, error)
	GetProductsByCategory(categoryID uint, page, limit int) ([]*Product, int64, error)
	GetProductsByCategoryIDs(categoryIDs []uint, page, limit int) ([]*Product, int64, error)
	GetProductsByShopID(shopID uint, page, limit int) ([]*Product, int64, error) // THÊM MỚI - Get products by shop
	Delete(id uint) error
}

// ProductSearchRepository defines the interface for product search operations
// Separated from ProductRepository to follow Interface Segregation Principle
type ProductSearchRepository interface {
	IndexProduct(product *Product) error
	SearchProducts(query string, filters map[string]interface{}) ([]*Product, error)
	DeleteFromIndex(id uint) error
}
