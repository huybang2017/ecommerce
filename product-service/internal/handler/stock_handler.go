package handler

import (
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StockHandler handles HTTP requests for stock operations
type StockHandler struct {
	stockService *service.StockService
	logger       *zap.Logger
}

// NewStockHandler creates a new stock handler
func NewStockHandler(stockService *service.StockService, logger *zap.Logger) *StockHandler {
	return &StockHandler{
		stockService: stockService,
		logger:       logger,
	}
}

// GetStock godoc
// @Summary Get stock for a product item
// @Description Get current stock quantity for a product item (SKU)
// @Tags stock
// @Produce json
// @Param id path int true "Product Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /product-items/{id}/stock [get]
func (h *StockHandler) GetStock(c *gin.Context) {
	productItemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_item_id"})
		return
	}

	stock, err := h.stockService.GetStock(c.Request.Context(), uint(productItemID))
	if err != nil {
		h.logger.Error("failed to get stock", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "product item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_item_id": productItemID,
		"stock":           stock,
	})
}

// CheckStock godoc
// @Summary Check stock availability
// @Description Check if enough stock is available for multiple items
// @Tags stock
// @Accept json
// @Produce json
// @Param request body domain.StockCheckRequest true "Stock check request"
// @Success 200 {object} domain.StockCheckResponse
// @Failure 400 {object} map[string]interface{}
// @Router /product-items/check-stock [post]
func (h *StockHandler) CheckStock(c *gin.Context) {
	var req domain.StockCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.stockService.CheckStock(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to check stock", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check stock"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ReserveStock godoc
// @Summary Reserve stock for an order
// @Description Temporarily reserve stock during checkout (15 minutes TTL)
// @Tags stock
// @Accept json
// @Produce json
// @Param request body domain.StockReserveRequest true "Stock reserve request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-items/reserve-stock [post]
func (h *StockHandler) ReserveStock(c *gin.Context) {
	var req domain.StockReserveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.stockService.ReserveStock(c.Request.Context(), &req); err != nil {
		h.logger.Error("failed to reserve stock", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "stock reserved successfully",
		"order_id": req.OrderID,
	})
}

// DeductStock godoc
// @Summary Deduct stock permanently
// @Description Deduct stock from product_item.qty_in_stock (after payment confirmed)
// @Tags stock
// @Accept json
// @Produce json
// @Param request body domain.StockDeductRequest true "Stock deduct request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-items/deduct-stock [post]
func (h *StockHandler) DeductStock(c *gin.Context) {
	var req domain.StockDeductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.stockService.DeductStock(c.Request.Context(), &req); err != nil {
		h.logger.Error("failed to deduct stock", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "stock deducted successfully",
		"order_id": req.OrderID,
	})
}

// ReleaseStock godoc
// @Summary Release reserved stock
// @Description Release stock reservation when order is cancelled or payment failed
// @Tags stock
// @Accept json
// @Produce json
// @Param request body domain.StockReleaseRequest true "Stock release request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-items/release-stock [post]
func (h *StockHandler) ReleaseStock(c *gin.Context) {
	var req domain.StockReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.stockService.ReleaseStock(c.Request.Context(), &req); err != nil {
		h.logger.Error("failed to release stock", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "stock reservations released successfully",
		"order_id": req.OrderID,
	})
}

// UpdateStock godoc
// @Summary Update stock quantity
// @Description Update stock quantity for a product item (for shop owners)
// @Tags stock
// @Accept json
// @Produce json
// @Param id path int true "Product Item ID"
// @Param request body map[string]int true "Stock update request {new_stock: 100}"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-items/{id}/stock [put]
func (h *StockHandler) UpdateStock(c *gin.Context) {
	productItemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_item_id"})
		return
	}

	var req struct {
		NewStock int `json:"new_stock" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.stockService.UpdateStock(c.Request.Context(), uint(productItemID), req.NewStock); err != nil {
		h.logger.Error("failed to update stock", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "stock updated successfully",
		"product_item_id": productItemID,
		"new_stock":       req.NewStock,
	})
}

