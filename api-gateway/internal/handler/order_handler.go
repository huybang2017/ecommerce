package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ = models.CreateOrderRequest{}
var _ = models.OrderResponse{}
var _ = models.ErrorResponse{}
var _ = models.SuccessResponse{}

// OrderHandler proxies order-related requests to Order Service
type OrderHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

func NewOrderHandler(gatewayService *service.GatewayService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetOrderByNumber handles GET /orders/number/:order_number
func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
