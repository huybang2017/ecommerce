package handler

import (
	"api-gateway/internal/models"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ = models.AddToCartRequest{}
var _ = models.UpdateCartRequest{}
var _ = models.SelectAllRequest{}
var _ = models.CartResponse{}
var _ = models.CartValidationResponse{}
var _ = models.ErrorResponse{}
var _ = models.SuccessResponse{}

// CartHandler xử lý các yêu cầu liên quan đến giỏ hàng thông qua Gateway
type CartHandler struct {
	gatewayService *service.GatewayService
	logger         *zap.Logger
}

func NewCartHandler(gatewayService *service.GatewayService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		gatewayService: gatewayService,
		logger:         logger,
	}
}

// GetCart handles GET /cart
func (h *CartHandler) GetCart(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ClearCart handles DELETE /cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// AddItem handles POST /cart/items
func (h *CartHandler) AddItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateItem handles PUT /cart/items/:product_item_id
func (h *CartHandler) UpdateItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// RemoveItem handles DELETE /cart/items/:product_item_id
func (h *CartHandler) RemoveItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ToggleItemSelection handles POST /cart/items/:product_item_id/toggle
func (h *CartHandler) ToggleItemSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SelectAll handles POST /cart/select_all
func (h *CartHandler) SelectAll(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SelectByShop handles POST /cart/shops/:shop_id/select
func (h *CartHandler) SelectByShop(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ClearSelected handles DELETE /cart/selected
func (h *CartHandler) ClearSelected(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetItemSelection handles PATCH /cart/items/:product_item_id/selection
func (h *CartHandler) SetItemSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetAllSelection handles PATCH /cart/selection
func (h *CartHandler) SetAllSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetShopSelection handles PATCH /cart/shops/:shop_id/selection
func (h *CartHandler) SetShopSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ValidateCart handles POST /cart/validate
func (h *CartHandler) ValidateCart(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
