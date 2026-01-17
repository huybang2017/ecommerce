package handler

import (
	"log"
	"net/http"
	"order-service/internal/domain"
	"order-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CartHandler handles HTTP requests for cart operations
// This is the transport layer - it knows HOW to handle HTTP (Gin framework)
// It delegates business logic to the service layer
type CartHandler struct {
	cartService *service.CartService
	logger      *zap.Logger
}

// NewCartHandler creates a new cart handler
// Dependency injection: we inject the service
func NewCartHandler(cartService *service.CartService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		logger:      logger,
	}
}

// AddItemRequest represents the request body for adding an item to cart
type AddItemRequest struct {
	ProductItemID uint `json:"product_item_id,omitempty"`
	Quantity      int  `json:"quantity" binding:"required,min=1"`
}

// UpdateItemRequest represents the request body for updating item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// GetCart handles GET /cart
// @Summary Get cart
// @Description Get the shopping cart for the current user
// @Tags Cart
// @Security UserAuth
// @Produce json
// @Success 200 {object} domain.ShoppingCart "Cart retrieved successfully"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	// Get user_id from header (set by API Gateway after JWT validation)
	userID := c.GetHeader("X-User-Id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItem handles POST /cart/items
// @Summary Add item to cart
// @Description Add a product item (SKU) to the shopping cart
// @Tags Cart
// @Security UserAuth
// @Accept json
// @Produce json
// @Param request body AddItemRequest true "Add Item Request"
// @Success 200 {object} handler.SuccessResponse "Item added successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request payload"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	// Get user_id from header (set by API Gateway after JWT validation)
	userID := c.GetHeader("X-User-Id")
	log.Println("")
	log.Println("userID:", userID)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use SKU-level ProductItemID for cart
	if req.ProductItemID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_item_id is required"})
		return
	}

	if err := h.cartService.AddToCart(
		c.Request.Context(),
		userID,
		req.ProductItemID,
		req.Quantity,
	); err != nil {
		h.logger.Error("failed to add item to cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// UpdateItem handles PUT /cart/items/:product_item_id
// @Summary Update item quantity
// @Description Update the quantity of an item in the cart
// @Tags Cart
// @Security UserAuth
// @Accept json
// @Produce json
// @Param product_item_id path int true "Product Item ID (SKU)"
// @Param request body UpdateItemRequest true "Update Item Request"
// @Success 200 {object} handler.SuccessResponse "Item updated successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request payload"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 404 {object} handler.ErrorResponse "Item not found"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/items/{product_item_id} [put]
func (h *CartHandler) UpdateItem(c *gin.Context) {
	// Get user_id from header (set by API Gateway after JWT validation)
	userID := c.GetHeader("X-User-Id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productItemIDStr := c.Param("product_item_id")
	productItemIDUint, err := strconv.ParseUint(productItemIDStr, 10, 32)
	if err != nil || productItemIDUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_item_id"})
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.UpdateItemQuantity(
		c.Request.Context(),
		userID,
		uint(productItemIDUint),
		req.Quantity,
	); err != nil {
		if err.Error() == "item not found in cart" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to update item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// RemoveItem handles DELETE /cart/items/:product_item_id
// @Summary Remove item from cart
// @Description Remove an item from the shopping cart
// @Tags Cart
// @Security UserAuth
// @Produce json
// @Param product_item_id path int true "Product Item ID (SKU)"
// @Success 200 {object} handler.SuccessResponse "Item removed successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 404 {object} handler.ErrorResponse "Item not found"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/items/{product_item_id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	// Get user_id from header (set by API Gateway after JWT validation)
	userID := c.GetHeader("X-User-Id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productItemIDStr := c.Param("product_item_id")
	productItemIDUint, err := strconv.ParseUint(productItemIDStr, 10, 32)
	if err != nil || productItemIDUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_item_id"})
		return
	}

	if err := h.cartService.RemoveFromCart(
		c.Request.Context(),
		userID,
		uint(productItemIDUint),
	); err != nil {
		if err.Error() == "item not found in cart" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to remove item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed successfully"})
}

// ClearCart handles DELETE /cart
// @Summary Clear cart
// @Description Remove all items from the shopping cart
// @Tags Cart
// @Security UserAuth
// @Produce json
// @Success 200 {object} handler.SuccessResponse "Cart cleared successfully"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	// Get user_id from header (set by API Gateway after JWT validation)
	userID := c.GetHeader("X-User-Id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to clear cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// ClearSelectedItems handles DELETE /cart/selected
// @Summary Clear selected items
// @Description Remove only items that are marked selected (used after checkout)
// @Tags Cart
// @Security UserAuth
// @Produce json
// @Success 200 {object} handler.SuccessResponse "Selected items cleared successfully"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/selected [delete]
func (h *CartHandler) ClearSelectedItems(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.cartService.ClearSelectedItems(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to clear selected items", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Selected items cleared successfully"})
}

// ToggleItemSelection handles POST /cart/items/:product_item_id/toggle
// @Summary Toggle item selection
// @Description Toggle the selection state of a cart item
// @Tags Cart
// @Deprecated This endpoint is deprecated; use PATCH /cart/items/{product_item_id}/selection instead
// @Security UserAuth
// @Produce json
// @Param product_item_id path int true "Product Item ID (SKU)"
// @Success 200 {object} handler.SuccessResponse "Item selection toggled successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request parameters"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 404 {object} handler.ErrorResponse "Item not found"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/items/{product_item_id}/toggle [post]
func (h *CartHandler) ToggleItemSelection(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productItemIDStr := c.Param("product_item_id")
	productItemIDUint, err := strconv.ParseUint(productItemIDStr, 10, 32)
	if err != nil || productItemIDUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_item_id"})
		return
	}

	if err := h.cartService.ToggleItemSelection(c.Request.Context(), userID, uint(productItemIDUint)); err != nil {
		if err.Error() == "item not found in cart" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to toggle item selection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item selection toggled successfully"})
}

// SetItemSelection handles PATCH /cart/items/:product_item_id/selection
// @Summary Set item selection
// @Description Set the selection state (selected/deselected) for a cart item
// @Tags Cart
// @Security UserAuth
// @Accept json
// @Produce json
// @Param product_item_id path int true "Product Item ID (SKU)"
// @Param request body SelectionRequest true "Selection Request"
// @Success 200 {object} handler.SuccessResponse "Item selection updated successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request payload"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 404 {object} handler.ErrorResponse "Item not found"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/items/{product_item_id}/selection [patch]
func (h *CartHandler) SetItemSelection(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productItemIDStr := c.Param("product_item_id")
	productItemIDUint, err := strconv.ParseUint(productItemIDStr, 10, 32)
	if err != nil || productItemIDUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_item_id"})
		return
	}

	var req SelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Selected == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "selected is required"})
		return
	}

	if err := h.cartService.SetItemSelection(c.Request.Context(), userID, uint(productItemIDUint), *req.Selected); err != nil {
		if err == domain.ErrCartItemNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to set item selection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item selection updated successfully"})
}

// SetAllSelection handles PATCH /cart/selection
// @Summary Set selection for all items
// @Description Set selection state for all items in the cart
// @Tags Cart
// @Security UserAuth
// @Accept json
// @Produce json
// @Param request body SelectAllRequest true "Select All Request"
// @Success 200 {object} handler.SuccessResponse "Selection updated successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request payload"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/selection [patch]
func (h *CartHandler) SetAllSelection(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SelectAllRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Selected == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "selected is required"})
		return
	}

	if err := h.cartService.SelectAllItems(c.Request.Context(), userID, *req.Selected); err != nil {
		h.logger.Error("failed to set all selections", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Selection updated successfully"})
}

// SetShopSelection handles PATCH /cart/shops/:shop_id/selection
// @Summary Set selection for shop items
// @Description Set selection state for all items from a shop
// @Tags Cart
// @Security UserAuth
// @Accept json
// @Produce json
// @Param shop_id path int true "Shop ID"
// @Param request body SelectShopRequest true "Select Shop Request"
// @Success 200 {object} handler.SuccessResponse "Selection updated successfully"
// @Failure 400 {object} handler.ErrorResponse "Invalid request payload"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/shops/{shop_id}/selection [patch]
func (h *CartHandler) SetShopSelection(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	shopIDStr := c.Param("shop_id")
	shopIDUint, err := strconv.ParseUint(shopIDStr, 10, 32)
	if err != nil || shopIDUint == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop_id"})
		return
	}

	var req SelectShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Selected == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "selected is required"})
		return
	}

	if err := h.cartService.SelectShopItems(c.Request.Context(), userID, uint(shopIDUint), *req.Selected); err != nil {
		h.logger.Error("failed to set shop selection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shop selection updated successfully"})
}

