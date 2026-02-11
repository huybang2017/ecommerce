package handler

import (
	"fmt"
	"identity-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService *service.AuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register handles POST /auth/register
// @Summary Register a new user
// @Description Register a new user with email, password, username, and full name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 409 {object} map[string]interface{} "Email or username already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		h.logger.Error("failed to register", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set role-scoped HttpOnly cookies
	setAuthCookies(c, response.SessionID, response.AccessToken, response.RefreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"user":    response.User,
	})
}

// Login handles POST /auth/login
// @Summary Login user
// @Description Login with email and password, receive JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		h.logger.Error("failed to login", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set role-scoped HttpOnly cookies
	setAuthCookies(c, response.SessionID, response.AccessToken, response.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    response.User,
	})
}

// RefreshToken handles POST /auth/refresh
// @Summary Refresh access token
// @Description Use session_id from cookie to get a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 401 {object} map[string]interface{} "Invalid or expired session"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	prefix := cookiePrefix(c)

	// Try session_id cookie (role-scoped)
	sessionID, err := getAuthCookie(c, "session_id")
	if err == nil && sessionID != "" {
		response, err := h.authService.RefreshAccessTokenBySession(sessionID)
		if err != nil {
			h.logger.Error("failed to refresh token by session", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set new access_token with role prefix
		c.SetCookie(prefix+"access_token", response.AccessToken, 900, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "token refreshed successfully",
			"user":    response.User,
		})
		return
	}

	// Fallback: refresh_token cookie (role-scoped)
	refreshToken, err := getAuthCookie(c, "refresh_token")
	if err != nil || refreshToken == "" {
		h.logger.Warn("neither session_id nor refresh_token found in cookie")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session or refresh token required"})
		return
	}

	response, err := h.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		h.logger.Error("failed to refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(prefix+"access_token", response.AccessToken, 900, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed successfully",
		"user":    response.User,
	})
}

// Logout handles POST /auth/logout
// @Summary Logout user
// @Description Revoke session and clear cookies
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Try session_id from role-scoped cookie
	sessionID, err := getAuthCookie(c, "session_id")
	if err == nil && sessionID != "" {
		if err := h.authService.LogoutBySession(sessionID); err != nil {
			h.logger.Error("failed to logout session", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
			return
		}

		clearAuthCookies(c)

		c.JSON(http.StatusOK, gin.H{
			"message": "logout successful",
		})
		return
	}

	// Fallback: Legacy logout using user_id (revoke all refresh tokens)
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert user_id to uint
	uid, ok := userID.(uint)
	if !ok {
		// Try string conversion (from API Gateway)
		if userIDStr, ok := userID.(string); ok {
			var uidInt int
			if _, err := fmt.Sscanf(userIDStr, "%d", &uidInt); err == nil {
				uid = uint(uidInt)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id format"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}
	}

	// Revoke all refresh tokens
	if err := h.authService.Logout(uid); err != nil {
		h.logger.Error("failed to logout", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
