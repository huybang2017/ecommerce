package service

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"order-service/pkg/product_client"
	"time"

	"go.uber.org/zap"
)

// CartService contains the business logic for cart operations
// This is the service layer - it orchestrates between repositories
type CartService struct {
	cartRepo       domain.CartRepository
	productClient  ProductClientInterface // THÊM MỚI - For marketplace: get shop_id
	logger         *zap.Logger
}

// ProductClientInterface defines interface for Product Service client
// This allows for easier testing and dependency injection
type ProductClientInterface interface {
	GetProductByID(productID uint) (*ProductInfo, error)
}

// ProductInfo represents product information needed for cart
type ProductInfo struct {
	ID     uint
	ShopID uint
	Name   string
	Price  float64
}

// NewCartService creates a new cart service with dependencies
func NewCartService(cartRepo domain.CartRepository, productClient ProductClientInterface, logger *zap.Logger) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		productClient: productClient,
		logger:        logger,
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
// Marketplace: Fetches shop_id from Product Service
func (s *CartService) AddItem(ctx context.Context, userID string, productID uint, name string, price float64, quantity int, image, sku string, productItemID uint) (*domain.Cart, error) {
	if userID == "" {
		return nil, errors.New("user_id is required - authentication required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	// MARKETPLACE: Get shop_id from Product Service
	var shopID uint
	if s.productClient != nil {
		product, err := s.productClient.GetProductByID(productID)
		if err != nil {
			s.logger.Warn("failed to get product info, using default shop_id", zap.Uint("product_id", productID), zap.Error(err))
			shopID = 1 // Fallback to default shop
		} else {
			shopID = product.ShopID
		}
	} else {
		s.logger.Warn("product client not available, using default shop_id")
		shopID = 1 // Fallback
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
		// Update shop_id if not set (backward compatibility)
		if existingItem.ShopID == 0 {
			existingItem.ShopID = shopID
		}
	} else {
		// Add new item with shop_id (marketplace)
		cart.Items[productID] = &domain.CartItem{
			ProductID:     productID,
			ProductItemID: productItemID, // SKU ID
			ShopID:        shopID,         // THÊM MỚI - Shop ID from Product Service
			Name:          name,
			Price:         price,
			Quantity:      quantity,
			Image:         image,
			SKU:           sku,
		}
	}

	// Recalculate total
	total := float64(0)
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}
	cart.Total = total
	cart.UpdatedAt = time.Now().Unix()

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

// ProductClientAdapter adapts product_client.ProductClient to ProductClientInterface
type ProductClientAdapter struct {
	Client *product_client.ProductClient
}

// GetProductByID implements ProductClientInterface
func (a *ProductClientAdapter) GetProductByID(productID uint) (*ProductInfo, error) {
	product, err := a.Client.GetProductByIDInternal(productID)
	if err != nil {
		return nil, err
	}
	return &ProductInfo{
		ID:     product.ID,
		ShopID: product.ShopID,
		Name:   product.Name,
		Price:  product.BasePrice,
	}, nil
}