// SelectAllRequest represents request to select/deselect all items
// Use pointer bool so `false` is accepted while still validating presence
type SelectAllRequest struct {
	Selected *bool `json:"selected" binding:"required"`
}

// SelectionRequest represents request to set selection for a single item
// Use pointer bool so `false` is accepted while still validating presence
type SelectionRequest struct {
	Selected *bool `json:"selected" binding:"required"`
}

// SelectAllItems handles POST /cart/select_all
// @Summary Select or deselect all items
// @Description Set selection state for all items in cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body SelectAllRequest true "Select All Request"
// @Success 200 {object} map[string]string "Selection updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/select_all [post]
// func (h *CartHandler) SelectAllItems(c *gin.Context) {
// 	userID := c.GetHeader("X-User-Id")
// 	if userID == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	var req SelectAllRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.cartService.SelectAllItems(c.Request.Context(), userID, req.Selected); err != nil {
// 		h.logger.Error("failed to select all items", zap.Error(err))
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Selection updated successfully"})
// }

// SelectShopRequest represents request to select/deselect shop items
// Use pointer bool so `false` is accepted while still validating presence
type SelectShopRequest struct {
	Selected *bool `json:"selected" binding:"required"`
}

// SelectShopItems handles POST /cart/shops/:shop_id/select
// @Summary Select or deselect all items of a shop
// @Description Set selection state for all items from a specific shop
// @Tags Cart
// @Accept json
// @Produce json
// @Param shop_id path int true "Shop ID"
// @Param request body SelectShopRequest true "Select Shop Request"
// @Success 200 {object} map[string]string "Selection updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/shops/{shop_id}/select [post]
// func (h *CartHandler) SelectShopItems(c *gin.Context) {
// 	userID := c.GetHeader("X-User-Id")
// 	if userID == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
// 		return
// 	}

// 	shopIDStr := c.Param("shop_id")
// 	shopIDUint, err := strconv.ParseUint(shopIDStr, 10, 32)
// 	if err != nil || shopIDUint == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop_id"})
// 		return
// 	}

// 	var req SelectShopRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := h.cartService.SelectShopItems(c.Request.Context(), userID, uint(shopIDUint), req.Selected); err != nil {
// 		h.logger.Error("failed to select shop items", zap.Error(err))
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Shop selection updated successfully"})
// }

// ValidateCart handles POST /cart/validate
// @Summary Validate cart
// @Description Validate selected items in the cart (ensure stock and availability) before checkout
// @Tags Cart
// @Security UserAuth
// @Produce json
// @Success 200 {object} handler.SuccessResponse "Cart is valid"
// @Failure 400 {object} handler.ErrorResponse "Cart invalid"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized"
// @Failure 500 {object} handler.ErrorResponse "Internal server error"
// @Router /cart/validate [post]
func (h *CartHandler) ValidateCart(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.cartService.ValidateCart(c.Request.Context(), userID); err != nil {
		if err == domain.ErrCartEmpty {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to validate cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart is valid"})
}

// HealthCheck handles GET /health
func (h *CartHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "order-service"})
}
