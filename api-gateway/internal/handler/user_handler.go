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
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserInfo "User profile"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateProfile handles PUT /users/profile
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} models.UserInfo "Profile updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ChangePassword handles PUT /users/password
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} models.SuccessResponse "Password changed successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package

