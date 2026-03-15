package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"product-service/internal/domain"
	"strings"

	"go.uber.org/zap"
)

type CategoryService struct {
	categoryRepo domain.CategoryRepository
	logger       *zap.Logger
}

func NewCategoryService(
	categoryRepo domain.CategoryRepository,
	logger *zap.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s *CategoryService) CreateParentCategory(ctx context.Context, category *domain.Category) error {
	if category.ParentID != nil {
		return errors.New("parent category cannot have a parent_id")
	}
	if category.Name == "" {
		return errors.New("category name is required")
	}
	if category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}
	existing, err := s.categoryRepo.GetBySlug(category.Slug)
	if err == nil && existing != nil {
		return errors.New("category with this slug already exists")
	}
	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("failed to create parent category in database", zap.Error(err))
		return fmt.Errorf("failed to create parent category: %w", err)
	}

	s.logger.Info("parent category created", zap.Uint("category_id", category.ID))
	return nil
}

func (s *CategoryService) CreateChildCategory(ctx context.Context, category *domain.Category) error {
	if category.ParentID == nil {
		return errors.New("child category must have a parent_id")
	}
	if category.Name == "" {
		return errors.New("category name is required")
	}
	parent, err := s.categoryRepo.GetByID(*category.ParentID)
	if err != nil || parent == nil {
		return errors.New("parent category not found")
	}
	if category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}
	existing, err := s.categoryRepo.GetBySlug(category.Slug)
	if err == nil && existing != nil {
		return errors.New("category with this slug already exists")
	}
	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("failed to create child category in database", zap.Error(err))
		return fmt.Errorf("failed to create child category: %w", err)
	}

	s.logger.Info("child category created",
		zap.Uint("category_id", category.ID),
		zap.Uint("parent_id", *category.ParentID),
	)
	return nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}

	if category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}

	existing, err := s.categoryRepo.GetBySlug(category.Slug)
	if err == nil && existing != nil {
		return errors.New("category with this slug already exists")
	}

	if category.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(*category.ParentID)
		if err != nil {
			return errors.New("parent category not found")
		}
		if parent == nil {
			return errors.New("parent category not found")
		}
	}

	if err := s.categoryRepo.Create(category); err != nil {
		s.logger.Error("failed to create category in database", zap.Error(err))
		return fmt.Errorf("failed to create category: %w", err)
	}

	s.logger.Info("category created", zap.Uint("category_id", category.ID))
	return nil
}

func (s *CategoryService) IsActiveCategory(ctx context.Context, id uint, isActive bool) error {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to find category", zap.Uint("category_id", id), zap.Error(err))
		return fmt.Errorf("failed to find category: %w", err)
	}
	if category == nil {
		return errors.New("category not found")
	}

	category.IsActive = isActive

	if err := s.categoryRepo.Update(category); err != nil {
		s.logger.Error("failed to update category status in database",
			zap.Uint("category_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update category status: %w", err)
	}

	s.logger.Info("category status patched",
		zap.Uint("category_id", id),
		zap.Bool("is_active", isActive),
	)

	return nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	existing, err := s.categoryRepo.GetByID(category.ID)
	if err != nil {
		return errors.New("category not found")
	}

	if category.Name != existing.Name && category.Slug == "" {
		category.Slug = s.generateSlug(category.Name)
	}

	if category.Slug != existing.Slug {
		existingBySlug, err := s.categoryRepo.GetBySlug(category.Slug)
		if err == nil && existingBySlug != nil && existingBySlug.ID != category.ID {
			return errors.New("category with this slug already exists")
		}
	}

	if category.ParentID != nil {
		if *category.ParentID == category.ID {
			return errors.New("category cannot be its own parent")
		}
		parent, err := s.categoryRepo.GetByID(*category.ParentID)
		if err != nil || parent == nil {
			return errors.New("parent category not found")
		}
	}

	category.CreatedAt = existing.CreatedAt

	if err := s.categoryRepo.Update(category); err != nil {
		s.logger.Error("failed to update category in database", zap.Error(err))
		return fmt.Errorf("failed to update category: %w", err)
	}

	s.logger.Info("category updated", zap.Uint("category_id", category.ID))
	return nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id uint) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

func (s *CategoryService) GetAdminCategories(ctx context.Context, params domain.CategoryQueryParams) ([]*domain.Category, int64, int, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.SortBy == "" {
		params.SortBy = "id"
	}
	if params.SortOrder == "" {
		params.SortOrder = "DESC"
	}

	categories, total, err := s.categoryRepo.GetAdminList(ctx, params)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to fetch admin categories: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return categories, total, totalPages, nil
}

func (s *CategoryService) GetAdminCategoryParent(ctx context.Context) ([]*domain.Category, error) {
	parents, err := s.categoryRepo.GetAdminCategoryParent(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin category parents: %w", err)
	}
	return parents, nil
}

func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}
	return category, nil
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all categories", zap.Error(err))
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}

func (s *CategoryService) GetCategoryChildren(ctx context.Context, parentID uint) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetChildren(parentID)
	if err != nil {
		s.logger.Error("failed to get category children", zap.Error(err))
		return nil, fmt.Errorf("failed to get category children: %w", err)
	}
	return categories, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	children, err := s.categoryRepo.GetChildren(id)
	if err == nil && len(children) > 0 {
		return errors.New("cannot delete category with children")
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete category", zap.Error(err))
		return fmt.Errorf("failed to delete category: %w", err)
	}

	s.logger.Info("category deleted", zap.Uint("category_id", id))
	return nil
}

func (s *CategoryService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
