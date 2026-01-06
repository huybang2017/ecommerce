package domain

type CartRepository interface {
	// Basic operations
	GetCart(userID string) (*ShoppingCart, error)
	SaveCart(cart *ShoppingCart) error
	DeleteCart(userID string) error
	ClearSelectedItems(userID string) error

	// Item operations
	AddItem(userID string, item *CartItem) error
	UpdateItemQuantity(userID string, productItemID uint, quantity int) error
	RemoveItem(userID string, productItemID uint) error

	// Selection operations
	ToggleItemSelection(userID string, productItemID uint) error
	SelectAllItems(userID string, selected bool) error
	GetSelectedItems(userID string) ([]*CartItem, error)

	// Utility
	GetCartItemCount(userID string) (int, error)
}
