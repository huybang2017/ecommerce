package handler

import "github.com/gin-gonic/gin"

// cookiePrefix returns the role-specific prefix for cookie names.
// Each frontend sends "X-App-Role" header (admin | buyer | seller).
// This enables multiple concurrent sessions in the same browser.
//
// Cookie names:
//   admin_session_id,  admin_access_token,  admin_refresh_token
//   buyer_session_id,  buyer_access_token,  buyer_refresh_token
//   seller_session_id, seller_access_token, seller_refresh_token
func cookiePrefix(c *gin.Context) string {
	role := c.GetHeader("X-App-Role")
	switch role {
	case "admin":
		return "admin_"
	case "seller":
		return "seller_"
	case "buyer":
		return "buyer_"
	default:
		return "" // backward compatible: no prefix
	}
}

// setAuthCookies sets session_id, access_token, refresh_token with role prefix.
// Also sets is_logged_in flag (non-httpOnly) for Next.js middleware.
func setAuthCookies(c *gin.Context, sessionID, accessToken, refreshToken string) {
	prefix := cookiePrefix(c)

	// ────────────────────────────────────────────────────────────────────────
	// Set httpOnly cookies (secure, not readable by JavaScript)
	// ────────────────────────────────────────────────────────────────────────
	c.SetCookie(prefix+"session_id", sessionID, 604800, "/", "", false, true)
	c.SetCookie(prefix+"access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie(prefix+"refresh_token", refreshToken, 604800, "/", "", false, true)

	// ────────────────────────────────────────────────────────────────────────
	// ✅ NEW: Set is_logged_in flag (non-httpOnly) for Next.js middleware
	// ────────────────────────────────────────────────────────────────────────
	// Why?
	// - Next.js middleware runs on Edge Runtime
	// - Edge Runtime cannot decrypt httpOnly cookies
	// - We need a plain flag to check auth status
	//
	// Security:
	// - This cookie only contains a boolean flag (not sensitive)
	// - Real authentication still uses httpOnly cookies
	// - This is the industry-standard approach for edge middleware
	//
	// Parameters:
	// - name: "is_logged_in"
	// - value: "true"
	// - maxAge: 604800 seconds (7 days, matches session_id)
	// - path: "/" (available on all routes)
	// - domain: "" (current domain)
	// - secure: false (set to true in production with HTTPS)
	// - httpOnly: false (✅ CRITICAL: Must be readable by middleware)
	// ────────────────────────────────────────────────────────────────────────
	c.SetCookie("is_logged_in", "true", 604800, "/", "", false, false)
}

// clearAuthCookies clears all auth cookies for the current role.
// Also clears is_logged_in flag.
func clearAuthCookies(c *gin.Context) {
	prefix := cookiePrefix(c)

	// ────────────────────────────────────────────────────────────────────────
	// Clear httpOnly cookies
	// ────────────────────────────────────────────────────────────────────────
	c.SetCookie(prefix+"session_id", "", -1, "/", "", false, true)
	c.SetCookie(prefix+"access_token", "", -1, "/", "", false, true)
	c.SetCookie(prefix+"refresh_token", "", -1, "/", "", false, true)

	// ────────────────────────────────────────────────────────────────────────
	// ✅ NEW: Clear is_logged_in flag
	// ────────────────────────────────────────────────────────────────────────
	// CRITICAL: Must clear this flag so Next.js middleware knows user is logged out
	// Setting maxAge to -1 tells the browser to delete the cookie
	// ────────────────────────────────────────────────────────────────────────
	c.SetCookie("is_logged_in", "", -1, "/", "", false, false)
}

// getAuthCookie reads a cookie by name with role prefix.
func getAuthCookie(c *gin.Context, name string) (string, error) {
	prefix := cookiePrefix(c)
	return c.Cookie(prefix + name)
}
