package middleware

import (
	"fmt"
	"identity-service/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token from cookie and sets user context
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access_token from HttpOnly cookie
		token, err := c.Cookie("access_token")
		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] access_token cookie present: %v\n", err == nil && token != "")

		if err != nil || token == "" {
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ❌ Missing access_token cookie!\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Validate token
		userID, role, err := authService.ValidateToken(token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ❌ Invalid token: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ✅ Token validated for user: %v\n", userID)

		// Set user context
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}
