package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// variationRepository implements the VariationRepository interface
type variationRepository struct {
	db *gorm.DB
}

// NewVariationRepository creates a new PostgreSQL variation repository
func NewVariationRepository(db *gorm.DB) domain.VariationRepository {
	return &variationRepository{db: db}
}

// Create inserts a new variation into the database
func (r *variationRepository) Create(variation *domain.Variation) error {
	return r.db.Create(variation).Error
}

// Update updates an existing variation
func (r *variationRepository) Update(variation *domain.Variation) error {
	return r.db.Save(variation).Error
}

// GetByID retrieves a variation by its ID
func (r *variationRepository) GetByID(id uint) (*domain.Variation, error) {
	var variation domain.Variation
	err := r.db.First(&variation, id).Error
	if err != nil {
		return nil, err
	}
	return &variation, nil
}

// GetByProductID retrieves all variations for a product
func (r *variationRepository) GetByProductID(productID uint) ([]*domain.Variation, error) {
	var variations []*domain.Variation
	err := r.db.Where("product_id = ?", productID).Find(&variations).Error
	if err != nil {
		return nil, err
	}
	return variations, nil
}

// Delete deletes a variation
func (r *variationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Variation{}, id).Error
}

