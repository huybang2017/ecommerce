package middleware

import (
	"fmt"
	"identity-service/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from Authorization header or cookie and sets user context
// Priority:
//  1. Authorization header (Bearer token) - forwarded by API Gateway
//  2. Role-prefixed cookie ({role}_access_token) - for direct service access
//  3. Legacy cookie (access_token) - backward compatibility
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var tokenSource string

		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] Request: %s %s\n", c.Request.Method, c.Request.URL.Path)

		// PRIORITY 1: Authorization header (forwarded by Gateway)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			if len(authHeader) > 7 && (authHeader[:7] == "Bearer " || authHeader[:7] == "bearer ") {
				token = authHeader[7:]
				tokenSource = "Authorization header"
				fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ✅ Token from Authorization header\n")
			}
		}

		// PRIORITY 2: Role-prefixed cookie (admin_access_token, buyer_access_token, seller_access_token)
		if token == "" {
			// Get role from Gateway's X-Resolved-Role header
			role := c.GetHeader("X-Resolved-Role")
			if role != "" {
				cookieName := role + "_access_token"
				if cookieToken, err := c.Cookie(cookieName); err == nil && cookieToken != "" {
					token = cookieToken
					tokenSource = cookieName + " cookie"
					fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ✅ Token from %s cookie\n", cookieName)
				}
			}
		}

		// PRIORITY 3: Legacy cookie (backward compatibility)
		if token == "" {
			if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
				token = cookieToken
				tokenSource = "access_token cookie"
				fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ✅ Token from access_token cookie (legacy)\n")
			}
		}

		// No token found
		if token == "" {
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ❌ No auth token found (checked: Authorization header, role cookies, legacy cookie)\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Validate token
		userID, role, err := authService.ValidateToken(token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ❌ Invalid token from %s: %v\n", tokenSource, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ✅ Token validated (source: %s, user: %v, role: %s)\n", tokenSource, userID, role)

		// Set user context
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}
