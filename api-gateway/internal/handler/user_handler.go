package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Import models for Swagger documentation generation
var _ = models.UserInfo{}
var _ = models.UpdateProfileRequest{}
var _ = models.ChangePasswordRequest{}
var _ = models.ErrorResponse{}

// UserHandler handles user-related requests
type UserHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(gatewayService *service.GatewayService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// GetProfile handles GET /users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateProfile handles PUT /users/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ChangePassword handles PUT /users/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package
