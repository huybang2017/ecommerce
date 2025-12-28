package handler

import (
	"net/http"
	"search-service/internal/domain"
	"search-service/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SearchHandler handles HTTP requests for search operations
// This is the transport layer - it knows HOW to handle HTTP (Gin framework)
// It delegates business logic to the service layer
type SearchHandler struct {
	searchService *service.SearchService
	logger        *zap.Logger
}

// NewSearchHandler creates a new search handler
// Dependency injection: we inject the service
func NewSearchHandler(searchService *service.SearchService, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
		logger:        logger,
	}
}

// SearchProducts handles GET /search
// @Summary Search products
// @Description Search products by keyword with filters (category, price range) and sort options
// @Tags Search
// @Produce json
// @Param q query string false "Search keyword"
// @Param category_id query int false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param status query string false "Filter by status (ACTIVE, INACTIVE)"
// @Param sort_field query string false "Sort field (price, name, created_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} domain.SearchResult "Search results"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /search [get]
func (h *SearchHandler) SearchProducts(c *gin.Context) {
	// Parse query parameters
	query := c.Query("q")

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Parse filters
	var filters *domain.SearchFilters
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryIDUint := uint(categoryID)
			if filters == nil {
				filters = &domain.SearchFilters{}
			}
			filters.CategoryID = &categoryIDUint
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			if filters == nil {
				filters = &domain.SearchFilters{}
			}
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			if filters == nil {
				filters = &domain.SearchFilters{}
			}
			filters.MaxPrice = &maxPrice
		}
	}

	if status := c.Query("status"); status != "" {
		if filters == nil {
			filters = &domain.SearchFilters{}
		}
		filters.Status = &status
	}

	// Parse sort
	var sort *domain.SearchSort
	if sortField := c.Query("sort_field"); sortField != "" {
		sort = &domain.SearchSort{
			Field: sortField,
			Order: c.DefaultQuery("sort_order", "asc"),
		}
	}

	// Build search request
	searchReq := &domain.SearchRequest{
		Query:   query,
		Filters: filters,
		Sort:    sort,
		Page:    page,
		Limit:   limit,
	}

	// Call service layer
	result, err := h.searchService.SearchProducts(c.Request.Context(), searchReq)
	if err != nil {
		h.logger.Error("failed to search products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HealthCheck handles GET /health
func (h *SearchHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "search-service"})
}


