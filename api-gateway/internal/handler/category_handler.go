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
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category "List of categories"
// @Router /categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategory handles GET /categories/:id
// @Summary Get category by ID
// @Description Get detailed information about a specific category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Category "Category details"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategoryBySlug handles GET /categories/slug/:slug
// @Summary Get category by slug
// @Description Get category information by its slug
// @Tags Categories
// @Accept json
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} models.Category "Category details"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategoryChildren handles GET /categories/:id/children
// @Summary Get child categories
// @Description Get all child categories of a parent category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Parent Category ID"
// @Success 200 {array} models.Category "List of child categories"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetCategoryProducts handles GET /categories/:id/products
// @Summary Get products by category
// @Description Get all products in a specific category with pagination
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} models.ProductsResponse "List of products in category"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Router /categories/{id}/products [get]
func (h *CategoryHandler) GetCategoryProducts(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// CreateCategory handles POST /categories
// @Summary Create a new category
// @Description Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body models.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} models.SuccessResponse "Category created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateCategory handles PUT /categories/:id
// @Summary Update an existing category
// @Description Update an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body models.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} models.SuccessResponse "Category updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// DeleteCategory handles DELETE /categories/:id
// @Summary Delete a category
// @Description Delete a category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.SuccessResponse "Category deleted successfully"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

