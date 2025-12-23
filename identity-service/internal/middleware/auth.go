package middleware

import (
	"fmt"
	"identity-service/internal/service"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] Authorization header present: %v\n", authHeader != "")
		if authHeader != "" {
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] Auth header: %s...\n", authHeader[:min(30, len(authHeader))])
		} else {
			// Debug: Log all headers to see what we received
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] ❌ Missing Authorization header!\n")
			fmt.Fprintf(os.Stderr, "[DEBUG Identity Service] All headers: %v\n", c.Request.Header)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Extract token (format: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		userID, role, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}

