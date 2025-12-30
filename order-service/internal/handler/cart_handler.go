package handler

import (
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
	ProductID     uint   `json:"product_id" binding:"required"`
	ProductItemID uint   `json:"product_item_id,omitempty"` // THÊM MỚI - SKU ID
	Name          string `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"required,min=0"`
	Quantity      int     `json:"quantity" binding:"required,min=1"`
	Image         string  `json:"image,omitempty"`
	SKU           string  `json:"sku,omitempty"`
}

// UpdateItemRequest represents the request body for updating item quantity
type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

// GetCart handles GET /cart
// @Summary Get cart
// @Description Get the shopping cart for the current user or session
// @Tags Cart
// @Produce json
// @Param user_id query string false "User ID (if authenticated)"
// @Param session_id query string false "Session ID (if guest)"
// @Success 200 {object} domain.Cart "Cart retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.Query("user_id")
	sessionID := c.Query("session_id")

	// If no user_id or session_id, try to get from headers or generate session
	if userID == "" && sessionID == "" {
		// For now, generate a session ID if not provided
		// In production, this would come from cookies or JWT token
		sessionID = c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id is required"})
			return
		}
	}

	cart, err := h.cartService.GetCart(c.Request.Context(), userID) // Đã sửa: chỉ userID
	if err != nil {
		h.logger.Error("failed to get cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItem handles POST /cart/items
// @Summary Add item to cart
// @Description Add a product to the shopping cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param user_id query string false "User ID (if authenticated)"
// @Param session_id query string false "Session ID (if guest)"
// @Param request body AddItemRequest true "Add Item Request"
// @Success 200 {object} domain.Cart "Item added successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	userID := c.Query("user_id")
	sessionID := c.Query("session_id")

	if userID == "" && sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id is required"})
			return
		}
	}

	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.AddItem(
		c.Request.Context(),
		userID, // Đã sửa: bỏ sessionID
		req.ProductID,
		req.Name,
		req.Price,
		req.Quantity,
		req.Image,
		req.SKU,
		req.ProductItemID, // THÊM MỚI - SKU ID
	)
	if err != nil {
		h.logger.Error("failed to add item to cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateItem handles PUT /cart/items/:product_id
// @Summary Update item quantity
// @Description Update the quantity of an item in the cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param user_id query string false "User ID (if authenticated)"
// @Param session_id query string false "Session ID (if guest)"
// @Param request body UpdateItemRequest true "Update Item Request"
// @Success 200 {object} domain.Cart "Item updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/items/{product_id} [put]
func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID := c.Query("user_id")
	sessionID := c.Query("session_id")

	if userID == "" && sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id is required"})
			return
		}
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.UpdateItemQuantity(
		c.Request.Context(),
		userID, // Đã sửa: bỏ sessionID
		uint(productID),
		req.Quantity,
	)
	if err != nil {
		if err.Error() == "item not found in cart" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to update item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// RemoveItem handles DELETE /cart/items/:product_id
// @Summary Remove item from cart
// @Description Remove an item from the shopping cart
// @Tags Cart
// @Produce json
// @Param product_id path int true "Product ID"
// @Param user_id query string false "User ID (if authenticated)"
// @Param session_id query string false "Session ID (if guest)"
// @Success 200 {object} domain.Cart "Item removed successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Item not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart/items/{product_id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID := c.Query("user_id")
	sessionID := c.Query("session_id")

	if userID == "" && sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id is required"})
			return
		}
	}

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	cart, err := h.cartService.RemoveItem(
		c.Request.Context(),
		userID, // Đã sửa: bỏ sessionID
		uint(productID),
	)
	if err != nil {
		if err.Error() == "item not found in cart" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to remove item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// ClearCart handles DELETE /cart
// @Summary Clear cart
// @Description Remove all items from the shopping cart
// @Tags Cart
// @Produce json
// @Param user_id query string false "User ID (if authenticated)"
// @Param session_id query string false "Session ID (if guest)"
// @Success 200 {object} map[string]string "Cart cleared successfully"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.Query("user_id")
	sessionID := c.Query("session_id")

	if userID == "" && sessionID == "" {
		sessionID = c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id is required"})
			return
		}
	}

	err := h.cartService.ClearCart(c.Request.Context(), userID) // Đã sửa: bỏ sessionID
	if err != nil {
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

