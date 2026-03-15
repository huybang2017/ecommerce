package postgres

import (
	"context"
	"fmt"
	"product-service/internal/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) GetAdminList(ctx context.Context, p domain.CategoryQueryParams) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Category{})

	if p.Search != "" {
		searchTerm := "%" + p.Search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", searchTerm, searchTerm)
	}

	if p.ParentID != nil {
		query = query.Where("parent_id = ?", *p.ParentID)
	}

	if p.IsActive != nil {
		query = query.Where("is_active = ?", *p.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", p.SortBy, p.SortOrder)
	query = query.Order(orderClause)

	offset := (p.Page - 1) * p.Limit

	err := query.Limit(p.Limit).Offset(offset).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) GetChildren(parentID uint) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *categoryRepository) GetAdminCategoryParent(ctx context.Context) ([]*domain.Category, error) {
	var parents []*domain.Category
	err := r.db.WithContext(ctx).Where("parent_id IS NULL").Find(&parents).Error
	if err != nil {
		return nil, err
	}
	return parents, nil
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}
