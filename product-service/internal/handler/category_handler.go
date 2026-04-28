package handler

import (
	"net/http"
	"product-service/internal/dto/request"
	"product-service/internal/dto/response"
	"product-service/internal/mapper"
	"product-service/internal/service"
	"product-service/pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
	logger          *zap.Logger
}

func NewCategoryHandler(categoryService *service.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

// mapServiceError maps service layer errors to appropriate HTTP status codes
func (h *CategoryHandler) mapServiceError(c *gin.Context, err error, defaultMsg string) {
	errMsg := err.Error()
	errLower := strings.ToLower(errMsg)

	// Check for validation errors (400 Bad Request)
	if strings.Contains(errLower, "required") ||
		strings.Contains(errLower, "invalid") ||
		strings.Contains(errLower, "must") ||
		strings.Contains(errLower, "cannot") {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", errMsg)
		return
	}

	// Check for conflict errors (409 Conflict)
	if strings.Contains(errLower, "already exists") ||
		strings.Contains(errLower, "duplicate") {
		utils.Error(c, http.StatusConflict, "CONFLICT", errMsg)
		return
	}

	// Check for not found errors (404 Not Found)
	if strings.Contains(errLower, "not found") {
		utils.Error(c, http.StatusNotFound, "NOT_FOUND", errMsg)
		return
	}

	// Default to 500 Internal Server Error
	h.logger.Error(defaultMsg, zap.Error(err))
	utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", errMsg)
}

// CreateParentCategory handles POST /categories/admin/parent
// @Summary Create a new parent category (ADMIN only)
// @Description Admin creates a new parent category (no parent_id)
// @Tags Admin Categories
// @Accept json
// @Produce json
// @Param request body request.CreateParentCategoryRequest true "Create Parent Category Request"
// @Success 201 {object} response.CategoryResponse "Parent category created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload or validation error"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin role required"
// @Failure 409 {object} utils.ErrorResponse "Conflict - category with this slug already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/admin/parent [post]
func (h *CategoryHandler) CreateParentCategory(c *gin.Context) {
	var req request.CreateParentCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	category := mapper.CreateParentCategoryRequestToDomain(&req)

	if err := h.categoryService.CreateParentCategory(c.Request.Context(), category); err != nil {
		h.mapServiceError(c, err, "failed to create parent category")
		return
	}

	c.JSON(http.StatusCreated, mapper.CategoryToResponse(category))
}

// CreateChildCategory handles POST /categories/:id/children
// @Summary Create a child category under a parent
// @Description Creates a child category under the specified parent category. The parent category ID is provided in the URL path, NOT in the request body.
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Parent Category ID"
// @Param request body request.CreateChildCategoryRequest true "Child Category Properties (no parent_id)"
// @Success 201 {object} response.CategoryResponse "Child category created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload, parent ID, or validation error"
// @Failure 404 {object} utils.ErrorResponse "Parent category not found"
// @Failure 409 {object} utils.ErrorResponse "Conflict - category with this slug already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/{id}/children [post]
func (h *CategoryHandler) CreateChildCategory(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid parent category ID")
		return
	}

	var req request.CreateChildCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// Convert DTO to domain entity with parentID from path
	category := mapper.CreateChildCategoryRequestToDomain(&req, uint(parentID))

	// Call service layer
	if err := h.categoryService.CreateChildCategory(c.Request.Context(), category); err != nil {
		h.mapServiceError(c, err, "failed to create child category")
		return
	}

	c.JSON(http.StatusCreated, mapper.CategoryToResponse(category))
}

// CreateCategory handles POST /categories (backward compatibility)
// @Summary Create a new category (generic endpoint)
// @Description Create a category with optional parent_id in the request body. If parent_id is null or omitted, creates a parent (root) category. If parent_id is provided, creates a child category under that parent.
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body request.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} response.CategoryResponse "Category created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload or validation error"
// @Failure 404 {object} utils.ErrorResponse "Parent category not found (if parent_id provided)"
// @Failure 409 {object} utils.ErrorResponse "Conflict - category with this slug already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	category := mapper.CreateCategoryRequestToDomain(&req)

	// Call service layer
	if err := h.categoryService.CreateCategory(c.Request.Context(), category); err != nil {
		h.mapServiceError(c, err, "failed to create category")
		return
	}

	c.JSON(http.StatusCreated, mapper.CategoryToResponse(category))
}

// GetAdminCategoryParent handles GET /categories/admin/parents
// @Summary Get parent categories for admin
// @Description Retrieve a list of parent categories for admin management
// @Tags Admin Categories
// @Accept json
// @Produce json
// @Success 200 {array} response.CategoryResponse "List of parent categories"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/admin/parents [get]
func (h *CategoryHandler) GetAdminCategoryParent(c *gin.Context) {
	parents, err := h.categoryService.GetAdminCategoryParent(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get admin category parents", zap.Error(err))
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response.CategoryListResponse{
		Data: mapper.CategoriesToResponse(parents),
	})
}

