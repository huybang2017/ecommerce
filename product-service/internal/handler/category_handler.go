package handler

import (
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CategoryHandler handles HTTP requests for category operations
// This is the transport layer - it knows HOW to handle HTTP (Gin framework)
// It delegates business logic to the service layer
type CategoryHandler struct {
	categoryService *service.CategoryService
	logger          *zap.Logger
}

// NewCategoryHandler creates a new category handler
// Dependency injection: we inject the service
func NewCategoryHandler(categoryService *service.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	Description string `json:"description"`
}

// CreateCategory handles POST /categories
// @Summary Create a new category
// @Description Create a new category with name, slug, optional parent_id, and description
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to domain entity
	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		Description: req.Description,
	}

	// Call service layer
	if err := h.categoryService.CreateCategory(c.Request.Context(), category); err != nil {
		h.logger.Error("failed to create category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "category created successfully",
		"category": category,
	})
}

// UpdateCategory handles PUT /categories/:id
// @Summary Update an existing category
// @Description Update an existing category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or category ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing category
	category, err := h.categoryService.GetCategory(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	// Call service layer
	if err := h.categoryService.UpdateCategory(c.Request.Context(), category); err != nil {
		h.logger.Error("failed to update category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "category updated successfully",
		"category": category,
	})
}

// GetCategory handles GET /categories/:id
// @Summary Get a category by ID
// @Description Get a specific category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} handler.CategoryResponse "Category details"
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryBySlug handles GET /categories/slug/:slug
// @Summary Get a category by slug
// @Description Get a specific category by its slug
// @Tags Categories
// @Produce json
// @Param slug path string true "Category Slug"
// @Success 200 {object} handler.CategoryResponse "Category details"
// @Failure 400 {object} map[string]string "Slug is required"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetAllCategories handles GET /categories
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {array} handler.CategoryResponse "List of categories"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get all categories", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryChildren handles GET /categories/:id/children
// @Summary Get child categories
// @Description Get all child categories of a parent category
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {array} handler.CategoryResponse "List of child categories"
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	children, err := h.categoryService.GetCategoryChildren(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get category children", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, children)
}

// DeleteCategory handles DELETE /categories/:id
// @Summary Delete a category
// @Description Delete a category by its ID (cannot delete if has children)
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "Category deleted successfully"
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Internal server error or category has children"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

