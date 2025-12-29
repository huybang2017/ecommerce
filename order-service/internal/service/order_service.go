package service

import (
	"errors"
	"fmt"
	"order-service/internal/domain"
	"order-service/internal/repository/postgres"
	"time"

	"go.uber.org/zap"
)

// OrderService handles business logic for orders
// This is the business logic layer - it contains domain rules and orchestrates operations
type OrderService struct {
	orderRepo      *postgres.OrderRepository
	cartRepo       domain.CartRepository
	eventPublisher domain.OrderEventPublisher
	logger         *zap.Logger
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo *postgres.OrderRepository,
	cartRepo domain.CartRepository,
	eventPublisher domain.OrderEventPublisher,
	logger *zap.Logger,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		cartRepo:       cartRepo,
		eventPublisher: eventPublisher,
		logger:         logger,
	}
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	UserID    *uint  `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"` // GIỮ LẠI deprecated
	
	// Shipping information
	ShippingName       string `json:"shipping_name" binding:"required"`
	ShippingPhone      string `json:"shipping_phone" binding:"required"`
	ShippingAddress    string `json:"shipping_address" binding:"required"`
	ShippingCity       string `json:"shipping_city" binding:"required"`
	ShippingProvince   string `json:"shipping_province,omitempty"`
	ShippingPostalCode string `json:"shipping_postal_code,omitempty"`
	ShippingCountry    string `json:"shipping_country,omitempty"`
	
	// Financial (theo db-diagram.db)
	ShippingFee      float64 `json:"shipping_fee,omitempty"`
	ShippingDiscount float64 `json:"shipping_discount,omitempty"` // THÊM MỚI - Mã freeship
	VoucherDiscount  float64 `json:"voucher_discount,omitempty"`  // THÊM MỚI - Mã giảm giá
	PaymentMethod    string  `json:"payment_method,omitempty"`     // THÊM MỚI
	
	// Backward compatibility
	Tax      float64 `json:"tax,omitempty"`      // GIỮ LẠI
	Discount float64 `json:"discount,omitempty"` // GIỮ LẠI
}

