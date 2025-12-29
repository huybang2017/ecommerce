package domain

import "time"

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"    // Order created, waiting for payment
	OrderStatusPaid       OrderStatus = "paid"       // Payment completed
	OrderStatusProcessing OrderStatus = "processing" // Order is being processed
	OrderStatusShipped    OrderStatus = "shipped"    // Order has been shipped
	OrderStatusDelivered  OrderStatus = "delivered"  // Order has been delivered
	OrderStatusCancelled  OrderStatus = "cancelled" // Order has been cancelled
)

// Order represents an order in the system (shop_order in db-diagram.db)
// This is the domain entity - it contains business logic and validation
// NOTE: Following db-diagram.db schema (SOURCE OF TRUTH)
type Order struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        uint        `json:"user_id" gorm:"index;not null"` // ĐỔI thành NOT NULL (bỏ guest orders)
	ShopID        uint        `json:"shop_id" gorm:"index;not null"` // THÊM MỚI - Order từ shop (theo db-diagram.db)
	ShippingAddressID *uint   `json:"shipping_address_id,omitempty" gorm:"index"` // THÊM MỚI - Reference address table
	
	// Order identification
	OrderNumber   string      `json:"order_number" gorm:"uniqueIndex;not null"` // Unique order number
	SessionID     string      `json:"session_id,omitempty" gorm:"index"` // GIỮ LẠI deprecated
	
	// Status
	Status        OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	
	// Financial breakdown (theo db-diagram.db)
	MerchandiseSubtotal float64 `json:"merchandise_subtotal" gorm:"column:merchandise_subtotal;type:decimal(15,2);not null"` // THÊM MỚI - Tổng tiền hàng
	ShippingFee         float64 `json:"shipping_fee" gorm:"type:decimal(15,2);default:0"` // Phí vận chuyển
	ShippingDiscount    float64 `json:"shipping_discount" gorm:"column:shipping_discount;type:decimal(15,2);default:0"` // THÊM MỚI - Mã freeship
	VoucherDiscount     float64 `json:"voucher_discount" gorm:"column:voucher_discount;type:decimal(15,2);default:0"` // THÊM MỚI - Mã giảm giá
	FinalAmount         float64 `json:"final_amount" gorm:"column:final_amount;type:decimal(15,2);not null"` // THÊM MỚI - Khách thực trả
	PlatformFee         float64 `json:"platform_fee" gorm:"column:platform_fee;type:decimal(15,2);default:0"` // THÊM MỚI - Phí sàn
	EarningAmount       float64 `json:"earning_amount" gorm:"column:earning_amount;type:decimal(15,2);not null"` // THÊM MỚI - Shop thực nhận
	
	// GIỮ LẠI để backward compatibility
	TotalAmount   float64     `json:"total_amount" gorm:"type:decimal(10,2);not null"` // GIỮ LẠI (sẽ sync với FinalAmount)
	Subtotal      float64     `json:"subtotal" gorm:"type:decimal(10,2);not null"` // GIỮ LẠI (sẽ sync với MerchandiseSubtotal)
	Tax           float64     `json:"tax" gorm:"type:decimal(10,2);default:0"` // GIỮ LẠI (không có trong diagram)
	Discount      float64     `json:"discount" gorm:"type:decimal(10,2);default:0"` // GIỮ LẠI
	
	// Payment & timestamps
	PaymentMethod string    `json:"payment_method" gorm:"size:50" json:"payment_method"` // THÊM MỚI
	OrderedAt     time.Time `json:"ordered_at" gorm:"column:ordered_at;index"` // THÊM MỚI
	
	// Shipping information (GIỮ LẠI)
	ShippingName       string `json:"shipping_name" gorm:"not null"`
	ShippingPhone      string `json:"shipping_phone" gorm:"not null"`
	ShippingAddress    string `json:"shipping_address" gorm:"not null"`
	ShippingCity       string `json:"shipping_city" gorm:"not null"`
	ShippingProvince   string `json:"shipping_province,omitempty"`
	ShippingPostalCode string `json:"shipping_postal_code,omitempty"`
	ShippingCountry    string `json:"shipping_country" gorm:"default:'VN'"`
	
	// Order items
	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderItem represents an item in an order (order_line in db-diagram.db)
// NOTE: Following db-diagram.db schema (SOURCE OF TRUTH)
type OrderItem struct {
	ID              uint    `json:"id" gorm:"primaryKey"`
	OrderID         uint    `json:"order_id" gorm:"index;not null"`
	ProductItemID   uint    `json:"product_item_id" gorm:"column:product_item_id;index;not null"` // THÊM MỚI - Reference product_item (SKU)
	ProductID       uint    `json:"product_id" gorm:"not null"` // GIỮ LẠI backward compatibility
	ProductName     string  `json:"product_name" gorm:"not null"` // GIỮ LẠI
	ProductSKU      string  `json:"product_sku,omitempty"` // GIỮ LẠI
	ProductImage    string  `json:"product_image,omitempty"` // GIỮ LẠI
	Quantity        int     `json:"quantity" gorm:"not null"`
	PriceAtPurchase float64 `json:"price_at_purchase" gorm:"column:price_at_purchase;type:decimal(15,2);not null"` // THÊM MỚI - Đúng tên theo diagram
	Price           float64 `json:"price" gorm:"type:decimal(10,2);not null"` // GIỮ LẠI backward compatibility (sync với PriceAtPurchase)
	Subtotal        float64 `json:"subtotal" gorm:"type:decimal(10,2);not null"` // GIỮ LẠI
	
	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

