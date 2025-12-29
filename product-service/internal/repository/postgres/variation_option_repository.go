package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// variationOptionRepository implements the VariationOptionRepository interface
type variationOptionRepository struct {
	db *gorm.DB
}

// NewVariationOptionRepository creates a new PostgreSQL variation option repository
func NewVariationOptionRepository(db *gorm.DB) domain.VariationOptionRepository {
	return &variationOptionRepository{db: db}
}

// Create inserts a new variation option into the database
func (r *variationOptionRepository) Create(option *domain.VariationOption) error {
	return r.db.Create(option).Error
}

// Update updates an existing variation option
func (r *variationOptionRepository) Update(option *domain.VariationOption) error {
	return r.db.Save(option).Error
}

// GetByID retrieves a variation option by its ID
func (r *variationOptionRepository) GetByID(id uint) (*domain.VariationOption, error) {
	var option domain.VariationOption
	err := r.db.First(&option, id).Error
	if err != nil {
		return nil, err
	}
	return &option, nil
}

// GetByVariationID retrieves all options for a variation
func (r *variationOptionRepository) GetByVariationID(variationID uint) ([]*domain.VariationOption, error) {
	var options []*domain.VariationOption
	err := r.db.Where("variation_id = ?", variationID).Find(&options).Error
	if err != nil {
		return nil, err
	}
	return options, nil
}

// Delete deletes a variation option
func (r *variationOptionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.VariationOption{}, id).Error
}

