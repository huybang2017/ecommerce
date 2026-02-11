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

const (
	RoleAdmin  = "ADMIN"
	RoleSeller = "SELLER"
	RoleBuyer  = "BUYER"
)

// resolvedPrefix reads the role_prefix set by RoleCookieRouter middleware.
// Falls back to deriving it from X-App-Role if the router hasn't run (backward compat).
func resolvedPrefix(c *gin.Context) string {
	if p, ok := c.Get("role_prefix"); ok {
		if prefix, ok := p.(string); ok {
			return prefix
		}
	}
	// fallback – should not happen if pipeline is wired correctly
	role := strings.ToLower(c.GetHeader("X-App-Role"))
	if role != "" {
		return role + "_"
	}
	return ""
}

// AuthMiddleware validates JWT tokens for protected routes.
// Expects RoleCookieRouter to have already run — cookies are pre-filtered
// to only contain the role-matching set.
//
// Token source priority:
//   1. Role-scoped HttpOnly cookie ({role}_access_token)
//   2. Authorization header (bearer token)
func AuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		prefix := resolvedPrefix(c)

		// PRIORITY 1: Role-scoped HttpOnly cookie (pre-filtered by RoleCookieRouter)
		if prefix != "" {
			if cookieToken, err := c.Cookie(prefix + "access_token"); err == nil && cookieToken != "" {
				tokenString = cookieToken
				log.Printf("[AUTH] Token from %saccess_token cookie", prefix)
			}
		}

		// PRIORITY 2: Authorization header (API clients, Swagger, etc.)
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				logger.Warn("Missing authorization credentials (no cookie or header)")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization credentials"})
				c.Abort()
				return
			}

			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				tokenString = authHeader[7:]
			} else {
				tokenString = strings.TrimSpace(authHeader)
			}
			log.Printf("[AUTH] Token from Authorization header")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization credentials"})
			c.Abort()
			return
		}

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})
		if err != nil {
			preview := tokenString
			if len(preview) > 20 {
				preview = preview[:20] + "..."
			}
			log.Printf("[AUTH] Token validation failed: %v, preview=%s", err, preview)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims → context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userIDFloat, ok := claims["user_id"].(float64); ok {
				userID := fmt.Sprintf("%.0f", userIDFloat)
				c.Set("user_id", userID)
				c.Set("user_id_uint", uint(userIDFloat))
				log.Printf("[AUTH] Authenticated user_id=%s", userID)
			}
			if email, ok := claims["email"].(string); ok {
				c.Set("email", email)
			}
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role)

				// ── Role-JWT consistency check ──────────────────────────
				// If RoleCookieRouter resolved a role, verify the JWT role
				// is consistent with the declared app role.
				if resolvedRole, exists := c.Get("resolved_role"); exists {
					headerRole := strings.ToUpper(resolvedRole.(string))
					if role != headerRole {
						log.Printf("[AUTH] SECURITY: role mismatch jwt_role=%s header_role=%s", role, headerRole)
						logger.Warn("Role mismatch between JWT and X-App-Role",
							zap.String("jwt_role", role),
							zap.String("header_role", headerRole),
						)
						c.JSON(http.StatusForbidden, gin.H{"error": "Role mismatch: token role does not match app role"})
						c.Abort()
						return
					}
				}
			}
		}

		// Store bearer token for forwarding to backend services
		c.Set("auth_header", "Bearer "+tokenString)
		log.Printf("[AUTH] Authentication successful")

		c.Next()
	}
}

// SessionMiddleware validates session_id cookie against Redis.
// Expects RoleCookieRouter to have already filtered cookies.
func SessionMiddleware(logger *zap.Logger, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			log.Printf("[SESSION] user_id not in context – AuthMiddleware must run first")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user_id in context"})
			c.Abort()
			return
		}
		userID, ok := userIDVal.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id type in context"})
			c.Abort()
			return
		}

		// Read role-scoped session_id (cookies already filtered by RoleCookieRouter)
		prefix := resolvedPrefix(c)
		var sessionID string
		var err error

		if prefix != "" {
			sessionID, err = c.Cookie(prefix + "session_id")
		}

		if err != nil || sessionID == "" {
			log.Printf("[SESSION] Missing %ssession_id cookie user_id=%s", prefix, userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session_id cookie"})
			c.Abort()
			return
		}

		// Validate against Redis
		key := fmt.Sprintf("session:%s", sessionID)
		sessionJSON, err := redisClient.Get(c.Request.Context(), key).Result()
		if err != nil {
			log.Printf("[SESSION] Not found key=%s user_id=%s err=%v", key, userID, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		var session SessionData
		if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
			log.Printf("[SESSION] Unmarshal failed key=%s err=%v", key, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session data"})
			c.Abort()
			return
		}

		sessionUserID := fmt.Sprintf("%d", session.UserID)
		if sessionUserID != userID {
			log.Printf("[SESSION] User mismatch token_uid=%s session_uid=%s", userID, sessionUserID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session user mismatch"})
			c.Abort()
			return
		}

		c.Set("session_id", sessionID)
		log.Printf("[SESSION] Validated user_id=%s session_id=%s", userID, sessionID)

		c.Next()
	}
}

// RoleMiddleware checks that the authenticated user has one of the allowed roles.
func RoleMiddleware(logger *zap.Logger, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			logger.Warn("Role not found in context")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: role not identified"})
			c.Abort()
			return
		}

		userRole, _ := roleVal.(string)

		isAllowed := false
		for _, r := range allowedRoles {
			if userRole == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			logger.Warn("Unauthorized role access attempt",
				zap.String("user_role", userRole),
				zap.Strings("allowed_roles", allowedRoles),
			)
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware allows requests with or without authentication.
func OptionalAuthMiddleware(cfg *config.JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try role-scoped cookie first
		prefix := resolvedPrefix(c)
		var tokenString string

		if prefix != "" {
			if cookieToken, err := c.Cookie(prefix + "access_token"); err == nil && cookieToken != "" {
				tokenString = cookieToken
			}
		}

		// Fallback to Authorization header
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.Next()
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				c.Next()
				return
			}
			tokenString = parts[1]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.Secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if uid, ok := claims["user_id"].(float64); ok {
					c.Set("user_id", fmt.Sprintf("%.0f", uid))
				}
				if email, ok := claims["email"].(string); ok {
					c.Set("email", email)
				}
				if role, ok := claims["role"].(string); ok {
					c.Set("role", role)
				}
			}
		}

		c.Next()
	}
}
