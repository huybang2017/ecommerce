package handler

import (
	"identity-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ShopHandler handles HTTP requests for shop operations
type ShopHandler struct {
	shopService *service.ShopService
	logger      *zap.Logger
}

// NewShopHandler creates a new shop handler
func NewShopHandler(shopService *service.ShopService, logger *zap.Logger) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
		logger:      logger,
	}
}

// CreateShop godoc
// @Summary Create a new shop
// @Description Create a new shop for a SELLER user (1 User = 1 Shop)
// @Tags shops
// @Accept json
// @Produce json
// @Param shop body service.CreateShopRequest true "Shop info"
// @Success 201 {object} domain.Shop
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /shops [post]
func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req service.CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_id from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	// Set owner_user_id from authenticated user
	req.OwnerUserID = userID.(uint)

	shop, err := h.shopService.CreateShop(&req)
	if err != nil {
		h.logger.Error("failed to create shop", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

// GetShop godoc
// @Summary Get shop by ID
// @Description Get shop details by shop ID
// @Tags shops
// @Produce json
// @Param id path int true "Shop ID"
// @Success 200 {object} domain.Shop
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /shops/{id} [get]
func (h *ShopHandler) GetShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
		return
	}

	shop, err := h.shopService.GetShop(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// GetMyShop godoc
// @Summary Get my shop
// @Description Get the shop of the authenticated user (1 User = 1 Shop)
// @Tags shops
// @Produce json
// @Success 200 {object} domain.Shop
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /shops/my-shop [get]
func (h *ShopHandler) GetMyShop(c *gin.Context) {
	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	shop, err := h.shopService.GetMyShop(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// ListShops godoc
// @Summary List all shops
// @Description Get all shops with pagination
// @Tags shops
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /shops [get]
func (h *ShopHandler) ListShops(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	shops, total, err := h.shopService.ListShops(page, limit)
	if err != nil {
		h.logger.Error("failed to list shops", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list shops"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shops": shops,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateShop godoc
// @Summary Update shop
// @Description Update shop information (only shop owner or ADMIN)
// @Tags shops
// @Accept json
// @Produce json
// @Param id path int true "Shop ID"
// @Param shop body service.UpdateShopRequest true "Shop info"
// @Success 200 {object} domain.Shop
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /shops/{id} [put]
func (h *ShopHandler) UpdateShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	var req service.UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shop, err := h.shopService.UpdateShop(uint(id), userID.(uint), &req)
	if err != nil {
		h.logger.Error("failed to update shop", zap.Error(err))
		if err.Error() == "only shop owner or ADMIN can update shop" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shop)
}

// DeleteShop godoc
// @Summary Delete shop
// @Description Soft delete shop by setting status to SUSPENDED (ADMIN only)
// @Tags shops
// @Produce json
// @Param id path int true "Shop ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /shops/{id} [delete]
func (h *ShopHandler) DeleteShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	if err := h.shopService.DeleteShop(uint(id), userID.(uint)); err != nil {
		h.logger.Error("failed to delete shop", zap.Error(err))
		if err.Error() == "only ADMIN can delete shop" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shop deleted successfully"})
}

// UpdateShopStatus godoc
// @Summary Update shop status
// @Description Update shop status (ADMIN only)
// @Tags shops
// @Accept json
// @Produce json
// @Param id path int true "Shop ID"
// @Param status body map[string]string true "Status" example({"status": "ACTIVE"})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /shops/{id}/status [put]
func (h *ShopHandler) UpdateShopStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
		return
	}

	// Get user_id from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.shopService.UpdateShopStatus(uint(id), req.Status, userID.(uint)); err != nil {
		h.logger.Error("failed to update shop status", zap.Error(err))
		if err.Error() == "only ADMIN can update shop status" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shop status updated successfully"})
}

