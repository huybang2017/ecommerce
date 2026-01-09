package handler

import (
	"net/http"
	"product-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SKUHandler handles HTTP requests for SKU-related operations (variations, SKUs)
type SKUHandler struct {
	productItemService *service.ProductItemService
	logger             *zap.Logger
}

// ProductItemWithVariations extends ProductItem with variation option IDs
type ProductItemWithVariations struct {
	ID                 uint    `json:"id"`
	ProductID          uint    `json:"product_id"`
	SKUCode            string  `json:"sku_code"`
	ImageURL           string  `json:"image_url"`
	Price              float64 `json:"price"`
	QtyInStock         int     `json:"qty_in_stock"`
	Status             string  `json:"status"`
	VariationOptionIDs []uint  `json:"variation_option_ids"` // [1, 5] = Size M + Color Red
}

// NewSKUHandler creates a new SKU handler
func NewSKUHandler(productItemService *service.ProductItemService, logger *zap.Logger) *SKUHandler {
	return &SKUHandler{
		productItemService: productItemService,
		logger:             logger,
	}
}

// CreateProductItem godoc
// @Summary Create a new product item (SKU)
// @Description Create a new SKU for a product with variation options
// @Tags skus
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param item body service.CreateProductItemRequest true "Product item info"
// @Success 201 {object} domain.ProductItem
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/items [post]
func (h *SKUHandler) CreateProductItem(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	var req service.CreateProductItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set product_id from path
	req.ProductID = uint(productID)

	item, err := h.productItemService.CreateProductItem(&req)
	if err != nil {
		h.logger.Error("failed to create product item", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetProductItems godoc
// @Summary Get all SKUs for a product
// @Description Get all product items (SKUs) for a specific product
// @Tags skus
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {array} domain.ProductItem
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/items [get]
func (h *SKUHandler) GetProductItems(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
		return
	}

	items, err := h.productItemService.GetProductItemsWithVariations(uint(productID))
	if err != nil {
		h.logger.Error("failed to get product items", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// GetProductItem godoc
// @Summary Get a specific SKU
// @Description Get product item (SKU) details by ID
// @Tags skus
// @Produce json
// @Param product_id path int true "Product ID"
// @Param item_id path int true "Product Item ID"
// @Success 200 {object} domain.ProductItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/items/{item_id} [get]
func (h *SKUHandler) GetProductItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	item, err := h.productItemService.GetProductItem(uint(itemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetProductItemBySKU godoc
// @Summary Get SKU by code
// @Description Get product item by SKU code
// @Tags skus
// @Produce json
// @Param sku_code path string true "SKU Code"
// @Success 200 {object} domain.ProductItem
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-items/{sku_code} [get]
func (h *SKUHandler) GetProductItemBySKU(c *gin.Context) {
	skuCode := c.Param("sku_code")

	item, err := h.productItemService.GetProductItemBySKU(skuCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetProductItemsBatch godoc
// @Summary Get multiple product items by IDs (batch)
// @Description Fetch multiple product items in one request for cart/order services
// @Tags skus
// @Produce json
// @Param ids query string true "Comma-separated product item IDs (e.g., 1,2,3)"
// @Success 200 {object} map[string]interface{} "items array with product details"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product-items/batch [get]
func (h *SKUHandler) GetProductItemsBatch(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids parameter is required"})
		return
	}

	// Parse comma-separated IDs
	var ids []uint
	idStrings := splitByComma(idsParam)
	for _, idStr := range idStrings {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format: " + idStr})
			return
		}
		ids = append(ids, uint(id))
	}

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid ids provided"})
		return
	}

	// Fetch items with product details
	items, err := h.productItemService.GetProductItemsWithProduct(ids)
	if err != nil {
		h.logger.Error("failed to get product items batch", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch product items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

// Helper function to split comma-separated string
func splitByComma(s string) []string {
	var result []string
	current := ""
	for _, char := range s {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if char != ' ' { // Skip spaces
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// UpdateProductItem godoc
// @Summary Update a SKU
// @Description Update product item (SKU) details
// @Tags skus
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param item_id path int true "Product Item ID"
// @Param item body service.UpdateProductItemRequest true "Product item info"
// @Success 200 {object} domain.ProductItem
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/items/{item_id} [put]
func (h *SKUHandler) UpdateProductItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	var req service.UpdateProductItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.productItemService.UpdateProductItem(uint(itemID), &req)
	if err != nil {
		h.logger.Error("failed to update product item", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteProductItem godoc
// @Summary Delete a SKU
// @Description Delete product item (SKU) and its configurations
// @Tags skus
// @Produce json
// @Param product_id path int true "Product ID"
// @Param item_id path int true "Product Item ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{product_id}/items/{item_id} [delete]
func (h *SKUHandler) DeleteProductItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	if err := h.productItemService.DeleteProductItem(uint(itemID)); err != nil {
		h.logger.Error("failed to delete product item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product item deleted successfully"})
}
