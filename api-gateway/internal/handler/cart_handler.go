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
// @Summary      Lấy nội dung giỏ hàng
// @Description  Lấy danh sách sản phẩm trong giỏ hàng của người dùng hiện tại (yêu cầu đăng nhập)
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  models.CartResponse  "Thông tin giỏ hàng"
// @Failure      401  {object}  models.ErrorResponse "Chưa xác thực"
// @Failure      500  {object}  models.ErrorResponse "Lỗi hệ thống"
// @Router       /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ClearCart handles DELETE /cart
// @Summary      Xóa toàn bộ giỏ hàng
// @Description  Xóa tất cả sản phẩm trong giỏ hàng của người dùng
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Failure      401  {object}  models.ErrorResponse "Chưa xác thực"
// @Failure      500  {object}  models.ErrorResponse "Lỗi hệ thống"
// @Router       /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// AddItem handles POST /cart/items
// @Summary      Thêm sản phẩm vào giỏ
// @Description  Thêm một sản phẩm cụ thể (SKU) vào giỏ hàng với số lượng nhất định
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        item  body      models.AddToCartRequest  true  "Thông tin sản phẩm thêm vào"
// @Success      201   {object}  models.CartResponse
// @Failure      400   {object}  models.ErrorResponse "Dữ liệu không hợp lệ hoặc thiếu hàng"
// @Router       /cart/items [post]
func (h *CartHandler) AddItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// UpdateItem handles PUT /cart/items/:product_item_id
// @Summary      Cập nhật số lượng sản phẩm
// @Description  Thay đổi số lượng của một SKU cụ thể trong giỏ hàng
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product_item_id  path      int  true  "ID của SKU sản phẩm"
// @Param        quantity         body      models.UpdateCartRequest  true  "Số lượng mới"
// @Success      200   {object}  models.CartResponse
// @Router       /cart/items/{product_item_id} [put]
func (h *CartHandler) UpdateItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// RemoveItem handles DELETE /cart/items/:product_item_id
// @Summary      Xóa sản phẩm khỏi giỏ hàng
// @Description  Xóa một sản phẩm cụ thể ra khỏi giỏ hàng
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Param        product_item_id  path      int  true  "ID của SKU sản phẩm"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse "Tham số không hợp lệ"
// @Failure      401   {object}  models.ErrorResponse "Chưa xác thực"
// @Failure      404   {object}  models.ErrorResponse "Không tìm thấy sản phẩm"
// @Failure      500   {object}  models.ErrorResponse "Lỗi hệ thống"
// @Router       /cart/items/{product_item_id} [delete]
func (h *CartHandler) RemoveItem(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ToggleItemSelection handles POST /cart/items/:product_item_id/toggle
// @Summary      Bật/Tắt chọn sản phẩm
// @Description  Thay đổi trạng thái được chọn (tick) của một sản phẩm để chuẩn bị thanh toán
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product_item_id  path      int  true  "ID của SKU sản phẩm"
// @Success      200   {object}  models.CartResponse
// @Router       /cart/items/{product_item_id}/toggle [post]
func (h *CartHandler) ToggleItemSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SelectAll handles POST /cart/select_all
// @Summary      Chọn/Bỏ chọn tất cả
// @Description  Chọn hoặc bỏ chọn toàn bộ sản phẩm có trong giỏ hàng
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.SelectAllRequest  true  "Trạng thái chọn (true/false)"
// @Success      200      {object}  models.SuccessResponse
// @Router       /cart/select_all [post]
func (h *CartHandler) SelectAll(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SelectByShop handles POST /cart/shops/:shop_id/select
// @Summary      Chọn sản phẩm theo Shop
// @Description  Chọn hoặc bỏ chọn tất cả sản phẩm thuộc về một cửa hàng cụ thể
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shop_id  path      int  true  "ID của cửa hàng"
// @Param        request  body      models.SelectAllRequest  true  "Trạng thái chọn"
// @Success      200      {object}  models.SuccessResponse
// @Router       /cart/shops/{shop_id}/select [post]
func (h *CartHandler) SelectByShop(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ClearSelected handles DELETE /cart/selected
// @Summary      Xóa các mục đã chọn
// @Description  Xóa tất cả sản phẩm đang được tick chọn ra khỏi giỏ hàng
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SuccessResponse
// @Router       /cart/selected [delete]
func (h *CartHandler) ClearSelected(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetItemSelection handles PATCH /cart/items/:product_item_id/selection
// @Summary      Set item selection
// @Description  Set the selection state (selected/deselected) for a cart item
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product_item_id  path      int  true  "ID của SKU sản phẩm"
// @Param        request          body      models.SelectionRequest  true  "Trạng thái chọn"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse "Dữ liệu không hợp lệ"
// @Failure      401   {object}  models.ErrorResponse "Chưa xác thực"
// @Failure      404   {object}  models.ErrorResponse "Không tìm thấy sản phẩm"
// @Failure      500   {object}  models.ErrorResponse "Lỗi hệ thống"
// @Router       /cart/items/{product_item_id}/selection [patch]
func (h *CartHandler) SetItemSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetAllSelection handles PATCH /cart/selection
// @Summary      Set selection for all items
// @Description  Set selection state for all items in the cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.SelectAllRequest  true  "Trạng thái chọn"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse "Dữ liệu không hợp lệ"
// @Failure      401   {object}  models.ErrorResponse "Chưa xác thực"
// @Failure      500   {object}  models.ErrorResponse "Lỗi hệ thống"
// @Router       /cart/selection [patch]
func (h *CartHandler) SetAllSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// SetShopSelection handles PATCH /cart/shops/:shop_id/selection
// @Summary      Set selection for shop items
// @Description  Set selection state for all items from a shop
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        shop_id  path      int  true  "ID của cửa hàng"
// @Param        request  body      models.SelectShopRequest  true  "Trạng thái chọn"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse "Dữ liệu không hợp lệ"
// @Failure      401   {object}  models.ErrorResponse "Chưa xác thực"
// @Failure      500   {object}  models.ErrorResponse "Lỗi hệ thống"
// @Router       /cart/shops/{shop_id}/selection [patch]
func (h *CartHandler) SetShopSelection(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}

// ValidateCart handles POST /cart/validate
// @Summary      Kiểm tra giỏ hàng trước thanh toán
// @Description  Kiểm tra tồn kho và trạng thái của các sản phẩm đã chọn trước khi chuyển sang bước thanh toán
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.CartValidationResponse "Kết quả kiểm tra"
// @Failure      400  {object}  models.ErrorResponse          "Sản phẩm không hợp lệ hoặc hết hàng"
// @Router       /cart/validate [post]
func (h *CartHandler) ValidateCart(c *gin.Context) {
	gatewayHandler := NewGatewayHandler(h.gatewayService, h.logger)
	gatewayHandler.ProxyRequest(c)
}
