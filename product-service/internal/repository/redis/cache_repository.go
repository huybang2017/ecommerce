package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"product-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

// cacheRepository handles Redis operations for product caching
// This is the infrastructure layer - it knows HOW to interact with Redis
type cacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new Redis cache repository
// Dependency injection: we inject the Redis client
func NewCacheRepository(client *redis.Client) *cacheRepository {
	return &cacheRepository{client: client}
}

// SetProduct caches a product in Redis with a TTL
// TTL prevents stale data and manages memory usage
func (r *cacheRepository) SetProduct(ctx context.Context, product *domain.Product, ttl time.Duration) error {
	key := fmt.Sprintf("product:%d", product.ID)

	// Serialize product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	// Set with expiration
	err = r.client.Set(ctx, key, productJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set product in cache: %w", err)
	}

	return nil
}

// GetProduct retrieves a product from Redis cache
// Returns nil if not found (cache miss)
func (r *cacheRepository) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	key := fmt.Sprintf("product:%d", id)

	// Get from Redis
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss - not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product from cache: %w", err)
	}

	// Deserialize JSON to Product
	var product domain.Product
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return &product, nil
}

// DeleteProduct removes a product from Redis cache
func (r *cacheRepository) DeleteProduct(ctx context.Context, id uint) error {
	key := fmt.Sprintf("product:%d", id)
	return r.client.Del(ctx, key).Err()
}

// AcquireLock acquires a distributed lock using Redis
// This is useful for preventing race conditions (e.g., inventory updates)
// Returns true if lock was acquired, false if already locked
func (r *cacheRepository) AcquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, error) {
	// Use SET with NX (only if not exists) and EX (expiration)
	result, err := r.client.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return result, nil
}

// ReleaseLock releases a distributed lock
func (r *cacheRepository) ReleaseLock(ctx context.Context, lockKey string) error {
	return r.client.Del(ctx, lockKey).Err()
}

// Get retrieves a raw value from Redis (generic helper)
func (r *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set sets a raw value in Redis with TTL (generic helper)
func (r *cacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

