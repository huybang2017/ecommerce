package handler

import (
	"identity-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AddressHandler handles HTTP requests for address operations
type AddressHandler struct {
	addressService *service.AddressService
	logger         *zap.Logger
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(addressService *service.AddressService, logger *zap.Logger) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		logger:         logger,
	}
}

// CreateAddress handles POST /addresses
// @Summary Create a new address
// @Description Create a new address for the current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateAddressRequest true "Address data"
// @Success 201 {object} map[string]interface{} "Address created successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /addresses [post]
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	var req service.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid create address request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.addressService.CreateAddress(userIDUint, &req)
	if err != nil {
		h.logger.Error("failed to create address", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "address created successfully",
		"data":    address,
	})
}

// GetAddresses handles GET /addresses
// @Summary Get all addresses
// @Description Get all addresses for the current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of addresses"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /addresses [get]
func (h *AddressHandler) GetAddresses(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	addresses, err := h.addressService.GetAddresses(userIDUint)
	if err != nil {
		h.logger.Error("failed to get addresses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": addresses,
	})
}

// GetAddress handles GET /addresses/:id
// @Summary Get address by ID
// @Description Get a specific address by ID for the current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} map[string]interface{} "Address details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Address not found"
// @Router /addresses/{id} [get]
func (h *AddressHandler) GetAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	address, err := h.addressService.GetAddress(userIDUint, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": address,
	})
}

// UpdateAddress handles PUT /addresses/:id
// @Summary Update address
// @Description Update an existing address for the current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Param request body service.UpdateAddressRequest true "Address update data"
// @Success 200 {object} map[string]interface{} "Address updated successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /addresses/{id} [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	var req service.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update address request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address, err := h.addressService.UpdateAddress(userIDUint, uint(id), &req)
	if err != nil {
		h.logger.Error("failed to update address", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "address updated successfully",
		"data":    address,
	})
}

// DeleteAddress handles DELETE /addresses/:id
// @Summary Delete address
// @Description Delete an address for the current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} map[string]interface{} "Address deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /addresses/{id} [delete]
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	if err := h.addressService.DeleteAddress(userIDUint, uint(id)); err != nil {
		h.logger.Error("failed to delete address", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "address deleted successfully",
	})
}

// SetDefaultAddress handles PUT /addresses/:id/default
// @Summary Set default address
// @Description Set an address as the default address for the current user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} map[string]interface{} "Default address set successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /addresses/{id}/default [put]
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")
	userIDUint := userID.(uint)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address ID"})
		return
	}

	if err := h.addressService.SetDefaultAddress(userIDUint, uint(id)); err != nil {
		h.logger.Error("failed to set default address", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "default address set successfully",
	})
}


