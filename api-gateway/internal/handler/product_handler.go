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

