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
// @Summary      Tạo đơn hàng từ giỏ hàng
// @Description  Tạo đơn hàng (hỗ trợ multi-shop). Hỗ trợ cả user_id (header X-User-Id) hoặc session_id cho guest
// @Tags         Order
// @Accept       json
// @Produce      json
// @Param        order  body      models.CreateOrderRequest  true  "Thông tin tạo đơn"
// @Success      201   {object}  models.OrderResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ListOrders handles GET /orders
// @Summary      Lấy danh sách đơn hàng
// @Description  Lấy danh sách đơn hàng theo user_id hoặc session_id
// @Tags         Order
// @Produce      json
// @Param        user_id    query     int     false  "User ID"
// @Param        session_id query     string  false  "Session ID"
// @Param        limit      query     int     false  "Limit"
// @Param        offset     query     int     false  "Offset"
// @Success      200   {object}  []models.OrderResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetOrder handles GET /orders/:id
// @Summary      Lấy đơn hàng theo ID
// @Tags         Order
// @Produce      json
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  models.OrderResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// GetOrderByNumber handles GET /orders/number/:order_number
// @Summary      Lấy đơn hàng theo số đơn
// @Tags         Order
// @Produce      json
// @Param        order_number  path  string  true  "Order number"
// @Success      200   {object}  models.OrderResponse
// @Failure      404   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /orders/number/{order_number} [get]
func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
