package handler

import (
	"log"
	"net/http"
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
// @Produce json
// @Param user_id query string true "User ID (authenticated)"
// @Success 200 {object} domain.Cart "Cart retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
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
// @Accept json
// @Produce json
// @Param user_id query string true "User ID (authenticated)"
// @Param request body AddItemRequest true "Add Item Request"
// @Success 200 {object} map[string]string "Item added successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	userID := c.Query("user_id")
	log.Println("userID:", userID)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
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
// @Accept json
// @Produce json
// @Param product_item_id path int true "Product Item ID (SKU)"
// @Param user_id query string true "User ID (authenticated)"
// @Param request body UpdateItemRequest true "Update Item Request"
// @Success 200 {object} map[string]string "Item updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/items/{product_item_id} [put]
func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
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
// @Produce json
// @Param product_item_id path int true "Product Item ID (SKU)"
// @Param user_id query string true "User ID (authenticated)"
// @Success 200 {object} map[string]string "Item removed successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/items/{product_item_id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
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
// @Produce json
// @Param user_id query string true "User ID (authenticated)"
// @Success 200 {object} map[string]string "Cart cleared successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if err := h.cartService.ClearCart(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to clear cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

// HealthCheck handles GET /health
func (h *CartHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "order-service"})
}
