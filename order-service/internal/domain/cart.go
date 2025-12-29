package domain

// CartItem represents a single item in the shopping cart
type CartItem struct {
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Image     string  `json:"image,omitempty"`
	SKU       string  `json:"sku,omitempty"`
}

// Cart represents a shopping cart
// Cart is stored in Redis with key: "cart:user:{user_id}"
// Business rule: Cart requires authentication - only authenticated users can have carts
type Cart struct {
	UserID    string              `json:"user_id"`              // User ID (required - authentication required)
	SessionID string              `json:"session_id,omitempty"` // Deprecated: No longer used, kept for backward compatibility
	Items     map[uint]*CartItem  `json:"items"`                // Map of product_id -> CartItem
	Total     float64             `json:"total"`                // Total price of all items
	UpdatedAt int64               `json:"updated_at"`            // Unix timestamp
}

// CartRepository defines the interface for cart data access
// Cart is stored in Redis, not PostgreSQL (for Gate 4)
// Business rule: Cart requires authentication - only userID is accepted (sessionID is deprecated)
type CartRepository interface {
	GetCart(userID string) (*Cart, error)
	SaveCart(cart *Cart) error
	DeleteCart(userID string) error
	ClearCartItems(userID string) error
}


