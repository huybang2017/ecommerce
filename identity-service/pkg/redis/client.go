package redis

import (
	"context"
	"fmt"
	"identity-service/config"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// clientInstance is the singleton Redis client
	clientInstance *redis.Client
	// once ensures the client is created only once
	once sync.Once
)

// GetClient returns the singleton Redis client
// This implements the Singleton pattern to ensure only one Redis connection pool exists
func GetClient(cfg *config.RedisConfig) (*redis.Client, error) {
	var err error

	once.Do(func() {
		clientInstance = redis.NewClient(&redis.Options{
			Addr:         cfg.GetAddress(),
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
		})

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = clientInstance.Ping(ctx).Err(); err != nil {
			log.Printf("Failed to connect to Redis: %v", err)
			return
		}

		log.Println("Redis connection established successfully")
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	return clientInstance, nil
}

// CloseClient closes the Redis client connection
// This should be called during graceful shutdown
func CloseClient() error {
	if clientInstance == nil {
		return nil
	}

	return clientInstance.Close()
}
