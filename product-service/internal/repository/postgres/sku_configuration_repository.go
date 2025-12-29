package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// skuConfigurationRepository implements the SKUConfigurationRepository interface
type skuConfigurationRepository struct {
	db *gorm.DB
}

// NewSKUConfigurationRepository creates a new PostgreSQL SKU configuration repository
func NewSKUConfigurationRepository(db *gorm.DB) domain.SKUConfigurationRepository {
	return &skuConfigurationRepository{db: db}
}

// Create inserts a new SKU configuration into the database
func (r *skuConfigurationRepository) Create(config *domain.SKUConfiguration) error {
	return r.db.Create(config).Error
}

// CreateBatch inserts multiple SKU configurations in a single transaction
func (r *skuConfigurationRepository) CreateBatch(configs []*domain.SKUConfiguration) error {
	return r.db.Create(configs).Error
}

// GetByProductItemID retrieves all configurations for a product item (SKU)
func (r *skuConfigurationRepository) GetByProductItemID(productItemID uint) ([]*domain.SKUConfiguration, error) {
	var configs []*domain.SKUConfiguration
	err := r.db.Where("product_item_id = ?", productItemID).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// GetByVariationOptionID retrieves all configurations for a variation option
func (r *skuConfigurationRepository) GetByVariationOptionID(optionID uint) ([]*domain.SKUConfiguration, error) {
	var configs []*domain.SKUConfiguration
	err := r.db.Where("variation_option_id = ?", optionID).Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// Delete deletes a specific SKU configuration
func (r *skuConfigurationRepository) Delete(productItemID uint, variationOptionID uint) error {
	return r.db.Where("product_item_id = ? AND variation_option_id = ?", productItemID, variationOptionID).
		Delete(&domain.SKUConfiguration{}).Error
}

// DeleteByProductItemID deletes all configurations for a product item (SKU)
func (r *skuConfigurationRepository) DeleteByProductItemID(productItemID uint) error {
	return r.db.Where("product_item_id = ?", productItemID).Delete(&domain.SKUConfiguration{}).Error
}

