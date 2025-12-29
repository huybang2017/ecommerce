package postgres

import (
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// categoryAttributeRepository implements the CategoryAttributeRepository interface
type categoryAttributeRepository struct {
	db *gorm.DB
}

// NewCategoryAttributeRepository creates a new PostgreSQL category attribute repository
func NewCategoryAttributeRepository(db *gorm.DB) domain.CategoryAttributeRepository {
	return &categoryAttributeRepository{db: db}
}

// Create inserts a new category attribute into the database
func (r *categoryAttributeRepository) Create(attr *domain.CategoryAttribute) error {
	return r.db.Create(attr).Error
}

// Update updates an existing category attribute
func (r *categoryAttributeRepository) Update(attr *domain.CategoryAttribute) error {
	return r.db.Save(attr).Error
}

// GetByID retrieves a category attribute by its ID
func (r *categoryAttributeRepository) GetByID(id uint) (*domain.CategoryAttribute, error) {
	var attr domain.CategoryAttribute
	err := r.db.First(&attr, id).Error
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// GetByCategoryID retrieves all attributes for a category
func (r *categoryAttributeRepository) GetByCategoryID(categoryID uint) ([]*domain.CategoryAttribute, error) {
	var attrs []*domain.CategoryAttribute
	err := r.db.Where("category_id = ?", categoryID).Find(&attrs).Error
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// GetFilterablesByCategoryID retrieves only filterable attributes for a category
func (r *categoryAttributeRepository) GetFilterablesByCategoryID(categoryID uint) ([]*domain.CategoryAttribute, error) {
	var attrs []*domain.CategoryAttribute
	err := r.db.Where("category_id = ? AND is_filterable = ?", categoryID, true).Find(&attrs).Error
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// Delete deletes a category attribute
func (r *categoryAttributeRepository) Delete(id uint) error {
	return r.db.Delete(&domain.CategoryAttribute{}, id).Error
}

