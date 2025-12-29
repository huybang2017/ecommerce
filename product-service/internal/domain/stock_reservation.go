package domain

import "time"

// StockReservation represents a temporary stock hold (stored in Redis)
// Used during checkout flow to prevent overselling
type StockReservation struct {
	OrderID       string    `json:"order_id"`       // Order ID that reserved this stock
	ProductItemID uint      `json:"product_item_id"` // SKU ID
	Quantity      int       `json:"quantity"`       // Reserved quantity
	ExpiresAt     time.Time `json:"expires_at"`     // Expiration time (auto-release after timeout)
}

// StockCheckRequest represents a request to check stock availability
type StockCheckRequest struct {
	Items []StockCheckItem `json:"items" binding:"required"`
}

// StockCheckItem represents a single item to check stock
type StockCheckItem struct {
	ProductItemID uint `json:"product_item_id" binding:"required"`
	Quantity      int  `json:"quantity" binding:"required,min=1"`
}

// StockCheckResponse represents the response for stock check
type StockCheckResponse struct {
	Available         bool                  `json:"available"`
	UnavailableItems  []UnavailableStockItem `json:"unavailable_items,omitempty"`
}

// UnavailableStockItem represents an item that doesn't have enough stock
type UnavailableStockItem struct {
	ProductItemID uint `json:"product_item_id"`
	Requested     int  `json:"requested"`
	Available     int  `json:"available"`
}

// StockReserveRequest represents a request to reserve stock
type StockReserveRequest struct {
	OrderID string            `json:"order_id" binding:"required"`
	Items   []StockReserveItem `json:"items" binding:"required"`
}

// StockReserveItem represents a single item to reserve
type StockReserveItem struct {
	ProductItemID uint `json:"product_item_id" binding:"required"`
	Quantity      int  `json:"quantity" binding:"required,min=1"`
}

// StockDeductRequest represents a request to deduct stock permanently
type StockDeductRequest struct {
	OrderID string           `json:"order_id" binding:"required"`
	Items   []StockDeductItem `json:"items" binding:"required"`
}

// StockDeductItem represents a single item to deduct
type StockDeductItem struct {
	ProductItemID uint `json:"product_item_id" binding:"required"`
	Quantity      int  `json:"quantity" binding:"required,min=1"`
}

// StockReleaseRequest represents a request to release reserved stock
type StockReleaseRequest struct {
	OrderID string `json:"order_id" binding:"required"`
}

