package handler

import (
	"net/http"
	"order-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderHandler handles HTTP requests for order operations
// This is the transport layer - it knows HOW to handle HTTP (Gin framework)
// It delegates business logic to the service layer
type OrderHandler struct {
	orderService *service.OrderService
	logger       *zap.Logger
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

// CreateOrder handles POST /orders
// @Summary Create order from cart
// @Description Create a new order from the shopping cart
// @Tags Order
// @Accept json
// @Produce json
// @Param order body service.CreateOrderRequest true "Order creation request"
// @Success 201 {object} domain.Order "Order created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Get user_id from query if not in body
	if req.UserID == nil {
		userIDStr := c.Query("user_id")
		if userIDStr != "" {
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err == nil {
				userIDUint := uint(userID)
				req.UserID = &userIDUint
			}
		}
	}

	// Get session_id from query if not in body
	if req.SessionID == "" {
		req.SessionID = c.Query("session_id")
	}

	order, err := h.orderService.CreateOrder(&req)
	if err != nil {
		h.logger.Error("failed to create order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /orders/:id
// @Summary Get order by ID
// @Description Get order details by order ID
// @Tags Order
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} domain.Order "Order retrieved successfully"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(uint(id))
	if err != nil {
		h.logger.Error("failed to get order", zap.Error(err), zap.Uint("order_id", uint(id)))
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderByOrderNumber handles GET /orders/number/:order_number
// @Summary Get order by order number
// @Description Get order details by order number
// @Tags Order
// @Produce json
// @Param order_number path string true "Order Number"
// @Success 200 {object} domain.Order "Order retrieved successfully"
// @Failure 404 {object} map[string]string "Order not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders/number/{order_number} [get]
func (h *OrderHandler) GetOrderByOrderNumber(c *gin.Context) {
	orderNumber := c.Param("order_number")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order number is required"})
		return
	}

	order, err := h.orderService.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		h.logger.Error("failed to get order", zap.Error(err), zap.String("order_number", orderNumber))
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListOrders handles GET /orders
// @Summary List orders
// @Description Get list of orders for a user or session
// @Tags Order
// @Produce json
// @Param user_id query int false "User ID"
// @Param session_id query string false "Session ID"
// @Param limit query int false "Limit (default: 20)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} map[string]interface{} "Orders retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Get user_id or session_id
	userIDStr := c.Query("user_id")
	sessionID := c.Query("session_id")

	if userIDStr == "" && sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id or session_id is required"})
		return
	}

	var userID *uint
	if userIDStr != "" {
		id, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		userIDUint := uint(id)
		userID = &userIDUint
	}

	// Get pagination params
	limit := 20
	offset := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	orders, total, err := h.orderService.ListOrders(userID, sessionID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list orders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

