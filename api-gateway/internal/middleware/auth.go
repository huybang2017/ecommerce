package middleware

import (
	"api-gateway/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// SessionData represents the session object stored in Redis
// We only care about user_id here
type SessionData struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"`
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
			log.Printf("[AUTH] Token found in cookie")
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
			log.Printf("[AUTH] Token found in Authorization header")
		}

		// Validate token is not empty
		if tokenString == "" {
			log.Printf("[AUTH] Empty token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization credentials"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("[AUTH] Invalid signing method: %v", token.Method)
				return nil, jwt.ErrSignatureInvalid
			}
			log.Printf("[AUTH] Validating token with secret (len=%d)", len(cfg.Secret))
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			preview := tokenString
			if len(preview) > 20 {
				preview = preview[:20] + "..."
			}
			log.Printf("[AUTH] Token validation failed: %v, token_preview=%s", err, preview)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Printf("[AUTH] Token is not valid")
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
				log.Printf("[AUTH] User authenticated user_id=%s", userID)
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
		log.Printf("[AUTH] Authentication successful")

		c.Next()
	}
}

// SessionMiddleware validates session_id từ cookie với Redis
// Kiểm tra user có sở hữu session này không
func SessionMiddleware(logger *zap.Logger, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy user_id từ context (set bởi AuthMiddleware)
		userIDVal, exists := c.Get("user_id")
		if !exists {
			log.Printf("[SESSION] user_id not found in context - auth middleware should run first")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id in context"})
			c.Abort()
			return
		}
		userID, ok := userIDVal.(string)
		if !ok {
			log.Printf("[SESSION] user_id in context is not string, type=%T, value=%v", userIDVal, userIDVal)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id type in context"})
			c.Abort()
			return
		}

		// Lấy session_id từ cookie
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			log.Printf("[SESSION] Missing session_id cookie user_id=%s err=%v", userID, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session_id cookie"})
			c.Abort()
			return
		}

		// Lấy session JSON từ Redis
		key := fmt.Sprintf("session:%s", sessionID)
		sessionJSON, err := redisClient.Get(c.Request.Context(), key).Result()
		if err != nil {
			log.Printf("[SESSION] Session not found or expired key=%s user_id=%s err=%v", key, userID, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}
		log.Printf("[SESSION] Loaded session from redis key=%s value=%s", key, sessionJSON)

		// Parse JSON -> SessionData
		var session SessionData
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			log.Printf("[SESSION] Failed to unmarshal session JSON key=%s err=%v", key, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}
		sessionUserID := fmt.Sprintf("%d", session.UserID)

		// Kiểm tra user_id từ token có match với session không
		if sessionUserID != userID {
			log.Printf("\n[SESSION ERROR] ===================================\n"+
				"| Reason: User Mismatch\n"+
				"| Token UID:   %s\n"+
				"| Session UID: %s\n"+
				"| Raw Session: %s\n"+
				"===================================================",
				userID, sessionUserID, sessionJSON)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session user mismatch"})
			c.Abort()
			return
		}

		// Lưu session_id vào context
		c.Set("session_id", sessionID)
		log.Printf("[SESSION] Session validated successfully user_id=%s session_id=%s", userID, sessionID)

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
