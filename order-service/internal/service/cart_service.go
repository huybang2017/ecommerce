package service

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/domain"

	"go.uber.org/zap"
)

// CartService contains the business logic for cart operations
type CartService struct {
	cartRepo      domain.CartRepository
	productClient ProductServiceClient
	logger        *zap.Logger
}

// ProductServiceClient defines interface to communicate with Product Service
type ProductServiceClient interface {
	// GetProductItem fetches single product item details (SKU-level)
	GetProductItem(productItemID uint) (*ProductItemDTO, error)

	// GetProductItems fetches multiple product items in batch (for performance)
	GetProductItems(productItemIDs []uint) (map[uint]*ProductItemDTO, error)
}

// ProductItemDTO represents product item data from Product Service
// This is fetched on-demand, NOT stored in Redis cart
// NOTE: This is DISPLAY-ONLY for cart. Order validation uses full DTO with Stock/IsActive.
type ProductItemDTO struct {
	ID          uint    `json:"id"`           // ProductItem ID (SKU)
	ShopID      uint    `json:"shop_id"`      // Shop that owns this product
	ProductName string  `json:"product_name"` // Product name
	SKUCode     string  `json:"sku_code"`     // SKU code
	Price       float64 `json:"price"`        // Current price (for display only)
	ImageURL    string  `json:"image_url"`    // Product image
	QtyInStock  int     `json:"qty_in_stock"` // Stock quantity
	Status      string  `json:"status"`       // ACTIVE, INACTIVE
}

// NewCartService creates a new cart service
func NewCartService(
	cartRepo domain.CartRepository,
	productClient ProductServiceClient,
	logger *zap.Logger,
) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		productClient: productClient,
		logger:        logger,
	}
}

// GetCart retrieves user's cart and enriches with product data from Product Service
func (s *CartService) GetCart(ctx context.Context, userID string) (*domain.ShoppingCart, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	// 1. Get cart from Redis (only contains product_item_id, quantity, is_selected)
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		s.logger.Error("failed to get cart from Redis",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// 2. If cart is empty, return immediately
	if len(cart.Items) == 0 {
		return cart, nil
	}

	// 3. Fetch product details from Product Service
	if err := s.enrichCartWithProductData(cart); err != nil {
		s.logger.Warn("failed to enrich cart with product data",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		// Return cart anyway, just missing product details
	}

	// 4. Calculate cart-level totals only (no checkout grouping here)
	cart.CalculateTotals()

	return cart, nil
}

// AddToCart adds a product item (SKU) to cart
func (s *CartService) AddToCart(ctx context.Context, userID string, productItemID uint, quantity int) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	if productItemID == 0 {
		return domain.ErrInvalidProductItem
	}

	if quantity <= 0 {
		return domain.ErrInvalidQuantity
	}

	if quantity > 999 {
		return domain.ErrQuantityExceedsLimit
	}

	// 4. Get cart from Redis
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// 5. Check if item already exists
	existingItem := cart.FindItemByProductItemID(productItemID)

	if existingItem != nil {
		// Update quantity
		newQuantity := existingItem.Quantity + quantity

		if newQuantity > 999 {
			return domain.ErrQuantityExceedsLimit
		}

		existingItem.Quantity = newQuantity

	} else {
		// Add new item (only store minimal data in Redis)
		newItem := &domain.CartItem{
			ProductItemID: productItemID,
			Quantity:      quantity,
			IsSelected:    true, // Auto-select new items
		}

		if err := newItem.Validate(); err != nil {
			return err
		}

		cart.Items = append(cart.Items, newItem)
	}

	// 6. Save cart to Redis
	if err := s.cartRepo.SaveCart(cart); err != nil {
		s.logger.Error("failed to save cart to Redis",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to save cart: %w", err)
	}

	s.logger.Info("item added to cart",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", productItemID),
		zap.Int("quantity", quantity),
	)

	return nil
}

// UpdateItemQuantity updates quantity of a cart item
func (s *CartService) UpdateItemQuantity(ctx context.Context, userID string, productItemID uint, quantity int) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	// If quantity is 0, remove item
	if quantity == 0 {
		return s.RemoveFromCart(ctx, userID, productItemID)
	}

	if quantity < 0 {
		return domain.ErrInvalidQuantity
	}

	if quantity > 999 {
		return domain.ErrQuantityExceedsLimit
	}

	// Get cart
	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Find item
	item := cart.FindItemByProductItemID(productItemID)
	if item == nil {
		return domain.ErrCartItemNotFound
	}

	// Update quantity
	item.Quantity = quantity

	// Save cart
	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	s.logger.Info("cart item quantity updated",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", productItemID),
		zap.Int("new_quantity", quantity),
	)

	return nil
}

// RemoveFromCart removes an item from cart
func (s *CartService) RemoveFromCart(ctx context.Context, userID string, productItemID uint) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Find and remove item
	newItems := make([]*domain.CartItem, 0, len(cart.Items))
	found := false

	for _, item := range cart.Items {
		if item.ProductItemID == productItemID {
			found = true
			continue
		}
		newItems = append(newItems, item)
	}

	if !found {
		return domain.ErrCartItemNotFound
	}

	cart.Items = newItems

	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	s.logger.Info("item removed from cart",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", productItemID),
	)

	return nil
}

