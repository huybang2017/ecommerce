package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Import models for Swagger documentation generation
var _ = models.Category{}
var _ = models.CreateCategoryRequest{}
var _ = models.UpdateCategoryRequest{}
var _ = models.ProductsResponse{}
var _ = models.ErrorResponse{}
var _ = models.SuccessResponse{}

// CategoryHandler handles category-related requests
type CategoryHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(gatewayService *service.GatewayService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// ListCategories handles GET /categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategory handles GET /categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategoryBySlug handles GET /categories/slug/:slug
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategoryChildren handles GET /categories/:id/children
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategoryProducts handles GET /categories/:id/products
func (h *CategoryHandler) GetCategoryProducts(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateCategory handles PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteCategory handles DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
