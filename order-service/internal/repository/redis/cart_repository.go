package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type cartRepository struct {
	client *redis.Client
	logger *zap.Logger
}

func NewCartRepository(client *redis.Client, logger *zap.Logger) domain.CartRepository {
	return &cartRepository{
		client: client,
		logger: logger,
	}
}

// Redis key format
func (r *cartRepository) getCartKey(userID string) string {
	return fmt.Sprintf("cart:user:%s", userID)
}

// GetCart retrieves a cart from Redis
func (r *cartRepository) GetCart(userID string) (*domain.ShoppingCart, error) {
	ctx := context.Background()
	key := r.getCartKey(userID)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Return empty cart
		return &domain.ShoppingCart{
			UserID:  userID,
			Items:   make([]*domain.CartItem, 0),
			Version: 1,
		}, nil
	}
	if err != nil {
		r.logger.Error("failed to get cart from Redis",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, fmt.Errorf("failed to get cart from Redis: %w", err)
	}

	var cart domain.ShoppingCart
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		r.logger.Error("failed to unmarshal cart",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	// Ensure UserID is set
	if cart.UserID == "" {
		cart.UserID = userID
	}

	return &cart, nil
}

// SaveCart saves a cart to Redis with TTL
func (r *cartRepository) SaveCart(cart *domain.ShoppingCart) error {
	if cart.UserID == "" {
		return fmt.Errorf("user_id is required - authentication required")
	}

	ctx := context.Background()
	key := r.getCartKey(cart.UserID)

	// Update metadata
	cart.Version++ // Increment version for optimistic locking

	// Create minimal cart for Redis storage (without computed fields)
	minimalCart := struct {
		UserID  string             `json:"user_id"`
		Items   []*domain.CartItem `json:"items"`
		Version int                `json:"version"`
	}{
		UserID:  cart.UserID,
		Items:   cart.Items,
		Version: cart.Version,
	}

	// Serialize to JSON
	cartJSON, err := json.Marshal(minimalCart)
	if err != nil {
		r.logger.Error("failed to marshal cart",
			zap.Error(err),
			zap.String("user_id", cart.UserID),
		)
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	// Save with 30 days TTL
	ttl := 30 * 24 * time.Hour
	if err := r.client.Set(ctx, key, cartJSON, ttl).Err(); err != nil {
		r.logger.Error("failed to save cart to Redis",
			zap.Error(err),
			zap.String("user_id", cart.UserID),
		)
		return fmt.Errorf("failed to save cart to Redis: %w", err)
	}

	r.logger.Info("cart saved successfully",
		zap.String("user_id", cart.UserID),
		zap.Int("item_count", len(cart.Items)),
		zap.Int("version", cart.Version),
	)

	return nil
}

// DeleteCart removes a cart from Redis
func (r *cartRepository) DeleteCart(userID string) error {
	ctx := context.Background()
	key := r.getCartKey(userID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		r.logger.Error("failed to delete cart",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	r.logger.Info("cart deleted successfully",
		zap.String("user_id", userID),
	)

	return nil
}

// ClearSelectedItems removes only selected items from cart
// This is called after successful checkout
func (r *cartRepository) ClearSelectedItems(userID string) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	// Filter out selected items
	unselectedItems := make([]*domain.CartItem, 0)
	for _, item := range cart.Items {
		if !item.IsSelected {
			unselectedItems = append(unselectedItems, item)
		}
	}

	// Update cart with only unselected items
	cart.Items = unselectedItems

	r.logger.Info("cleared selected items from cart",
		zap.String("user_id", userID),
		zap.Int("remaining_items", len(unselectedItems)),
	)

	return r.SaveCart(cart)
}

// AddItem adds a new item to cart or updates quantity if exists
func (r *cartRepository) AddItem(userID string, item *domain.CartItem) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	// Check if item already exists
	found := false
	for _, existingItem := range cart.Items {
		if existingItem.ProductItemID == item.ProductItemID {
			// Update quantity
			existingItem.Quantity += item.Quantity
			found = true
			break
		}
	}

	// Add new item if not found
	if !found {
		cart.Items = append(cart.Items, item)
	}

	r.logger.Info("item added to cart",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", item.ProductItemID),
		zap.Int("quantity", item.Quantity),
	)

	return r.SaveCart(cart)
}

// UpdateItemQuantity updates the quantity of a specific item
func (r *cartRepository) UpdateItemQuantity(userID string, productItemID uint, quantity int) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	found := false
	for _, item := range cart.Items {
		if item.ProductItemID == productItemID {
			item.Quantity = quantity
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found in cart")
	}

	r.logger.Info("cart item quantity updated",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", productItemID),
		zap.Int("new_quantity", quantity),
	)

	return r.SaveCart(cart)
}

// RemoveItem removes a specific item from cart
func (r *cartRepository) RemoveItem(userID string, productItemID uint) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	// Filter out the item
	newItems := make([]*domain.CartItem, 0)
	removed := false
	for _, item := range cart.Items {
		if item.ProductItemID != productItemID {
			newItems = append(newItems, item)
		} else {
			removed = true
		}
	}

	if !removed {
		return fmt.Errorf("item not found in cart")
	}

	cart.Items = newItems

	r.logger.Info("item removed from cart",
		zap.String("user_id", userID),
		zap.Uint("product_item_id", productItemID),
	)

	return r.SaveCart(cart)
}

// ToggleItemSelection toggles the selection state of an item
func (r *cartRepository) ToggleItemSelection(userID string, productItemID uint) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	found := false
	for _, item := range cart.Items {
		if item.ProductItemID == productItemID {
			item.IsSelected = !item.IsSelected
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found in cart")
	}

	return r.SaveCart(cart)
}

// SelectAllItems selects or deselects all items in cart
func (r *cartRepository) SelectAllItems(userID string, selected bool) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	for _, item := range cart.Items {
		item.IsSelected = selected
	}

	r.logger.Info("all items selection updated",
		zap.String("user_id", userID),
		zap.Bool("selected", selected),
	)

	return r.SaveCart(cart)
}

// GetSelectedItems returns only selected items
func (r *cartRepository) GetSelectedItems(userID string) ([]*domain.CartItem, error) {
	cart, err := r.GetCart(userID)
	if err != nil {
		return nil, err
	}

	selectedItems := make([]*domain.CartItem, 0)
	for _, item := range cart.Items {
		if item.IsSelected {
			selectedItems = append(selectedItems, item)
		}
	}

	return selectedItems, nil
}

// GetCartItemCount returns total number of items in cart
func (r *cartRepository) GetCartItemCount(userID string) (int, error) {
	cart, err := r.GetCart(userID)
	if err != nil {
		return 0, err
	}

	totalCount := 0
	for _, item := range cart.Items {
		totalCount += item.Quantity
	}

	return totalCount, nil
}