// ClearCart removes all items from cart
func (s *CartService) ClearCart(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	cart.Items = make([]*domain.CartItem, 0)

	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	s.logger.Info("cart cleared", zap.String("user_id", userID))

	return nil
}

// ClearSelectedItems removes only selected items (after checkout)
func (s *CartService) ClearSelectedItems(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Keep only unselected items
	unselectedItems := make([]*domain.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		if !item.IsSelected {
			unselectedItems = append(unselectedItems, item)
		}
	}

	cart.Items = unselectedItems

	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to clear selected items: %w", err)
	}

	s.logger.Info("selected items cleared",
		zap.String("user_id", userID),
		zap.Int("remaining_items", len(unselectedItems)),
	)

	return nil
}

// ToggleItemSelection toggles selection state of an item
func (s *CartService) ToggleItemSelection(ctx context.Context, userID string, productItemID uint) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	item := cart.FindItemByProductItemID(productItemID)
	if item == nil {
		return domain.ErrCartItemNotFound
	}

	item.IsSelected = !item.IsSelected

	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	s.logger.Info("item selection toggled",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", productItemID),
		zap.Bool("is_selected", item.IsSelected),
	)

	return nil
}

// SelectAllItems selects/deselects all items
func (s *CartService) SelectAllItems(ctx context.Context, userID string, selected bool) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	for i := range cart.Items {
		cart.Items[i].IsSelected = selected
	}

	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	s.logger.Info("all items selection updated",
		zap.String("user_id", userID),
		zap.Bool("selected", selected),
	)

	return nil
}

// SelectShopItems selects/deselects all items from a specific shop
func (s *CartService) SelectShopItems(ctx context.Context, userID string, shopID uint, selected bool) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.cartRepo.GetCart(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	// Update selection for items from this shop
	for i := range cart.Items {
		if cart.Items[i].ShopID == shopID {
			cart.Items[i].IsSelected = selected
		}
	}

	if err := s.cartRepo.SaveCart(cart); err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	s.logger.Info("shop items selection updated",
		zap.String("user_id", userID),
		zap.Uint("shop_id", shopID),
		zap.Bool("selected", selected),
	)

	return nil
}

// ValidateCart validates all items in cart
func (s *CartService) ValidateCart(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user_id is required")
	}

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	if cart.IsEmpty() {
		return domain.ErrCartEmpty
	}

	return s.validateSelectedItems(cart)
}

// enrichCartWithProductData fetches product details from Product Service
func (s *CartService) enrichCartWithProductData(cart *domain.ShoppingCart) error {
	if len(cart.Items) == 0 {
		return nil
	}

	// Collect all product item IDs
	productItemIDs := make([]uint, 0, len(cart.Items))
	for _, item := range cart.Items {
		productItemIDs = append(productItemIDs, item.ProductItemID)
	}

	// Batch fetch from Product Service
	productItems, err := s.productClient.GetProductItems(productItemIDs)
	if err != nil {
		s.logger.Error("failed to fetch product items from Product Service",
			zap.Uints("product_item_ids", productItemIDs),
			zap.Error(err),
		)
		return fmt.Errorf("failed to fetch product items: %w", err)
	}

	s.logger.Info("fetched product items from Product Service",
		zap.Int("requested_count", len(productItemIDs)),
		zap.Int("received_count", len(productItems)),
	)

	// Enrich cart items with product data
	for _, item := range cart.Items {
		if productItem, ok := productItems[item.ProductItemID]; ok {
			item.ShopID = productItem.ShopID
			item.ProductName = productItem.ProductName
			item.SKUCode = productItem.SKUCode
			item.Price = productItem.Price
			item.ImageURL = productItem.ImageURL
			s.logger.Debug("enriched cart item",
				zap.Uint("product_item_id", item.ProductItemID),
				zap.Uint("shop_id", item.ShopID),
				zap.String("product_name", item.ProductName),
			)
		} else {
			s.logger.Warn("product item not found in Product Service response",
				zap.Uint("product_item_id", item.ProductItemID),
			)
		}
	}

	return nil
}

// validateSelectedItems validates all selected items in the cart
func (s *CartService) validateSelectedItems(cart *domain.ShoppingCart) error {
	// Collect all product item IDs from selected items
	productItemIDs := make([]uint, 0, len(cart.Items))
	for _, item := range cart.Items {
		if item.IsSelected {
			productItemIDs = append(productItemIDs, item.ProductItemID)
		}
	}

	if len(productItemIDs) == 0 {
		return nil // No selected items, nothing to validate
	}

	return nil
}
