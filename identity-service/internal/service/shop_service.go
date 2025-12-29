package service

import (
	"errors"
	"fmt"
	"identity-service/internal/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ShopService contains the business logic for shop operations
// Following Clean Architecture: business logic is independent of infrastructure
type ShopService struct {
	shopRepo domain.ShopRepository
	userRepo domain.UserRepository
	logger   *zap.Logger
}

// NewShopService creates a new shop service
func NewShopService(
	shopRepo domain.ShopRepository,
	userRepo domain.UserRepository,
	logger *zap.Logger,
) *ShopService {
	return &ShopService{
		shopRepo: shopRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateShopRequest represents the request to create a new shop
type CreateShopRequest struct {
	OwnerUserID  uint   `json:"owner_user_id" binding:"required"`
	Name         string `json:"name" binding:"required,min=3,max=100"`
	Description  string `json:"description"`
	LogoURL      string `json:"logo_url"`
	CoverURL     string `json:"cover_url"`
}

// UpdateShopRequest represents the request to update a shop
type UpdateShopRequest struct {
	Name         string `json:"name" binding:"omitempty,min=3,max=100"`
	Description  string `json:"description"`
	LogoURL      string `json:"logo_url"`
	CoverURL     string `json:"cover_url"`
}

// CreateShop creates a new shop
// Business rules:
// - 1 User can only have 1 Shop (unique constraint on owner_user_id)
// - Only SELLER role can create shop
// - User must exist and be active
func (s *ShopService) CreateShop(req *CreateShopRequest) (*domain.Shop, error) {
	// Validate user exists and is active
	user, err := s.userRepo.GetByID(req.OwnerUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check user status
	if user.Status != "ACTIVE" {
		return nil, errors.New("user is not active")
	}

	// Check user role (only SELLER can create shop)
	if user.Role != "SELLER" && user.Role != "ADMIN" {
		return nil, errors.New("only SELLER or ADMIN can create shop")
	}

	// Check if user already has a shop (1 User = 1 Shop)
	existingShop, err := s.shopRepo.GetByOwnerUserID(req.OwnerUserID)
	if err == nil && existingShop != nil {
		return nil, errors.New("user already has a shop")
	}

	// Create shop
	shop := &domain.Shop{
		OwnerUserID:  req.OwnerUserID,
		Name:         req.Name,
		Description:  req.Description,
		LogoURL:      req.LogoURL,
		CoverURL:     req.CoverURL,
		IsOfficial:   false,
		Rating:       0,
		ResponseRate: 0,
		Status:       "ACTIVE",
	}

	if err := s.shopRepo.Create(shop); err != nil {
		s.logger.Error("failed to create shop", zap.Error(err))
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}

	s.logger.Info("shop created", zap.Uint("shop_id", shop.ID), zap.Uint("owner_user_id", shop.OwnerUserID))

	return shop, nil
}

// UpdateShop updates an existing shop
// Business rule: Only shop owner or ADMIN can update
func (s *ShopService) UpdateShop(shopID uint, ownerUserID uint, req *UpdateShopRequest) (*domain.Shop, error) {
	// Get existing shop
	shop, err := s.shopRepo.GetByID(shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	// Validate ownership (only owner or ADMIN can update)
	user, err := s.userRepo.GetByID(ownerUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if shop.OwnerUserID != ownerUserID && user.Role != "ADMIN" {
		return nil, errors.New("only shop owner or ADMIN can update shop")
	}

	// Update fields
	if req.Name != "" {
		shop.Name = req.Name
	}
	if req.Description != "" {
		shop.Description = req.Description
	}
	if req.LogoURL != "" {
		shop.LogoURL = req.LogoURL
	}
	if req.CoverURL != "" {
		shop.CoverURL = req.CoverURL
	}

	if err := s.shopRepo.Update(shop); err != nil {
		s.logger.Error("failed to update shop", zap.Error(err))
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	s.logger.Info("shop updated", zap.Uint("shop_id", shop.ID))

	return shop, nil
}

// GetShop retrieves a shop by ID
func (s *ShopService) GetShop(id uint) (*domain.Shop, error) {
	shop, err := s.shopRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return shop, nil
}

// GetMyShop retrieves the shop of the current user (1 User = 1 Shop)
func (s *ShopService) GetMyShop(userID uint) (*domain.Shop, error) {
	shop, err := s.shopRepo.GetByOwnerUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not have a shop")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}
	return shop, nil
}

// ListShops retrieves all shops with pagination
func (s *ShopService) ListShops(page, limit int) ([]*domain.Shop, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	shops, total, err := s.shopRepo.GetAll(page, limit)
	if err != nil {
		s.logger.Error("failed to list shops", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list shops: %w", err)
	}

	return shops, total, nil
}

// DeleteShop soft deletes a shop (sets status to SUSPENDED)
// Business rule: Only ADMIN can delete shop
func (s *ShopService) DeleteShop(shopID uint, userID uint) error {
	// Validate user is ADMIN
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role != "ADMIN" {
		return errors.New("only ADMIN can delete shop")
	}

	// Soft delete (set status to SUSPENDED)
	if err := s.shopRepo.Delete(shopID); err != nil {
		s.logger.Error("failed to delete shop", zap.Error(err))
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	s.logger.Info("shop deleted", zap.Uint("shop_id", shopID), zap.Uint("deleted_by", userID))

	return nil
}

// UpdateShopStatus updates the status of a shop
// Business rule: Only ADMIN can update status
func (s *ShopService) UpdateShopStatus(shopID uint, status string, userID uint) error {
	// Validate user is ADMIN
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.Role != "ADMIN" {
		return errors.New("only ADMIN can update shop status")
	}

	// Validate status
	if status != "ACTIVE" && status != "SUSPENDED" {
		return errors.New("invalid status: must be ACTIVE or SUSPENDED")
	}

	if err := s.shopRepo.UpdateStatus(shopID, status); err != nil {
		s.logger.Error("failed to update shop status", zap.Error(err))
		return fmt.Errorf("failed to update shop status: %w", err)
	}

	s.logger.Info("shop status updated", zap.Uint("shop_id", shopID), zap.String("status", status))

	return nil
}

