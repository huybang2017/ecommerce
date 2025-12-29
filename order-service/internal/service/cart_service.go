package service

import (
	"context"
	"errors"
	"order-service/internal/domain"

	"go.uber.org/zap"
)

// CartService contains the business logic for cart operations
// This is the service layer - it orchestrates between repositories
type CartService struct {
	cartRepo domain.CartRepository
	logger   *zap.Logger
}

// NewCartService creates a new cart service with dependencies
func NewCartService(cartRepo domain.CartRepository, logger *zap.Logger) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		logger:   logger,
	}
}

// GetCart retrieves a cart for a user
// Business rule: Cart requires authentication - only user_id is accepted
func (s *CartService) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	if userID == "" {
		return nil, errors.New("user_id is required - authentication required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		s.logger.Error("failed to get cart", zap.Error(err))
		return nil, err
	}

	return cart, nil
}

// AddItem adds a product to the cart
// Business rule: Cart requires authentication - only user_id is accepted
func (s *CartService) AddItem(ctx context.Context, userID string, productID uint, name string, price float64, quantity int, image, sku string) (*domain.Cart, error) {
	if userID == "" {
		return nil, errors.New("user_id is required - authentication required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	// Get existing cart
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return nil, err
	}

	// Initialize items map if nil
	if cart.Items == nil {
		cart.Items = make(map[uint]*domain.CartItem)
	}

	// Update or add item
	if existingItem, exists := cart.Items[productID]; exists {
		// Update quantity
		existingItem.Quantity += quantity
	} else {
		// Add new item
		cart.Items[productID] = &domain.CartItem{
			ProductID: productID,
			Name:      name,
			Price:     price,
			Quantity:  quantity,
			Image:     image,
			SKU:       sku,
		}
	}

	// Save cart
	err = s.cartRepo.SaveCart(cart)
	if err != nil {
		s.logger.Error("failed to save cart", zap.Error(err))
		return nil, err
	}

	return cart, nil
}

// UpdateItemQuantity updates the quantity of an item in the cart
// Business rule: Cart requires authentication - only user_id is accepted
func (s *CartService) UpdateItemQuantity(ctx context.Context, userID string, productID uint, quantity int) (*domain.Cart, error) {
	if userID == "" {
		return nil, errors.New("user_id is required - authentication required")
	}
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	// Get existing cart
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return nil, err
	}

	// Check if item exists
	item, exists := cart.Items[productID]
	if !exists {
		return nil, errors.New("item not found in cart")
	}

	if quantity == 0 {
		// Remove item if quantity is 0
		delete(cart.Items, productID)
	} else {
		// Update quantity
		item.Quantity = quantity
	}

	// Save cart
	err = s.cartRepo.SaveCart(cart)
	if err != nil {
		s.logger.Error("failed to save cart", zap.Error(err))
		return nil, err
	}

	return cart, nil
}

// RemoveItem removes an item from the cart
// Business rule: Cart requires authentication - only user_id is accepted
func (s *CartService) RemoveItem(ctx context.Context, userID string, productID uint) (*domain.Cart, error) {
	if userID == "" {
		return nil, errors.New("user_id is required - authentication required")
	}

	// Get existing cart
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return nil, err
	}

	// Check if item exists
	if _, exists := cart.Items[productID]; !exists {
		return nil, errors.New("item not found in cart")
	}

	// Remove item
	delete(cart.Items, productID)

	// Save cart
	err = s.cartRepo.SaveCart(cart)
	if err != nil {
		s.logger.Error("failed to save cart", zap.Error(err))
		return nil, err
	}

	return cart, nil
}

// ClearCart removes all items from the cart
// Business rule: Cart requires authentication - only user_id is accepted
func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id is required - authentication required")
	}

	err := s.cartRepo.ClearCartItems(userID)
	if err != nil {
		s.logger.Error("failed to clear cart", zap.Error(err))
		return err
	}

	return nil
}


