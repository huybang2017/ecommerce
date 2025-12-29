package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// productAttributeValueRepository implements the ProductAttributeValueRepository interface
type productAttributeValueRepository struct {
	db *gorm.DB
}

// NewProductAttributeValueRepository creates a new PostgreSQL product attribute value repository
func NewProductAttributeValueRepository(db *gorm.DB) domain.ProductAttributeValueRepository {
	return &productAttributeValueRepository{db: db}
}

// Create inserts a new product attribute value into the database
func (r *productAttributeValueRepository) Create(value *domain.ProductAttributeValue) error {
	return r.db.Create(value).Error
}

// CreateBatch inserts multiple product attribute values in a single transaction
func (r *productAttributeValueRepository) CreateBatch(values []*domain.ProductAttributeValue) error {
	return r.db.Create(values).Error
}

// Update updates an existing product attribute value
func (r *productAttributeValueRepository) Update(value *domain.ProductAttributeValue) error {
	return r.db.Save(value).Error
}

// GetByID retrieves a product attribute value by its ID
func (r *productAttributeValueRepository) GetByID(id uint) (*domain.ProductAttributeValue, error) {
	var value domain.ProductAttributeValue
	err := r.db.First(&value, id).Error
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// GetByProductID retrieves all attribute values for a product
func (r *productAttributeValueRepository) GetByProductID(productID uint) ([]*domain.ProductAttributeValue, error) {
	var values []*domain.ProductAttributeValue
	err := r.db.Where("product_id = ?", productID).Find(&values).Error
	if err != nil {
		return nil, err
	}
	return values, nil
}

// GetByAttributeID retrieves all values for a specific attribute
func (r *productAttributeValueRepository) GetByAttributeID(attributeID uint) ([]*domain.ProductAttributeValue, error) {
	var values []*domain.ProductAttributeValue
	err := r.db.Where("attribute_id = ?", attributeID).Find(&values).Error
	if err != nil {
		return nil, err
	}
	return values, nil
}

// SearchByAttributeValue searches for products by attribute value
// This uses the compound index (attribute_id, value) for fast search
func (r *productAttributeValueRepository) SearchByAttributeValue(attributeID uint, value string) ([]*domain.ProductAttributeValue, error) {
	var values []*domain.ProductAttributeValue
	err := r.db.Where("attribute_id = ? AND value = ?", attributeID, value).Find(&values).Error
	if err != nil {
		return nil, err
	}
	return values, nil
}

// Delete deletes a product attribute value
func (r *productAttributeValueRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ProductAttributeValue{}, id).Error
}

// DeleteByProductID deletes all attribute values for a product
func (r *productAttributeValueRepository) DeleteByProductID(productID uint) error {
	return r.db.Where("product_id = ?", productID).Delete(&domain.ProductAttributeValue{}).Error
}

