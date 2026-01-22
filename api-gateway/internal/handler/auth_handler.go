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
func (h *AuthHandler) Register(c *gin.Context) {
	// This will proxy to Identity Service
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	// This will proxy to Identity Service
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Proxy to Identity Service - cookies will be forwarded automatically
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Proxy to Identity Service - auth middleware will add user_id to context
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package
