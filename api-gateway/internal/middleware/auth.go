package middleware

import (
	"api-gateway/config"
	"fmt"
	"net/http"
	"strings"

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
// Supports both Cookie-based (preferred) and Authorization header authentication
func AuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// PRIORITY 1: Try to get token from HttpOnly cookie (most secure)
		if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
			tokenString = cookieToken
			logger.Debug("Token found in cookie")
		} else {
			// PRIORITY 2: Fallback to Authorization header (for compatibility)
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				logger.Warn("Missing authorization credentials (no cookie or header)")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization credentials"})
				c.Abort()
				return
			}

			// Normalize Authorization header: auto-add "Bearer " prefix if missing
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else if strings.HasPrefix(authHeader, "bearer ") {
				tokenString = strings.TrimPrefix(strings.ToLower(authHeader), "bearer ")
			} else {
				tokenString = strings.TrimSpace(authHeader)
			}
			logger.Debug("Token found in Authorization header")
		}

		// Validate token is not empty
		if tokenString == "" {
			logger.Warn("Empty token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization credentials"})
			c.Abort()
			return
		}

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

				// Also set as uint for backend services compatibility
				c.Set("user_id_uint", uint(userIDFloat))

				logger.Debug("User authenticated", zap.String("user_id", userID))
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)
			}
		}

		// Store token for forwarding to backend services
		// Create Bearer token format for header forwarding
		bearerToken := "Bearer " + tokenString
		c.Set("auth_header", bearerToken)
		logger.Debug("Authentication successful")

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
