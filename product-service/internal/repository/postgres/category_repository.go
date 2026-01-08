package postgres

import (
	"fmt"
	"product-service/internal/domain"

	"gorm.io/gorm"
)

// categoryRepository implements the CategoryRepository interface
// This is the infrastructure layer - it knows HOW to interact with PostgreSQL
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new PostgreSQL category repository
// Dependency injection: we inject the database connection
func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

// Create inserts a new category into the database
func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

// Update updates an existing category
func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

// GetByID retrieves a category by its ID
func (r *categoryRepository) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Parent").Preload("Children").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	// Debug: check if parent loaded
	if category.ParentID != nil {
		fmt.Printf("[DEBUG] Category %d has parent_id=%d, Parent loaded: %v\n",
			category.ID, *category.ParentID, category.Parent != nil)
	}
	return &category, nil
}

// GetBySlug retrieves a category by its slug
func (r *categoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Parent").Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAll retrieves all categories
func (r *categoryRepository) GetAll() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetChildren retrieves all child categories of a parent category
func (r *categoryRepository) GetChildren(parentID uint) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Delete deletes a category (hard delete)
// Note: In production, you might want to check if category has products before deleting
func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
