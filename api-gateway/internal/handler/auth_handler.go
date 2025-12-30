package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Import models for Swagger documentation generation
var _ = models.RegisterRequest{}
var _ = models.LoginRequest{}
var _ = models.LoginResponse{}
var _ = models.RegisterResponse{}
var _ = models.RefreshResponse{}
var _ = models.ErrorResponse{}

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(gatewayService *service.GatewayService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account. Returns access_token in response body (store in memory) and sets refresh_token as HttpOnly cookie (7 days). The Gateway forwards the request to Identity Service which validates data, creates user, and returns tokens. Frontend should store access_token in memory (NOT localStorage) for security.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data (email, password, username, full_name)"
// @Success 201 {object} models.RegisterResponse "User registered successfully. Response contains access_token (15min, store in memory) and user info. refresh_token sent as HttpOnly cookie (7 days)"
// @Header 201 {string} Set-Cookie "refresh_token=<jwt>; Path=/; Max-Age=604800; HttpOnly"
// @Failure 400 {object} models.ErrorResponse "Bad request - invalid input format"
// @Failure 409 {object} models.ErrorResponse "User already exists - email or username taken"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// This will proxy to Identity Service
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user with email and password. Returns access_token in response body (15 min, store in memory) and sets refresh_token as HttpOnly cookie (7 days). The Gateway forwards credentials to Identity Service which validates them and returns tokens. Frontend should store access_token in memory (NOT localStorage) and send it via Authorization header. refresh_token cookie is sent automatically by browser.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials (email and password)"
// @Success 200 {object} models.LoginResponse "Login successful. Response contains access_token (15min, store in memory) and user info. refresh_token sent as HttpOnly cookie (7 days)"
// @Header 200 {string} Set-Cookie "refresh_token=<jwt>; Path=/; Max-Age=604800; HttpOnly"
// @Failure 400 {object} models.ErrorResponse "Bad request - invalid input format"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials - wrong email or password"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// This will proxy to Identity Service
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Use refresh_token from HttpOnly cookie to obtain a new access_token. The Gateway automatically forwards the cookie to Identity Service which validates the refresh token and issues a new access_token (15 min) in response body. The refresh_token cookie remains unchanged. Frontend should update the access_token in memory.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Cookie header string true "refresh_token cookie (HttpOnly, automatically sent by browser)" default(refresh_token=<your_refresh_token>)
// @Success 200 {object} models.RefreshResponse "Token refreshed successfully. New access_token returned in body (store in memory). refresh_token cookie unchanged"
// @Failure 401 {object} models.ErrorResponse "Invalid, expired, or revoked refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Proxy to Identity Service - cookies will be forwarded automatically
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Logout handles user logout
// @Summary Logout user
// @Description Revoke all refresh tokens for the authenticated user and clear refresh_token cookie. Requires valid access_token in Authorization header (Bearer token). The Gateway forwards the request to Identity Service which revokes all refresh tokens from the database. Frontend should discard the access_token from memory.
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer <access_token>" default(Bearer <your_access_token>)
// @Success 200 {object} map[string]interface{} "Logout successful. All refresh tokens revoked and refresh_token cookie cleared. Frontend should discard access_token from memory"
// @Header 200 {string} Set-Cookie "refresh_token=; Path=/; Max-Age=-1; HttpOnly (cleared)"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - missing or invalid access token"
// @Failure 500 {object} models.ErrorResponse "Internal server error - failed to revoke tokens"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Proxy to Identity Service - auth middleware will add user_id to context
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package
