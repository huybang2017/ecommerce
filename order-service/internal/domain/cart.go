package domain

import (
	"errors"
)

// CartItem represents a single item in the shopping cart
// Stored in Redis: MINIMAL data only (product_item_id, quantity, is_selected)
type CartItem struct {
	// ✅ STORED in Redis
	ProductItemID uint `json:"product_item_id"`
	Quantity      int  `json:"quantity"`
	IsSelected    bool `json:"is_selected"`

	// ❌ NOT stored in Redis - Fetched from Product Service on-demand
	ShopID      uint    `json:"shop_id,omitempty" redis:"-"`
	ProductName string  `json:"product_name,omitempty" redis:"-"`
	SKUCode     string  `json:"sku_code,omitempty" redis:"-"`
	ImageURL    string  `json:"image_url,omitempty" redis:"-"`
	Price       float64 `json:"price,omitempty" redis:"-"`
}

// ShoppingCart represents a shopping cart
// Stored in Redis with key: "cart:user:{user_id}"
type ShoppingCart struct {
	UserID  string      `json:"user_id"`
	Items   []*CartItem `json:"items"`
	Version int         `json:"version"`

	// ✅ COMPUTED FIELDS - Cart Service responsibility ONLY
	// These are DISPLAY metrics for Cart UI/Badge, NOT checkout logic
	ItemCount          int     `json:"item_count" redis:"-"`           // Total items (rows)
	TotalQuantity      int     `json:"total_quantity" redis:"-"`       // Sum of all quantities
	SelectedItemCount  int     `json:"selected_item_count" redis:"-"`  // Selected items (rows)
	SelectedQuantity   int     `json:"selected_quantity" redis:"-"`    // Sum of selected quantities
	TotalPrice         float64 `json:"total_price" redis:"-"`          // Total price (all)
	SelectedTotalPrice float64 `json:"selected_total_price" redis:"-"` // Total price (selected)
}

// CalculateTotals computes CART-LEVEL metrics ONLY
// This method focuses on counting and summing - NOT business logic like grouping
// Grouping by shop, shipping calculation, etc. belong to Checkout Service
func (c *ShoppingCart) CalculateTotals() {
	// Reset
	c.ItemCount = 0
	c.TotalQuantity = 0
	c.SelectedItemCount = 0
	c.SelectedQuantity = 0
	c.TotalPrice = 0
	c.SelectedTotalPrice = 0

	for _, item := range c.Items {
		// ✅ Cart Service chỉ quan tâm: đếm và tính tổng
		qty := item.Quantity
		linePrice := item.Price * float64(qty)

		// Count all items
		c.ItemCount++
		c.TotalQuantity += qty
		c.TotalPrice += linePrice

		// Count selected items
		if item.IsSelected {
			c.SelectedItemCount++
			c.SelectedQuantity += qty
			c.SelectedTotalPrice += linePrice
		}
	}
}

// ==========================================
// HELPER METHODS - CART DOMAIN ONLY
// ==========================================

// IsEmpty checks if cart has no items
func (c *ShoppingCart) IsEmpty() bool {
	return len(c.Items) == 0
}

// HasSelectedItems checks if any items are selected
func (c *ShoppingCart) HasSelectedItems() bool {
	for _, item := range c.Items {
		if item.IsSelected {
			return true
		}
	}
	return false
}

// GetSelectedItems returns only selected items (for passing to Checkout Service)
func (c *ShoppingCart) GetSelectedItems() []*CartItem {
	selected := make([]*CartItem, 0)
	for _, item := range c.Items {
		if item.IsSelected {
			selected = append(selected, item)
		}
	}
	return selected
}

// FindItemByProductItemID finds item by product item ID
func (c *ShoppingCart) FindItemByProductItemID(productItemID uint) *CartItem {
	for i := range c.Items {
		if c.Items[i].ProductItemID == productItemID {
			return c.Items[i]
		}
	}
	return nil
}

// Validate validates cart item
func (ci *CartItem) Validate() error {
	if ci.ProductItemID == 0 {
		return ErrInvalidProductItem
	}
	if ci.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	if ci.Quantity > 999 {
		return ErrQuantityExceedsLimit
	}
	return nil
}

// ==========================================
// DOMAIN ERRORS
// ==========================================

var (
	ErrCartNotFound         = errors.New("cart not found")
	ErrCartItemNotFound     = errors.New("cart item not found")
	ErrInvalidProductItem   = errors.New("invalid product item")
	ErrInvalidQuantity      = errors.New("quantity must be greater than 0")
	ErrQuantityExceedsLimit = errors.New("quantity exceeds maximum limit (999)")
	ErrCartEmpty            = errors.New("cart is empty")
	ErrNoItemsSelected      = errors.New("no items selected for checkout")
	ErrProductOutOfStock    = errors.New("product is out of stock")
	ErrInsufficientStock    = errors.New("insufficient stock for requested quantity")
)
