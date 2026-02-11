package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StripInternalHeaders removes any incoming headers that are intended
// to be injected only by the API Gateway. This prevents clients from
// spoofing internal identity headers like X-User-Id.
func StripInternalHeaders(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Collect keys to delete to avoid mutating map during iteration
		var toDelete []string
		for key := range c.Request.Header {
			lower := strings.ToLower(key)
			if strings.HasPrefix(lower, "x-user-") || strings.HasPrefix(lower, "x-internal-") {
				toDelete = append(toDelete, key)
			}
			// X-Resolved-Role is injected by RoleCookieRouter — never trust from client
			if lower == "x-resolved-role" {
				toDelete = append(toDelete, key)
			}
		}
		for _, k := range toDelete {
			c.Request.Header.Del(k)
		}

		// Ensure hop-by-hop headers are removed (case-insensitive)
		for _, h := range []string{"connection", "keep-alive", "transfer-encoding", "upgrade"} {
			c.Request.Header.Del(h)
			c.Request.Header.Del(strings.Title(h))
		}

		c.Next()
	}
}
