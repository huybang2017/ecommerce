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
	SessionID string `json:"session_id,omitempty"` // Deprecated
	
	// Shipping information
	ShippingName       string `json:"shipping_name" binding:"required"`
	ShippingPhone      string `json:"shipping_phone" binding:"required"`
	ShippingAddress    string `json:"shipping_address" binding:"required"`
	ShippingCity       string `json:"shipping_city" binding:"required"`
	ShippingProvince   string `json:"shipping_province,omitempty"`
	ShippingPostalCode string `json:"shipping_postal_code,omitempty"`
	ShippingCountry    string `json:"shipping_country,omitempty"`
	ShippingAddressID  *uint  `json:"shipping_address_id,omitempty"` // THÊM MỚI - Reference address table
	
	// Financial (theo db-diagram.db)
	ShippingFee      float64 `json:"shipping_fee,omitempty"`
	ShippingDiscount float64 `json:"shipping_discount,omitempty"` // Mã freeship
	VoucherDiscount  float64 `json:"voucher_discount,omitempty"`  // Mã giảm giá
	PaymentMethod    string  `json:"payment_method,omitempty"`
}

// CreateOrderResponse represents the response after creating orders
// MARKETPLACE: Can return multiple shop_orders
type CreateOrderResponse struct {
	Orders       []*domain.Order `json:"orders"`        // Multiple shop_orders (1 per shop)
	OrderNumbers []string         `json:"order_numbers"` // Order numbers for each shop_order
}