// UpdateCategory handles PUT /categories/:id
// @Summary Update an existing category
// @Description Update an existing category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body request.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} response.CategoryResponse "Category updated successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request payload or category ID"
// @Failure 404 {object} utils.ErrorResponse "Category not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid category ID")
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// Get existing category
	category, err := h.categoryService.GetCategory(c.Request.Context(), uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "NOT_FOUND", "category not found")
		return
	}

	// Update fields using mapper
	category = mapper.UpdateCategoryRequestToDomain(category, &req)

	// Call service layer
	if err := h.categoryService.UpdateCategory(c.Request.Context(), category); err != nil {
		h.mapServiceError(c, err, "failed to update category")
		return
	}

	c.JSON(http.StatusOK, mapper.CategoryToResponse(category))
}

// PatchCategoryActive handles PATCH /categories/:id/active
// @Summary Update category active status
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body map[string]bool true "Active status" example({"is_active": true})
// @Success 200 {object} map[string]interface{}
// @Router /categories/{id}/active [patch]
func (h *CategoryHandler) PatchCategoryActive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid category ID")
		return
	}

	var req struct {
		IsActive bool `json:"is_active" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if err := h.categoryService.IsActiveCategory(c.Request.Context(), uint(id), req.IsActive); err != nil {
		h.logger.Error("failed to patch category status", zap.Uint("id", uint(id)), zap.Error(err))
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "category status updated successfully",
		"is_active": req.IsActive,
	})
}

// GetCategory handles GET /categories/:id
// @Summary Get a category by ID
// @Description Get a specific category by its ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.CategoryResponse "Category details"
// @Failure 400 {object} utils.ErrorResponse "Invalid category ID"
// @Failure 404 {object} utils.ErrorResponse "Category not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "NOT_FOUND", "category not found")
		return
	}

	c.JSON(http.StatusOK, mapper.CategoryToResponse(category))
}

// GetCategoryBySlug handles GET /categories/slug/:slug
// @Summary Get a category by slug
// @Description Get a specific category by its slug
// @Tags Categories
// @Produce json
// @Param slug path string true "Category Slug"
// @Success 200 {object} response.CategoryResponse "Category details"
// @Failure 400 {object} utils.ErrorResponse "Slug is required"
// @Failure 404 {object} utils.ErrorResponse "Category not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "slug is required")
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "NOT_FOUND", "category not found")
		return
	}

	c.JSON(http.StatusOK, mapper.CategoryToResponse(category))
}

// GetAllCategories handles GET /categories
// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} response.CategoryListResponse "List of categories"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get all categories", zap.Error(err))
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response.CategoryListResponse{
		Data: mapper.CategoriesToResponse(categories),
	})
}

// GetAdminCategories handles GET /admin/categories
// @Summary Get categories for admin with pagination, filter and sort
// @Tags Admin Categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by name or slug"
// @Param parent_id query int false "Filter by parent ID"
// @Param is_active query bool false "Filter by active status"
// @Param sort_by query string false "Sort field (id, name, created_at)" default(id)
// @Param sort_order query string false "Sort order (asc/desc)" default(desc)
// @Success 200 {object} response.AdminCategoryListResponse
// @Router /categories/admin [get]
func (h *CategoryHandler) GetAdminCategories(c *gin.Context) {
	var req request.CategoryQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	// Convert DTO to domain params
	params := mapper.CategoryQueryRequestToDomain(&req)

	// Call service
	categories, total, totalPages, err := h.categoryService.GetAdminCategories(c.Request.Context(), params)
	if err != nil {
		h.logger.Error("failed to get admin categories", zap.Error(err))
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	// Return response with pagination metadata
	c.JSON(http.StatusOK, response.AdminCategoryListResponse{
		Data:       mapper.CategoriesToResponse(categories),
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

// GetCategoryChildren handles GET /categories/:id/children
// @Summary Get child categories
// @Description Get all child categories of a parent category
// @Tags Categories
// @Produce json
// @Param id path int true "Parent Category ID"
// @Success 200 {object} response.CategoryListResponse "List of child categories"
// @Failure 400 {object} utils.ErrorResponse "Invalid category ID"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid category ID")
		return
	}

	children, err := h.categoryService.GetCategoryChildren(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("failed to get category children", zap.Error(err))
		utils.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response.CategoryListResponse{
		Data: mapper.CategoriesToResponse(children),
	})
}

// DeleteCategory handles DELETE /categories/:id
// @Summary Delete a category
// @Description Delete a category by its ID (cannot delete if has children)
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string "Category deleted successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid category ID"
// @Failure 404 {object} utils.ErrorResponse "Category not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error or category has children"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "BAD_REQUEST", "invalid category ID")
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		h.mapServiceError(c, err, "failed to delete category")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

// @Summary Get category tree
// @Description Get all categories in a hierarchical tree structure
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} response.CategoryTreeResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get category tree", zap.Error(err))
		utils.Error(
			c,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Failed to get category tree",
		)
		return
	}
	c.JSON(http.StatusOK, response.CategoryTreeResponse{
		Data: mapper.BuildCategoryTree(categories),
	})
}
