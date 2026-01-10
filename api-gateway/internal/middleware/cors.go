package middleware

import (
	"api-gateway/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CORSMiddleware creates a custom CORS middleware with proper credentials support
func CORSMiddleware(cfg *config.CORSConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Set CORS headers for all requests with Origin header
		if origin != "" {
			allowedOrigin := getMatchedOrigin(origin, cfg.AllowedOrigins)

			if allowedOrigin == "" {
				// Origin not in allowed list - use it anyway for development
				allowedOrigin = origin
				logger.Warn("Origin not in allowed list",
					zap.String("origin", origin),
					zap.Strings("allowed_origins", cfg.AllowedOrigins),
				)
			}

			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", allowedOrigin)

			// Debug: log config
			logger.Debug("CORS config",
				zap.Strings("allowed_methods", cfg.AllowedMethods),
				zap.Int("allowed_methods_count", len(cfg.AllowedMethods)),
			)

			// Hardcode methods if config is empty (fallback for safety)
			allowedMethods := strings.Join(cfg.AllowedMethods, ", ")
			if allowedMethods == "" {
				allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
				logger.Warn("AllowedMethods config is empty, using default")
			}
			h.Set("Access-Control-Allow-Methods", allowedMethods)

			h.Set("Access-Control-Allow-Credentials", "true")
			h.Set("Access-Control-Expose-Headers", strings.Join(cfg.ExposeHeaders, ", "))

			// Use requested headers if provided, otherwise use config
			if reqHeaders := c.Request.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
				h.Set("Access-Control-Allow-Headers", reqHeaders)
			} else {
				h.Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			}

			if cfg.MaxAge > 0 {
				h.Set("Access-Control-Max-Age", strconv.Itoa(int(cfg.MaxAge.Seconds())))
			}
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// getMatchedOrigin checks if origin is in allowed list
func getMatchedOrigin(origin string, allowedOrigins []string) string {
	for _, allowed := range allowedOrigins {
		if allowed == origin {
			return origin
		}
	}
	return ""
}
