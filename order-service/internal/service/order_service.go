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
	productClient  OrderProductServiceClient
	eventPublisher domain.OrderEventPublisher
	logger         *zap.Logger
}

// OrderProductServiceClient defines interface to communicate with Product Service
// NOTE: OrderService needs FULL product data for validation (Stock, IsActive)
type OrderProductServiceClient interface {
	// GetProductItem fetches single product item details (SKU-level)
	GetProductItem(productItemID uint) (*OrderProductItemDTO, error)

	// GetProductItems fetches multiple product items in batch (for performance)
	GetProductItems(productItemIDs []uint) (map[uint]*OrderProductItemDTO, error)
}

// OrderProductItemDTO represents FULL product item data from Product Service
// This includes validation fields (Stock, IsActive) required for order creation
type OrderProductItemDTO struct {
	ID          uint    `json:"id"`           // ProductItem ID (SKU)
	ProductID   uint    `json:"product_id"`   // Base product ID
	ShopID      uint    `json:"shop_id"`      // Shop that owns this product
	ProductName string  `json:"product_name"` // Product name
	SKU         string  `json:"sku"`          // SKU code
	Price       float64 `json:"price"`        // Current price
	Stock       int     `json:"stock"`        // Available stock (REQUIRED for validation)
	ImageURL    string  `json:"image_url"`    // Product image
	IsActive    bool    `json:"is_active"`    // Product active status (REQUIRED for validation)
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo *postgres.OrderRepository,
	cartRepo domain.CartRepository,
	productClient OrderProductServiceClient,
	eventPublisher domain.OrderEventPublisher,
	logger *zap.Logger,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		cartRepo:       cartRepo,
		productClient:  productClient,
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
	OrderNumbers []string        `json:"order_numbers"` // Order numbers for each shop_order
}

