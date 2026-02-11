package postgres

import (
	"context"
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

// GetByID retrieves a category by its ID (NO Preload to avoid N+1)
func (r *categoryRepository) GetByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug retrieves a category by its slug (NO Preload)
func (r *categoryRepository) GetBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
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

func (r *categoryRepository) GetAdminList(ctx context.Context, p domain.CategoryQueryParams) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var total int64

	// 1. Khởi tạo query từ model Category
	// Dùng WithContext để quản lý timeout/cancel từ phía Service
	query := r.db.WithContext(ctx).Model(&domain.Category{})

	// 2. Xử lý FILTER: Tìm kiếm theo Tên hoặc Slug
	if p.Search != "" {
		// ILIKE là toán tử PostgreSQL tìm kiếm không phân biệt hoa thường
		searchTerm := "%" + p.Search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", searchTerm, searchTerm)
	}

	// 3. Xử lý FILTER: Trạng thái Active/Inactive
	if p.Status == "active" {
		query = query.Where("is_active = ?", true)
	} else if p.Status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	// 4. COUNT: Đếm tổng số bản ghi TRƯỚC khi thực hiện Limit/Offset
	// Đây là con số quan trọng để Frontend hiển thị: "1-10 of 500 rows"
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 5. Xử lý SORT: Sắp xếp động
	// Lưu ý: Cần kiểm tra SortBy để tránh SQL Injection nếu param lấy trực tiếp từ URL
	orderClause := fmt.Sprintf("%s %s", p.SortBy, p.SortOrder)
	query = query.Order(orderClause)

	// 6. Xử lý PAGINATION: Phân trang
	offset := (p.Page - 1) * p.Limit

	// 7. EXECUTE: Lấy dữ liệu cuối cùng
	err := query.Limit(p.Limit).Offset(offset).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
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
