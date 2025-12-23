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
// @Summary Get all addresses
// @Description Get all addresses for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.AddressInfo "List of addresses"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /addresses [get]
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateAddress handles POST /addresses
// @Summary Create a new address
// @Description Create a new address for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateAddressRequest true "Create Address Request"
// @Success 201 {object} models.AddressInfo "Address created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /addresses [post]
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetAddress handles GET /addresses/:id
// @Summary Get address by ID
// @Description Get a specific address by ID for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} models.AddressInfo "Address details"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Address not found"
// @Router /addresses/{id} [get]
func (h *AddressHandler) GetAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateAddress handles PUT /addresses/:id
// @Summary Update an address
// @Description Update an existing address for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Param request body models.UpdateAddressRequest true "Update Address Request"
// @Success 200 {object} models.AddressInfo "Address updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Address not found"
// @Router /addresses/{id} [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteAddress handles DELETE /addresses/:id
// @Summary Delete an address
// @Description Delete an address for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} models.SuccessResponse "Address deleted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Address not found"
// @Router /addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetDefaultAddress handles PUT /addresses/:id/default
// @Summary Set default address
// @Description Set an address as the default address for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} models.AddressInfo "Address set as default"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Address not found"
// @Router /addresses/{id}/default [put]
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// Models are now in api-gateway/internal/models package

