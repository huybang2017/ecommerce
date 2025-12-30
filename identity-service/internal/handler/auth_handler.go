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

	// ONLY set HttpOnly cookie for refresh_token (long-lived, 7 days)
	// access_token is returned in response body for frontend to store in memory
	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		604800, // 7 days
		"/",
		"",
		false, // secure (true in production)
		true,  // httpOnly
	)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "user registered successfully",
		"access_token": response.AccessToken, // SHORT-LIVED (15 min) - store in memory
		"user":         response.User,
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

	// ONLY set HttpOnly cookie for refresh_token (long-lived, 7 days)
	// access_token is returned in response body for frontend to store in memory
	c.SetCookie(
		"refresh_token",       // name
		response.RefreshToken, // value
		604800,                // maxAge in seconds (7 days)
		"/",                   // path
		"",                    // domain
		false,                 // secure (true in production with HTTPS)
		true,                  // httpOnly (prevents JavaScript access)
	)

	// Return access_token in response body + user info
	// Frontend will store access_token in memory (NOT localStorage)
	c.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": response.AccessToken, // SHORT-LIVED (15 min) - store in memory
		"user":         response.User,
	})
}

// RefreshToken handles POST /auth/refresh
// @Summary Refresh access token
// @Description Use refresh token from cookie to get a new access token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 401 {object} map[string]interface{} "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("access_token")
	if err != nil || refreshToken == "" {
		h.logger.Warn("refresh token not found in cookie")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
		return
	}

	// Refresh access token
	response, err := h.authService.RefreshAccessToken(refreshToken)
	if err != nil {
		h.logger.Error("failed to refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return new access_token in response body (frontend stores in memory)
	// refresh_token cookie remains unchanged
	c.JSON(http.StatusOK, gin.H{
		"message":      "token refreshed successfully",
		"access_token": response.AccessToken, // NEW access token
		"user":         response.User,
	})
}

// Logout handles POST /auth/logout
// @Summary Logout user
// @Description Revoke all refresh tokens and clear cookies
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
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

	// Clear only refresh_token cookie (access_token is in memory, will be discarded by frontend)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