// CreateOrder creates orders from the cart with MARKETPLACE logic
// Business logic:
// 1. Get cart from Redis
// 2. Group cart items by shop_id
// 3. For each shop, create a shop_order with financial breakdown
// 4. Clear cart
// 5. Publish OrderCreated event for each shop_order
// Returns CreateOrderResponse with multiple shop_orders
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
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
	
	userID := uint(0)
	if req.UserID != nil {
		userID = *req.UserID
	}
	
	// 2. MARKETPLACE: Group cart items by shop_id
	itemsByShop := make(map[uint][]*domain.CartItem)
	for productID, cartItem := range cart.Items {
		if cartItem == nil {
			continue
		}
		
		shopID := cartItem.ShopID
		if shopID == 0 {
			// Backward compatibility: if shop_id not set, use default
			s.logger.Warn("cart item missing shop_id, using default", zap.Uint("product_id", productID))
			shopID = 1
		}
		
		if itemsByShop[shopID] == nil {
			itemsByShop[shopID] = make([]*domain.CartItem, 0)
		}
		itemsByShop[shopID] = append(itemsByShop[shopID], cartItem)
	}

	if len(itemsByShop) == 0 {
		return nil, errors.New("no items found in cart")
	}

	// 3. Create shop_order for each shop
	createdOrders := make([]*domain.Order, 0, len(itemsByShop))
	orderNumbers := make([]string, 0, len(itemsByShop))

	for shopID, shopItems := range itemsByShop {
		// Calculate financial breakdown for this shop
		merchandiseSubtotal := float64(0)
		for _, item := range shopItems {
			merchandiseSubtotal += item.Price * float64(item.Quantity)
		}

		// Shipping fee (can be per shop or shared - for now, divide equally)
		shippingFee := req.ShippingFee / float64(len(itemsByShop))
		if shippingFee < 0 {
			shippingFee = 0
		}

		// Discounts (can be per shop or shared - for now, divide equally)
		shippingDiscount := req.ShippingDiscount / float64(len(itemsByShop))
		if shippingDiscount < 0 {
			shippingDiscount = 0
		}
		voucherDiscount := req.VoucherDiscount / float64(len(itemsByShop))
		if voucherDiscount < 0 {
			voucherDiscount = 0
		}

		// Final amount for this shop
		finalAmount := merchandiseSubtotal + shippingFee - shippingDiscount - voucherDiscount
		if finalAmount < 0 {
			finalAmount = 0
		}

		// Platform fee: 5% of merchandise_subtotal (per shop)
		platformFee := merchandiseSubtotal * 0.05

		// Earning amount: Shop receives final_amount - platform_fee
		earningAmount := finalAmount - platformFee
		if earningAmount < 0 {
			earningAmount = 0
		}

		// Generate order number for this shop_order
		orderNumber := s.generateOrderNumber()

		// Create shop_order
		order := &domain.Order{
			UserID:    userID,
			ShopID:    shopID, // Shop ID
			SessionID: req.SessionID, // Deprecated
			OrderNumber: orderNumber,
			Status:     domain.OrderStatusPending,
			
			// Financial breakdown (per shop)
			MerchandiseSubtotal: merchandiseSubtotal,
			ShippingFee:         shippingFee,
			ShippingDiscount:    shippingDiscount,
			VoucherDiscount:     voucherDiscount,
			FinalAmount:         finalAmount,
			PlatformFee:         platformFee,
			EarningAmount:       earningAmount,
			
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
			ShippingAddressID: req.ShippingAddressID, // Reference address table
			
			Items: make([]domain.OrderItem, 0, len(shopItems)),
		}

		// Set default country if not provided
		if order.ShippingCountry == "" {
			order.ShippingCountry = "VN"
		}

		// Convert shop items to order items
		for _, cartItem := range shopItems {
			productItemID := cartItem.ProductItemID
			if productItemID == 0 {
				// Backward compatibility: use product_id as product_item_id
				productItemID = cartItem.ProductID
			}

			orderItem := domain.OrderItem{
				ProductID:       cartItem.ProductID,
				ProductItemID:   productItemID, // SKU ID
				ProductName:     cartItem.Name,
				ProductSKU:      cartItem.SKU,
				ProductImage:    cartItem.Image,
				Price:           cartItem.Price,
				PriceAtPurchase: cartItem.Price, // Price at time of order
				Quantity:        cartItem.Quantity,
				Subtotal:        cartItem.Price * float64(cartItem.Quantity),
			}
			order.Items = append(order.Items, orderItem)
		}

		// Save shop_order to database
		if err := s.orderRepo.Create(order); err != nil {
			s.logger.Error("failed to create shop_order", zap.Uint("shop_id", shopID), zap.Error(err))
			// Continue with other shops even if one fails
			continue
		}

		createdOrders = append(createdOrders, order)
		orderNumbers = append(orderNumbers, orderNumber)

		s.logger.Info("shop_order created",
			zap.Uint("order_id", order.ID),
			zap.Uint("shop_id", shopID),
			zap.String("order_number", orderNumber),
			zap.Float64("final_amount", order.FinalAmount),
			zap.Float64("platform_fee", order.PlatformFee),
			zap.Float64("earning_amount", order.EarningAmount),
		)

		// Publish OrderCreated event for this shop_order (async)
		go func(shopOrder *domain.Order, sID uint) {
			event := &domain.OrderEvent{
				EventType: "order_created",
				OrderID:   shopOrder.ID,
				OrderData: shopOrder,
				Timestamp: time.Now(),
			}
			
			if err := s.eventPublisher.PublishOrderEvent(event); err != nil {
				s.logger.Error("failed to publish order_created event",
					zap.Uint("order_id", shopOrder.ID),
					zap.Uint("shop_id", sID),
					zap.Error(err),
				)
			} else {
				s.logger.Info("order_created event published",
					zap.Uint("order_id", shopOrder.ID),
					zap.Uint("shop_id", sID),
				)
			}
		}(order, shopID)
	}

	if len(createdOrders) == 0 {
		return nil, errors.New("failed to create any orders")
	}

	// 4. Clear cart (async)
	go func() {
		userIDStr := ""
		if req.UserID != nil {
			userIDStr = fmt.Sprintf("%d", *req.UserID)
		}
		_ = s.cartRepo.DeleteCart(userIDStr)
	}()

	return &CreateOrderResponse{
		Orders:       createdOrders,
		OrderNumbers: orderNumbers,
	}, nil
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