// CreateOrder creates an order from the cart
// Business logic:
// 1. Get cart from Redis
// 2. Validate cart is not empty
// 3. Create order with items
// 4. Clear cart
// 5. Publish OrderCreated event
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*domain.Order, error) {
	// 1. Get cart (chỉ dùng userID - đã bỏ sessionID)
	var cart *domain.Cart
	var err error
	
	userIDStr := ""
	if req.UserID != nil {
		userIDStr = fmt.Sprintf("%d", *req.UserID)
	}
	
	cart, err = s.cartRepo.GetCart(userIDStr) // Đã sửa: chỉ userID
	
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}
	
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}
	
	// 2. Calculate financial breakdown (theo db-diagram.db)
	// merchandise_subtotal: Tổng tiền hàng
	merchandiseSubtotal := cart.Total
	
	// shipping_fee: Phí vận chuyển
	shippingFee := req.ShippingFee
	if shippingFee < 0 {
		shippingFee = 0
	}
	
	// shipping_discount: Mã freeship (mặc định 0)
	shippingDiscount := float64(0)
	if req.ShippingDiscount > 0 {
		shippingDiscount = req.ShippingDiscount
	}
	
	// voucher_discount: Mã giảm giá (mặc định 0)
	voucherDiscount := float64(0)
	if req.VoucherDiscount > 0 {
		voucherDiscount = req.VoucherDiscount
	}
	
	// final_amount: Khách thực trả = merchandise_subtotal + shipping_fee - shipping_discount - voucher_discount
	finalAmount := merchandiseSubtotal + shippingFee - shippingDiscount - voucherDiscount
	if finalAmount < 0 {
		finalAmount = 0
	}
	
	// platform_fee: Phí sàn (5%)
	platformFee := finalAmount * 0.05
	
	// earning_amount: Shop thực nhận = final_amount - platform_fee
	earningAmount := finalAmount - platformFee
	
	// Backward compatibility fields
	subtotal := merchandiseSubtotal
	tax := req.Tax
	if tax < 0 {
		tax = 0
	}
	discount := req.Discount
	if discount < 0 {
		discount = 0
	}
	totalAmount := finalAmount // Sync với FinalAmount
	
	// 3. Generate order number
	orderNumber := s.generateOrderNumber()
	
	// 4. Create order (với financial breakdown theo db-diagram.db)
	userID := uint(0)
	if req.UserID != nil {
		userID = *req.UserID
	}
	
	// TODO: Get shop_id từ cart items (tất cả products trong cart phải cùng shop)
	// Tạm thời hardcode shop_id = 1
	shopID := uint(1)
	
	order := &domain.Order{
		UserID:    userID,
		ShopID:    shopID, // THÊM MỚI
		SessionID: req.SessionID, // GIỮ LẠI deprecated
		OrderNumber: orderNumber,
		Status:     domain.OrderStatusPending,
		
		// Financial breakdown (theo db-diagram.db)
		MerchandiseSubtotal: merchandiseSubtotal,
		ShippingFee:         shippingFee,
		ShippingDiscount:    shippingDiscount,
		VoucherDiscount:     voucherDiscount,
		FinalAmount:         finalAmount,
		PlatformFee:         platformFee,
		EarningAmount:       earningAmount,
		
		// Backward compatibility
		TotalAmount: totalAmount,
		Subtotal:    subtotal,
		Tax:         tax,
		Discount:    discount,
		
		// Payment & timestamps
		PaymentMethod: req.PaymentMethod,
		OrderedAt:     time.Now(),
		
		// Shipping info
		ShippingName:     req.ShippingName,
		ShippingPhone:    req.ShippingPhone,
		ShippingAddress:  req.ShippingAddress,
		ShippingCity:     req.ShippingCity,
		ShippingProvince: req.ShippingProvince,
		ShippingPostalCode: req.ShippingPostalCode,
		ShippingCountry:   req.ShippingCountry,
		Items: make([]domain.OrderItem, 0, len(cart.Items)),
	}
	
	// Set default country if not provided
	if order.ShippingCountry == "" {
		order.ShippingCountry = "VN"
	}
	
	// 5. Convert cart items to order items
	// cart.Items is a map[uint]*CartItem, so we iterate over values
	for productID, cartItem := range cart.Items {
		if cartItem == nil {
			continue
		}
		
		// TODO: Get product_item_id từ product_id (sau khi có SKU system)
		// Tạm thời hardcode product_item_id = product_id
		productItemID := productID
		
		orderItem := domain.OrderItem{
			ProductID:   productID, // Use the key from map
			ProductItemID: productItemID, // THÊM MỚI
			ProductName: cartItem.Name,
			ProductSKU:  cartItem.SKU,
			ProductImage: cartItem.Image,
			Price:       cartItem.Price,
			Quantity:    cartItem.Quantity,
			Subtotal:    cartItem.Price * float64(cartItem.Quantity),
		}
		order.Items = append(order.Items, orderItem)
	}
	
	// 6. Save order to database
	if err := s.orderRepo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	s.logger.Info("order created",
		zap.Uint("order_id", order.ID),
		zap.String("order_number", order.OrderNumber),
		zap.Float64("total_amount", order.TotalAmount),
	)
	
	// 7. Clear cart (async - don't block on cart clearing)
	go func() {
		userIDStr := ""
		if req.UserID != nil {
			userIDStr = fmt.Sprintf("%d", *req.UserID)
		}
		_ = s.cartRepo.DeleteCart(userIDStr) // Đã sửa: chỉ userID
	}()
	
	// 8. Publish OrderCreated event (async - event-driven communication)
	go func() {
		event := &domain.OrderEvent{
			EventType: "order_created",
			OrderID:   order.ID,
			OrderData: order,
			Timestamp: time.Now(),
		}
		
		if err := s.eventPublisher.PublishOrderEvent(event); err != nil {
			s.logger.Error("failed to publish order_created event",
				zap.Uint("order_id", order.ID),
				zap.Error(err),
			)
		} else {
			s.logger.Info("order_created event published",
				zap.Uint("order_id", order.ID),
			)
		}
	}()
	
	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(orderID uint) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// GetOrderByOrderNumber retrieves an order by order number
func (s *OrderService) GetOrderByOrderNumber(orderNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByOrderNumber(orderNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// ListOrders retrieves orders for a user or session
func (s *OrderService) ListOrders(userID *uint, sessionID string, limit, offset int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64
	var err error
	
	if userID != nil {
		orders, total, err = s.orderRepo.GetByUserID(*userID, limit, offset)
	} else if sessionID != "" {
		orders, total, err = s.orderRepo.GetBySessionID(sessionID, limit, offset)
	} else {
		return nil, 0, errors.New("user_id or session_id is required")
	}
	
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	
	return orders, total, nil
}

// generateOrderNumber generates a unique order number
// Format: ORD-YYYYMMDD-HHMMSS-XXXX (where XXXX is a random 4-digit number)
func (s *OrderService) generateOrderNumber() string {
	now := time.Now()
	timestamp := now.Format("20060102-150405")
	// Simple random number (in production, use crypto/rand)
	random := now.Nanosecond() % 10000
	return fmt.Sprintf("ORD-%s-%04d", timestamp, random)
}

