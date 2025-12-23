package service

import (
	"context"
	"errors"
	"fmt"
	"product-service/internal/domain"
	"strings"

	"go.uber.org/zap"
)

// CategoryService contains the business logic for category operations
// This is the service layer - it orchestrates between repositories
type CategoryService struct {
	categoryRepo domain.CategoryRepository
	logger       *zap.Logger
}

// NewCategoryService creates a new category service with all dependencies
func NewCategoryService(
	categoryRepo domain.CategoryRepository,
	logger *zap.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	// Business logic validation
	if category.Name == "" {
		return errors.New("category name is required")
	}

	// Generate slug from name if not provided
	if category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}

	// Check if slug already exists
	existing, err := s.categoryRepo.GetBySlug(category.Slug)
	if err == nil && existing != nil {
		return errors.New("category with this slug already exists")
	}

	// Validate parent_id if provided
	if category.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(*category.ParentID)
		if err != nil {
			return errors.New("parent category not found")
		}
		if parent == nil {
			return errors.New("parent category not found")
		}
	}

	// Create category
	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("failed to create category in database", zap.Error(err))
		return fmt.Errorf("failed to create category: %w", err)
	}

	s.logger.Info("category created", zap.Uint("category_id", category.ID))
	return nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	// Validate category exists
	existing, err := s.categoryRepo.GetByID(category.ID)
	if err != nil {
		return errors.New("category not found")
	}

	// Generate slug from name if name changed and slug not provided
	if category.Name != existing.Name && category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}

	// Check if slug already exists (excluding current category)
	if category.Slug != existing.Slug {
		existingBySlug, err := s.categoryRepo.GetBySlug(category.Slug)
		if err == nil && existingBySlug != nil && existingBySlug.ID != category.ID {
			return errors.New("category with this slug already exists")
		}
	}

	// Validate parent_id if provided (prevent circular reference)
	if category.ParentID != nil {
		if *category.ParentID == category.ID {
			return errors.New("category cannot be its own parent")
		}
		parent, err := s.categoryRepo.GetByID(*category.ParentID)
		if err != nil || parent == nil {
			return errors.New("parent category not found")
		}
	}

	// Preserve created_at
	category.CreatedAt = existing.CreatedAt

	// Update category
	if err := s.categoryRepo.Update(category); err != nil {
		s.logger.Error("failed to update category in database", zap.Error(err))
		return fmt.Errorf("failed to update category: %w", err)
	}

	s.logger.Info("category updated", zap.Uint("category_id", category.ID))
	return nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(ctx context.Context, id uint) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all categories", zap.Error(err))
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}

// GetCategoryChildren retrieves child categories of a parent category
func (s *CategoryService) GetCategoryChildren(ctx context.Context, parentID uint) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetChildren(parentID)
	if err != nil {
		s.logger.Error("failed to get category children", zap.Error(err))
		return nil, fmt.Errorf("failed to get category children: %w", err)
	}
	return categories, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// Check if category has children
	children, err := s.categoryRepo.GetChildren(id)
	if err == nil && len(children) > 0 {
		return errors.New("cannot delete category with children")
	}

	// Delete category
	if err := s.categoryRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete category", zap.Error(err))
		return fmt.Errorf("failed to delete category: %w", err)
	}

	s.logger.Info("category deleted", zap.Uint("category_id", id))
	return nil
}

// generateSlug generates a URL-friendly slug from a name
func (s *CategoryService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

