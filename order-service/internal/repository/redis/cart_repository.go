package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

// cartRepository handles Redis operations for cart storage
// This is the infrastructure layer - it knows HOW to interact with Redis
type cartRepository struct {
	client *redis.Client
}

// NewCartRepository creates a new Redis cart repository
// Dependency injection: we inject the Redis client
func NewCartRepository(client *redis.Client) *cartRepository {
	return &cartRepository{client: client}
}

// getCartKey generates the Redis key for a cart
// Format: "cart:user:{user_id}" - only authenticated users
// Business rule: Cart requires authentication - session_id is no longer supported
func (r *cartRepository) getCartKey(userID string) string {
	return fmt.Sprintf("cart:user:%s", userID)
}

// GetCart retrieves a cart from Redis
// Business rule: Only authenticated users - userID is required
func (r *cartRepository) GetCart(userID string) (*domain.Cart, error) {
	ctx := context.Background()
	key := r.getCartKey(userID)

	// Get from Redis
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cart doesn't exist, return empty cart
		return &domain.Cart{
			UserID:    userID,
			Items:     make(map[uint]*domain.CartItem),
			Total:     0,
			UpdatedAt: time.Now().Unix(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cart from Redis: %w", err)
	}

	// Deserialize JSON to Cart
	var cart domain.Cart
	err = json.Unmarshal([]byte(val), &cart)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	// Ensure UserID is set (for backward compatibility)
	if cart.UserID == "" {
		cart.UserID = userID
	}

	return &cart, nil
}

// SaveCart saves a cart to Redis
// Cart expires after 30 days of inactivity
// Business rule: Only authenticated users - UserID is required
func (r *cartRepository) SaveCart(cart *domain.Cart) error {
	if cart.UserID == "" {
		return fmt.Errorf("user_id is required - authentication required")
	}

	ctx := context.Background()
	key := r.getCartKey(cart.UserID)

	// Update timestamp
	cart.UpdatedAt = time.Now().Unix()

	// Calculate total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	// Serialize cart to JSON
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart: %w", err)
	}

	// Set with expiration (30 days)
	ttl := 30 * 24 * time.Hour
	err = r.client.Set(ctx, key, cartJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save cart to Redis: %w", err)
	}

	return nil
}

// DeleteCart removes a cart from Redis
// Business rule: Only authenticated users - userID is required
func (r *cartRepository) DeleteCart(userID string) error {
	ctx := context.Background()
	key := r.getCartKey(userID)
	return r.client.Del(ctx, key).Err()
}

// ClearCartItems clears all items from a cart but keeps the cart structure
// Business rule: Only authenticated users - userID is required
func (r *cartRepository) ClearCartItems(userID string) error {
	cart, err := r.GetCart(userID)
	if err != nil {
		return err
	}

	cart.Items = make(map[uint]*domain.CartItem)
	cart.Total = 0
	return r.SaveCart(cart)
}


