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
func (h *ProductHandler) ListProducts(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetProduct handles GET /products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateProduct handles PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteProduct handles DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SearchProducts handles GET /products/search
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateInventory handles PATCH /products/:id/inventory
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ==================== PRODUCT ITEMS (SKU) ====================

// GetProductItems handles GET /products/:id/items
func (h *ProductHandler) GetProductItems(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetProductItem handles GET /products/:id/items/:item_id
func (h *ProductHandler) GetProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateProductItem handles POST /products/:id/items
func (h *ProductHandler) CreateProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateProductItem handles PUT /products/:id/items/:item_id
func (h *ProductHandler) UpdateProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteProductItem handles DELETE /products/:id/items/:item_id
func (h *ProductHandler) DeleteProductItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetProductVariations handles GET /products/:id/variations
func (h *ProductHandler) GetProductVariations(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
