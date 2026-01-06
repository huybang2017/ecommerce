package domain

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // Order created, waiting for payment
	OrderStatusPaid       OrderStatus = "paid"       // Payment completed
	OrderStatusProcessing OrderStatus = "processing" // Order is being processed
	OrderStatusShipped    OrderStatus = "shipped"    // Order has been shipped
	OrderStatusDelivered  OrderStatus = "delivered"  // Order has been delivered
	OrderStatusCancelled  OrderStatus = "cancelled"  // Order has been cancelled
)

// Order represents an order in the system (shop_order in db-diagram.db)
// This is the domain entity - it contains business logic and validation
// NOTE: Following db-diagram.db schema (SOURCE OF TRUTH)
type Order struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// Business identifiers
	OrderNumber string `json:"order_number" gorm:"size:50;uniqueIndex;not null"`

	// Ownership
	UserID uint `json:"user_id" gorm:"index;not null"`
	ShopID uint `json:"shop_id" gorm:"index;not null"`

	// Shipping
	ShippingAddressID uint `json:"shipping_address_id" gorm:"index;not null"`

	// Status
	Status OrderStatus `json:"status" gorm:"type:varchar(20);not null"`

	// Financial snapshot (SOURCE OF TRUTH)
	MerchandiseSubtotal float64 `json:"merchandise_subtotal" gorm:"type:decimal(15,2);not null"`
	ShippingFee         float64 `json:"shipping_fee" gorm:"type:decimal(15,2);not null"`
	ShippingDiscount    float64 `json:"shipping_discount" gorm:"type:decimal(15,2);not null"`
	VoucherDiscount     float64 `json:"voucher_discount" gorm:"type:decimal(15,2);not null"`

	FinalAmount   float64 `json:"final_amount" gorm:"type:decimal(15,2);not null"`
	PlatformFee   float64 `json:"platform_fee" gorm:"type:decimal(15,2);not null"`
	EarningAmount float64 `json:"earning_amount" gorm:"type:decimal(15,2);not null"`

	// Payment
	PaymentMethod string `json:"payment_method" gorm:"size:50;not null"`

	// Time
	OrderedAt time.Time `json:"ordered_at" gorm:"index;not null"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

// OrderItem represents an item in an order (order_line in db-diagram.db)
// NOTE: Following db-diagram.db schema (SOURCE OF TRUTH)
type OrderItem struct {
	ID uint `json:"id" gorm:"primaryKey"`

	OrderID       uint `json:"order_id" gorm:"index;not null"`
	ProductItemID uint `json:"product_item_id" gorm:"index;not null"`

	Quantity        int     `json:"quantity" gorm:"not null"`
	PriceAtPurchase float64 `json:"price_at_purchase" gorm:"type:decimal(15,2);not null"`

	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for Order
// NOTE: Đổi từ "orders" sang "shop_order" theo db-diagram.db
func (Order) TableName() string {
	return "shop_order"
}

// TableName specifies the table name for OrderItem
// NOTE: Đổi từ "order_items" sang "order_line" theo db-diagram.db
func (OrderItem) TableName() string {
	return "order_line"
}
