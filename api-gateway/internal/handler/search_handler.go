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
func (h *SearchHandler) SearchProducts(c *gin.Context) {
	h.logger.Info("SearchHandler.SearchProducts called",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
	)
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
