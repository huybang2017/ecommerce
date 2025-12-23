package handler

import (
	"encoding/json"
	"net/http"
	"product-service/internal/domain"
	"product-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// ProductHandler handles HTTP requests for product operations
// This is the transport layer - it knows HOW to handle HTTP (Gin framework)
// It delegates business logic to the service layer
type ProductHandler struct {
	productService *service.ProductService
	logger         *zap.Logger
}

// NewProductHandler creates a new product handler
// Dependency injection: we inject the service
func NewProductHandler(productService *service.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,min=0"`
	SKU         string   `json:"sku" binding:"required"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	Status      string   `json:"status"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock"`
	IsActive    bool     `json:"is_active"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"min=0"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	Status      string   `json:"status"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock"`
	IsActive    *bool    `json:"is_active"`
}

// ProductResponse represents the product response for Swagger
type ProductResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	SKU         string   `json:"sku"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	Status      string   `json:"status"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// CategoryResponse represents the category response for Swagger
type CategoryResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	ParentID    *uint   `json:"parent_id,omitempty"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateProduct handles POST /products
// @Summary Create a new product
// @Description Create a new product with name, description, price, SKU, category, status, images, and stock
// @Tags Products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Create Product Request"
// @Success 201 {object} map[string]interface{} "Product created successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}

	// Convert images []string to datatypes.JSON
	var imagesJSON datatypes.JSON
	if len(req.Images) > 0 {
		imagesBytes, err := json.Marshal(req.Images)
		if err != nil {
			h.logger.Warn("failed to marshal images", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid images format"})
			return
		}
		imagesJSON = datatypes.JSON(imagesBytes)
	}

	// Convert request to domain entity
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.SKU,
		CategoryID:  req.CategoryID,
		Status:      status,
		Images:      imagesJSON,
		Stock:       req.Stock,
		IsActive:    req.IsActive,
	}

	// Call service layer (business logic)
	if err := h.productService.CreateProduct(c.Request.Context(), product); err != nil {
		h.logger.Error("failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"product": product,
	})
}

// UpdateProduct handles PUT /products/:id
// @Summary Update an existing product
// @Description Update an existing product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body UpdateProductRequest true "Update Product Request"
// @Success 200 {object} map[string]interface{} "Product updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or product ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing product
	product, err := h.productService.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	if req.Images != nil {
		imagesBytes, err := json.Marshal(req.Images)
		if err == nil {
			product.Images = datatypes.JSON(imagesBytes)
		}
	}
	if req.Stock != 0 {
		product.Stock = req.Stock
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Call service layer
	if err := h.productService.UpdateProduct(c.Request.Context(), product); err != nil {
		h.logger.Error("failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated successfully",
		"product": product,
	})
}

// GetProduct handles GET /products/:id
// @Summary Get a product by ID
// @Description Get a specific product by its ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} handler.ProductResponse "Product details"
// @Failure 400 {object} map[string]string "Invalid product ID"
// @Failure 404 {object} map[string]string "Product not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetAllProducts handles GET /products (deprecated - use ListProducts instead)
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get all products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// ListProducts handles GET /products with pagination and filters
// @Summary List products with pagination and filters
// @Description Get a paginated list of products with optional filters (category_id, status, min_price, max_price, search)
// @Tags Products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param category_id query int false "Filter by category ID"
// @Param status query string false "Filter by status (ACTIVE, INACTIVE)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param search query string false "Search in name and description"
// @Success 200 {object} map[string]interface{} "List of products with pagination"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Build filters from query parameters
	filters := make(map[string]interface{})
	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			filters["category_id"] = uint(id)
		}
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["min_price"] = price
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = price
		}
	}
	if search := c.Query("search"); search != "" {
		filters["search"] = search
	}

	products, total, err := h.productService.ListProducts(c.Request.Context(), filters, page, limit)
	if err != nil {
		h.logger.Error("failed to list products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":     limit,
	})
}

// GetProductsByCategory handles GET /categories/:id/products
// @Summary Get products by category
// @Description Get a paginated list of products filtered by category ID
// @Tags Products
// @Produce json
// @Param id path int true "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{} "List of products with pagination"
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /categories/{id}/products [get]
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	products, total, err := h.productService.GetProductsByCategory(c.Request.Context(), uint(categoryID), page, limit)
	if err != nil {
		h.logger.Error("failed to get products by category", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":     limit,
	})
}

// SearchProducts handles GET /products/search
// @Summary Search products using Elasticsearch
// @Description Search products by keyword and optional category filter using Elasticsearch
// @Tags Products
// @Produce json
// @Param q query string false "Search query"
// @Param category query string false "Filter by category name"
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")

	filters := make(map[string]interface{})
	if category != "" {
		filters["category"] = category
	}

	products, err := h.productService.SearchProducts(c.Request.Context(), query, filters)
	if err != nil {
		h.logger.Error("failed to search products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}

// UpdateInventory handles PATCH /products/:id/inventory
// @Summary Update product inventory
// @Description Update product stock quantity with distributed locking
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body object true "Update Inventory Request" example({"quantity": 10})
// @Success 200 {object} map[string]string "Inventory updated successfully"
// @Failure 400 {object} map[string]string "Invalid request payload or product ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /products/{id}/inventory [patch]
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productService.UpdateInventory(c.Request.Context(), uint(id), req.Quantity); err != nil {
		h.logger.Error("failed to update inventory", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "inventory updated successfully"})
}

