package middleware

import (
	"net/http"
	"sync"
	"time"
	"api-gateway/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"go.uber.org/zap"
)

// rateLimiter stores rate limiters per IP address
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	config   *config.RateLimitConfig
}

// newRateLimiter creates a new rate limiter
func newRateLimiter(cfg *config.RateLimitConfig) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   cfg,
	}
}

// getLimiter returns a rate limiter for the given IP
func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		// Create a new limiter: requests per minute converted to requests per second
		limiter = rate.NewLimiter(
			rate.Limit(rl.config.RequestsPerMinute)/60,
			rl.config.Burst,
		)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// cleanup removes old limiters periodically
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// In production, you'd want more sophisticated cleanup logic
			// For now, we keep all limiters in memory
			rl.mu.Unlock()
		}
	}()
}

var globalRateLimiter *rateLimiter

// RateLimitMiddleware implements rate limiting per IP address
// This prevents abuse and ensures fair resource usage
func RateLimitMiddleware(cfg *config.RateLimitConfig, logger *zap.Logger) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	if globalRateLimiter == nil {
		globalRateLimiter = newRateLimiter(cfg)
		globalRateLimiter.cleanup()
	}

	return func(c *gin.Context) {
		// Get client IP
		ip := c.ClientIP()

		// Get or create limiter for this IP
		limiter := globalRateLimiter.getLimiter(ip)

		// Check if request is allowed
		if !limiter.Allow() {
			logger.Warn("Rate limit exceeded", zap.String("ip", ip))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

