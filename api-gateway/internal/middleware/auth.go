package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"api-gateway/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AuthMiddleware validates JWT tokens for protected routes
// This implements authentication for the API Gateway
func AuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Warn("Invalid signing method", zap.String("method", fmt.Sprintf("%v", token.Method)))
				return nil, jwt.ErrSignatureInvalid
			}
			logger.Debug("Validating token with secret", zap.String("secret_length", fmt.Sprintf("%d", len(cfg.Secret))))
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			logger.Warn("Token validation failed", zap.Error(err), zap.String("token_preview", tokenString[:min(20, len(tokenString))]+"..."))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Warn("Token is not valid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims and store in context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Convert user_id to string for consistency
			if userIDFloat, ok := claims["user_id"].(float64); ok {
				userID := fmt.Sprintf("%.0f", userIDFloat)
				c.Set("user_id", userID)
				logger.Debug("User authenticated", zap.String("user_id", userID), zap.String("email", claims["email"].(string)))
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}

		// CRITICAL: Preserve Authorization header in context for forwarding to backend services
		// This ensures the header is available even if something modifies c.Request.Header
		// IMPORTANT: Use the original authHeader variable, not c.Request.Header.Get again
		// because c.Request.Header might have been modified
		c.Set("auth_header", authHeader)
		logger.Debug("Preserved Authorization header in context", zap.String("header_preview", authHeader[:min(30, len(authHeader))]))

		c.Next()
	}
}

// OptionalAuthMiddleware allows requests with or without authentication
// Useful for routes that have optional authentication
func OptionalAuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID := fmt.Sprintf("%.0f", claims["user_id"].(float64))
				c.Set("user_id", userID)
				c.Set("email", claims["email"])
				c.Set("role", claims["role"])
			}
		}

		c.Next()
	}
}

