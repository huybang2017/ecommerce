package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// productItemRepository implements the ProductItemRepository interface
type productItemRepository struct {
	db *gorm.DB
}

// NewProductItemRepository creates a new PostgreSQL product item repository
func NewProductItemRepository(db *gorm.DB) domain.ProductItemRepository {
	return &productItemRepository{db: db}
}

// Create inserts a new product item (SKU) into the database
func (r *productItemRepository) Create(item *domain.ProductItem) error {
	return r.db.Create(item).Error
}

// Update updates an existing product item
func (r *productItemRepository) Update(item *domain.ProductItem) error {
	return r.db.Save(item).Error
}

// GetByID retrieves a product item by its ID
func (r *productItemRepository) GetByID(id uint) (*domain.ProductItem, error) {
	var item domain.ProductItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetBySKUCode retrieves a product item by its SKU code
func (r *productItemRepository) GetBySKUCode(skuCode string) (*domain.ProductItem, error) {
	var item domain.ProductItem
	err := r.db.Where("sku_code = ?", skuCode).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByProductID retrieves all product items (SKUs) for a product
func (r *productItemRepository) GetByProductID(productID uint) ([]*domain.ProductItem, error) {
	var items []*domain.ProductItem
	err := r.db.Where("product_id = ?", productID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Delete deletes a product item
func (r *productItemRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ProductItem{}, id).Error
}

// UpdateStock updates the stock quantity atomically
func (r *productItemRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&domain.ProductItem{}).Where("id = ?", id).Update("qty_in_stock", quantity).Error
}

