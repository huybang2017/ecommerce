package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Import models for Swagger documentation generation
var _ = models.AddressInfo{}
var _ = models.CreateAddressRequest{}
var _ = models.UpdateAddressRequest{}
var _ = models.ErrorResponse{}

// AddressHandler handles address-related requests
type AddressHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(gatewayService *service.GatewayService, logger *zap.Logger) *AddressHandler {
	return &AddressHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// GetAddresses handles GET /addresses
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateAddress handles POST /addresses
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetAddress handles GET /addresses/:id
func (h *AddressHandler) GetAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateAddress handles PUT /addresses/:id
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteAddress handles DELETE /addresses/:id
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetDefaultAddress handles PUT /addresses/:id/default
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package
