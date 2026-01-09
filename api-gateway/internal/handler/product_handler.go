package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Import models for Swagger documentation generation
var _ = models.Product{}
var _ = models.CreateProductRequest{}
var _ = models.UpdateProductRequest{}
var _ = models.ProductsResponse{}
var _ = models.ErrorResponse{}
var _ = models.SuccessResponse{}

// ProductHandler handles product-related requests
type ProductHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewProductHandler creates a new product handler
func NewProductHandler(gatewayService *service.GatewayService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// ListProducts handles GET /products
// @Summary List products with pagination and filters
// @Description Get a paginated list of products with optional filters (category, status)
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param category_id query int false "Filter by category ID"
// @Param status query string false "Filter by status (ACTIVE, INACTIVE, DRAFT)"
// @Success 200 {object} models.ProductsResponse "List of products"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetProduct handles GET /products/:id
// @Summary Get product by ID
// @Description Get detailed information about a specific product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product "Product details"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateProduct handles POST /products
// @Summary Create a new product
// @Description Create a new product (requires authentication)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateProductRequest true "Create Product Request"
// @Success 201 {object} models.SuccessResponse "Product created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateProduct handles PUT /products/:id
// @Summary Update an existing product
// @Description Update an existing product (requires authentication)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body models.UpdateProductRequest true "Update Product Request"
// @Success 200 {object} models.SuccessResponse "Product updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteProduct handles DELETE /products/:id
// @Summary Delete a product
// @Description Delete a product (requires authentication)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} models.SuccessResponse "Product deleted successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SearchProducts handles GET /products/search
// @Summary Search products
// @Description Search products by query string and optional category
// @Tags Products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param category query string false "Filter by category slug"
// @Success 200 {array} models.Product "List of matching products"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateInventory handles PATCH /products/:id/inventory
// @Summary Update product inventory
// @Description Update product stock quantity (requires authentication)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param stock body object true "Stock update" example({"stock": 100})
// @Success 200 {object} models.SuccessResponse "Inventory updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /products/{id}/inventory [patch]
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ==================== PRODUCT ITEMS (SKU) ====================

// GetProductItems handles GET /products/:id/items
// @Summary Get all product items (SKUs) for a product
// @Description Get list of all SKUs/variations for a specific product
// @Tags Product Items
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} object "List of product items"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /products/{id}/items [get]
func (h *ProductHandler) GetProductItems(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetProductItem handles GET /products/:id/items/:item_id
// @Summary Get specific product item (SKU)
// @Description Get details of a specific SKU by product ID and item ID
// @Tags Product Items
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param item_id path int true "Product Item ID"
// @Success 200 {object} object "Product item details"
// @Failure 404 {object} models.ErrorResponse "Product item not found"
// @Router /products/{id}/items/{item_id} [get]
func (h *ProductHandler) GetProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateProductItem handles POST /products/:id/items
// @Summary Create new product item (SKU)
// @Description Create a new SKU for a product (Auth required)
// @Tags Product Items
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 201 {object} object "Created product item"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /products/{id}/items [post]
func (h *ProductHandler) CreateProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateProductItem handles PUT /products/:id/items/:item_id
// @Summary Update product item (SKU)
// @Description Update an existing SKU (Auth required)
// @Tags Product Items
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param item_id path int true "Product Item ID"
// @Success 200 {object} object "Updated product item"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Product item not found"
// @Router /products/{id}/items/{item_id} [put]
func (h *ProductHandler) UpdateProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteProductItem handles DELETE /products/:id/items/:item_id
// @Summary Delete product item (SKU)
// @Description Delete a SKU (Auth required)
// @Tags Product Items
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param item_id path int true "Product Item ID"
// @Success 204 "Product item deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Product item not found"
// @Router /products/{id}/items/{item_id} [delete]
func (h *ProductHandler) DeleteProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetProductVariations handles GET /products/:id/variations
// @Summary Get product variations with options
// @Description Get all variations (Color, Size, etc.) with their options for a product
// @Tags Product Variations
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} object "List of variations with options"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Router /products/{id}/variations [get]
func (h *ProductHandler) GetProductVariations(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
