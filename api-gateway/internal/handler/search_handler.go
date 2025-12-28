package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Import models for Swagger documentation generation
var _ = models.SearchResponse{}
var _ = models.ErrorResponse{}

// SearchHandler handles search requests
type SearchHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(gatewayService *service.GatewayService, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// SearchProducts handles GET /api/v1/search
// @Summary Search products
// @Description Search products by keyword with filters (category, price range, status) and sort options. Uses Elasticsearch for full-text search. This service indexes products from Kafka events and provides fast search capabilities.
// @Tags Search
// @Produce json
// @Param q query string false "Search keyword (searches in product name and description using Elasticsearch full-text search)"
// @Param category_id query int false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param status query string false "Filter by status (ACTIVE, INACTIVE)"
// @Param sort_field query string false "Sort field (price, name, created_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} models.SearchResponse "Search results with products and pagination"
// @Failure 400 {object} models.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /search [get]
func (h *SearchHandler) SearchProducts(c *gin.Context) {
	h.logger.Info("SearchHandler.SearchProducts called",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

