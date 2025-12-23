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
var _ = models.AuthResponse{}
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
// @Description Register a new user account with email, password, username, and full name
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} models.AuthResponse "User registered successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 409 {object} models.ErrorResponse "User already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// This will proxy to Identity Service
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user with email and password, receive JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} models.AuthResponse "Login successful"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// This will proxy to Identity Service
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package

