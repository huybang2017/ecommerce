package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SkipOptionsLoggingMiddleware skips logging for CORS preflight OPTIONS requests
// This reduces log noise while maintaining security
// NOTE: This middleware only skips LOGGING, not processing. CORS middleware still handles OPTIONS.
func SkipOptionsLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For OPTIONS requests, just pass through without logging
		// CORS middleware will handle it
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Log all other requests
		logger.Debug("Request received",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()
	}
}
