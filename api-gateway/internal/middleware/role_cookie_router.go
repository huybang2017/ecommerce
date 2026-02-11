package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ──────────────────────────────────────────────────────────────────────────────
// Role-Aware Cookie Router Middleware
//
// Responsibilities:
//   1. Read & normalize X-App-Role header
//   2. Validate role (buyer | admin | seller)
//   3. Map role → cookie name prefix
//   4. Filter cookies: keep ONLY the role-matching cookies, drop all others
//   5. Inject resolved role into request context + headers for downstream
//   6. Reject requests with missing/invalid role or missing role cookie
//
// Cookie naming convention (set by identity-service):
//   {role}_access_token   e.g. admin_access_token
//   {role}_session_id     e.g. admin_session_id
//   {role}_refresh_token  e.g. admin_refresh_token
//
// This middleware MUST run before AuthMiddleware and SessionMiddleware.
// ──────────────────────────────────────────────────────────────────────────────

// validRoles is the set of accepted X-App-Role values.
var validRoles = map[string]bool{
	"buyer":  true,
	"admin":  true,
	"seller": true,
}

// roleCookieSuffixes lists the cookie name suffixes owned by each role.
var roleCookieSuffixes = []string{"access_token", "session_id", "refresh_token"}

// RoleCookieRouter resolves X-App-Role, filters cookies by role, and injects
// role context for downstream middleware (AuthMiddleware, SessionMiddleware).
//
// After this middleware:
//   - c.Get("resolved_role")   → "buyer" | "admin" | "seller"
//   - c.Get("role_prefix")     → "buyer_" | "admin_" | "seller_"
//   - X-Resolved-Role header   → set on request for backend forwarding
//   - Cookie header            → contains ONLY the role-matching cookies
//
// For auth endpoints (login, register, refresh) that have no cookie yet,
// use RoleCookieRouterLax which validates the role header but does NOT
// require the cookie to be present.
func RoleCookieRouter(logger *zap.Logger) gin.HandlerFunc {
	return roleCookieRouterInternal(logger, true)
}

// RoleCookieRouterLax is the same as RoleCookieRouter but does NOT reject
// requests missing a role cookie. Used for auth endpoints (login, register,
// refresh) where the cookie does not exist yet.
func RoleCookieRouterLax(logger *zap.Logger) gin.HandlerFunc {
	return roleCookieRouterInternal(logger, false)
}

func roleCookieRouterInternal(logger *zap.Logger, requireCookie bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── 1. Read & normalize X-App-Role ──────────────────────────────
		rawRole := c.GetHeader("X-App-Role")
		role := strings.ToLower(strings.TrimSpace(rawRole))

		if role == "" {
			logger.Warn("[ROLE-ROUTER] Missing X-App-Role header",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing X-App-Role header"})
			c.Abort()
			return
		}

		// ── 2. Validate role ────────────────────────────────────────────
		if !validRoles[role] {
			logger.Warn("[ROLE-ROUTER] Invalid X-App-Role",
				zap.String("role", role),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Invalid role: %s. Allowed: buyer, admin, seller", role),
			})
			c.Abort()
			return
		}

		// ── 3. Map role → cookie prefix ─────────────────────────────────
		prefix := role + "_" // e.g. "admin_"

		// ── 4. Filter cookies ───────────────────────────────────────────
		// Build a set of cookie names that belong to THIS role.
		ownedNames := make(map[string]bool, len(roleCookieSuffixes))
		for _, suffix := range roleCookieSuffixes {
			ownedNames[prefix+suffix] = true
		}

		// Parse all cookies, keep only role-owned ones.
		allCookies := c.Request.Cookies()
		var kept []*http.Cookie
		var dropped []string

		for _, ck := range allCookies {
			if ownedNames[ck.Name] {
				kept = append(kept, ck)
			} else {
				// Check if this cookie belongs to ANOTHER role (cross-role)
				for otherRole := range validRoles {
					if otherRole == role {
						continue
					}
					otherPrefix := otherRole + "_"
					for _, suffix := range roleCookieSuffixes {
						if ck.Name == otherPrefix+suffix {
							dropped = append(dropped, ck.Name)
						}
					}
				}
				// Non-auth cookies (e.g. _ga, theme) are kept as-is
				isForeignAuth := false
				for _, d := range dropped {
					if d == ck.Name {
						isForeignAuth = true
						break
					}
				}
				if !isForeignAuth {
					kept = append(kept, ck)
				}
			}
		}

		if len(dropped) > 0 {
			logger.Debug("[ROLE-ROUTER] Dropped cross-role cookies",
				zap.String("role", role),
				zap.Strings("dropped", dropped),
			)
		}

		// ── 5. Check role cookie exists ─────────────────────────────────
		if requireCookie {
			hasRoleCookie := false
			for _, ck := range kept {
				if ck.Name == prefix+"access_token" || ck.Name == prefix+"session_id" {
					hasRoleCookie = true
					break
				}
			}
			if !hasRoleCookie {
				logger.Warn("[ROLE-ROUTER] No cookie found for role",
					zap.String("role", role),
					zap.String("expected_prefix", prefix),
					zap.String("path", c.Request.URL.Path),
				)
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": fmt.Sprintf("No auth cookie found for role: %s", role),
				})
				c.Abort()
				return
			}
		}

		// ── 6. Rewrite Cookie header with filtered cookies ──────────────
		// Rebuild the Cookie header so downstream code (including proxy)
		// only sees role-matching + non-auth cookies.
		c.Request.Header.Del("Cookie")
		var parts []string
		for _, ck := range kept {
			parts = append(parts, ck.Name+"="+ck.Value)
		}
		if len(parts) > 0 {
			c.Request.Header.Set("Cookie", strings.Join(parts, "; "))
		}

		// ── 7. Inject role context ──────────────────────────────────────
		c.Set("resolved_role", role)
		c.Set("role_prefix", prefix)

		// Set X-Resolved-Role on request for backend forwarding
		c.Request.Header.Set("X-Resolved-Role", role)

		// Ensure X-App-Role is normalized (lowercase)
		c.Request.Header.Set("X-App-Role", role)

		logger.Debug("[ROLE-ROUTER] Role resolved, cookies filtered",
			zap.String("role", role),
			zap.String("prefix", prefix),
			zap.Int("kept_cookies", len(kept)),
			zap.Int("dropped_cookies", len(dropped)),
		)

		c.Next()
	}
}