// CreateOrder creates orders from the cart with MARKETPLACE logic (REFACTORED - SENIOR LEVEL)
// Business logic (CORRECT FLOW):
// 1. Load cart from Redis
// 2. Filter SELECTED items only
// 3. Load SKU snapshots from Product Service & validate (price, stock, active status)
// 4. Group by shop_id
// 5. For each shop: calculate financials using server-side rules & snapshot prices
// 6. Create shop_orders in DB
// 7. Publish events (SYNC for MVP, TODO: outbox pattern)
// 8. Clear cart (SYNC)
// Returns CreateOrderResponse with multiple shop_orders
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// Validate required fields
	if req.UserID == nil {
		return nil, errors.New("user_id is required")
	}

	if req.ShippingAddressID == nil {
		return nil, errors.New("shipping_address_id is required")
	}

	userID := *req.UserID
	userIDStr := fmt.Sprintf("%d", userID)

	// STEP 1: Load cart from Redis
	cart, err := s.cartRepo.GetCart(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if cart == nil || cart.IsEmpty() {
		return nil, domain.ErrCartEmpty
	}

	// STEP 2: Filter SELECTED items only (B5 fix)
	selectedItems := cart.GetSelectedItems()
	if len(selectedItems) == 0 {
		return nil, domain.ErrNoItemsSelected
	}

	// STEP 3: Load SKU snapshots from Product Service & validate (B1 + B2 fix)
	productItemIDs := make([]uint, 0, len(selectedItems))
	for _, item := range selectedItems {
		productItemIDs = append(productItemIDs, item.ProductItemID)
	}

	// Batch load product items
	productItems, err := s.productClient.GetProductItems(productItemIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to load product items: %w", err)
	}

	// Validate each selected item
	for _, item := range selectedItems {
		sku, exists := productItems[item.ProductItemID]
		if !exists {
			return nil, fmt.Errorf("product item %d not found", item.ProductItemID)
		}

		// Validate SKU status
		if !sku.IsActive {
			return nil, fmt.Errorf("product %s is not available", sku.ProductName)
		}

		// Validate stock
		if sku.Stock <= 0 {
			return nil, fmt.Errorf("product %s is out of stock", sku.ProductName)
		}

		if item.Quantity > sku.Stock {
			return nil, fmt.Errorf("insufficient stock for %s (requested: %d, available: %d)",
				sku.ProductName, item.Quantity, sku.Stock)
		}
	}

	// STEP 4: Group selected items by shop_id
	itemsByShop := make(map[uint][]*domain.CartItem)
	for _, item := range selectedItems {
		sku := productItems[item.ProductItemID]
		shopID := sku.ShopID

		if shopID == 0 {
			s.logger.Warn("SKU missing shop_id, using default",
				zap.Uint("product_item_id", item.ProductItemID))
			shopID = 1
		}

		itemsByShop[shopID] = append(itemsByShop[shopID], item)
	}

	if len(itemsByShop) == 0 {
		return nil, errors.New("no valid items to checkout")
	}

	// STEP 5: Create shop_order for each shop
	createdOrders := make([]*domain.Order, 0, len(itemsByShop))
	orderNumbers := make([]string, 0, len(itemsByShop))

	for shopID, shopItems := range itemsByShop {
		// Calculate merchandise subtotal using SKU snapshot prices (B1 fix - server-side pricing)
		merchandiseSubtotal := float64(0)
		for _, item := range shopItems {
			sku := productItems[item.ProductItemID]
			// Use price from Product Service, NOT from cart
			lineTotal := sku.Price * float64(item.Quantity)
			merchandiseSubtotal += lineTotal
		}

		// Calculate shipping & discounts (B3/B4 fix - server-side rules, MVP: simple flat rate)
		// TODO: Call ShippingService for accurate per-shop shipping fee
		// TODO: Call PromotionService for voucher validation & discount calculation
		shippingFee := 30000.0  // MVP: flat 30k VND per shop
		shippingDiscount := 0.0 // MVP: no freeship
		voucherDiscount := 0.0  // MVP: no voucher

		// Final amount
		finalAmount := merchandiseSubtotal + shippingFee - shippingDiscount - voucherDiscount
		if finalAmount < 0 {
			finalAmount = 0
		}

		// Platform fee: 5% of merchandise
		platformFee := merchandiseSubtotal * 0.05

		// Shop earning
		earningAmount := finalAmount - platformFee
		if earningAmount < 0 {
			earningAmount = 0
		}

		// Generate order number
		orderNumber := s.generateOrderNumber()

		// Create Order aggregate
		order := &domain.Order{
			OrderNumber:       orderNumber,
			UserID:            userID,
			ShopID:            shopID,
			ShippingAddressID: *req.ShippingAddressID,
			Status:            domain.OrderStatusPending,

			// Financial snapshot
			MerchandiseSubtotal: merchandiseSubtotal,
			ShippingFee:         shippingFee,
			ShippingDiscount:    shippingDiscount,
			VoucherDiscount:     voucherDiscount,
			FinalAmount:         finalAmount,
			PlatformFee:         platformFee,
			EarningAmount:       earningAmount,

			PaymentMethod: req.PaymentMethod,
			OrderedAt:     time.Now(),

			Items: make([]domain.OrderItem, 0, len(shopItems)),
		}

		// Set default payment method if not provided
		if order.PaymentMethod == "" {
			order.PaymentMethod = "COD"
		}

		// Create OrderItems with snapshot price
		for _, item := range shopItems {
			sku := productItems[item.ProductItemID]

			orderItem := domain.OrderItem{
				ProductItemID:   item.ProductItemID,
				Quantity:        item.Quantity,
				PriceAtPurchase: sku.Price, // Snapshot price from Product Service
			}
			order.Items = append(order.Items, orderItem)
		}

		// STEP 6: Save shop_order to database
		if err := s.orderRepo.Create(order); err != nil {
			s.logger.Error("failed to create shop_order",
				zap.Uint("shop_id", shopID),
				zap.Error(err))
			// For MVP: fail fast if any shop order fails
			// TODO: Consider partial success handling
			return nil, fmt.Errorf("failed to create order for shop %d: %w", shopID, err)
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
	}

	if len(createdOrders) == 0 {
		return nil, errors.New("failed to create any orders")
	}

	// STEP 7: Publish OrderCreated events (B7 fix - SYNC for MVP, no goroutine)
	// TODO: Implement outbox pattern for reliable event delivery
	for _, order := range createdOrders {
		event := &domain.OrderEvent{
			EventType: "order_created",
			OrderID:   order.ID,
			OrderData: order,
			Timestamp: time.Now(),
		}

		if err := s.eventPublisher.PublishOrderEvent(event); err != nil {
			s.logger.Error("failed to publish order_created event",
				zap.Uint("order_id", order.ID),
				zap.Uint("shop_id", order.ShopID),
				zap.Error(err),
			)
			// For MVP: log error but don't fail the order
			// TODO: Retry with outbox pattern
		} else {
			s.logger.Info("order_created event published",
				zap.Uint("order_id", order.ID),
				zap.Uint("shop_id", order.ShopID),
			)
		}
	}

	// STEP 8: Clear cart (B7 fix - SYNC, handle error)
	if err := s.cartRepo.DeleteCart(userIDStr); err != nil {
		s.logger.Warn("failed to clear cart after order creation",
			zap.String("user_id", userIDStr),
			zap.Error(err),
		)
		// Don't fail order creation if cart clear fails
	}

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
