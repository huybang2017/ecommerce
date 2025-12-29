package domain

// ProductItem represents a SKU - a specific variation combination with its own price and stock
// Example: Product "T-Shirt" -> ProductItem "T-Shirt Size M Color Red" (SKU: TS-M-RED-001)
// Following db-diagram.db schema (SOURCE OF TRUTH)
type ProductItem struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	ProductID  uint    `gorm:"index;not null" json:"product_id"`
	SKUCode    string  `gorm:"column:sku_code;size:50;uniqueIndex;not null" json:"sku_code"` // Unique SKU
	ImageURL   string  `gorm:"column:image_url;size:255" json:"image_url"` // Image for this specific SKU
	Price      float64 `gorm:"type:decimal(15,2);not null" json:"price"` // Price for this SKU
	QtyInStock int     `gorm:"column:qty_in_stock;default:0" json:"qty_in_stock"` // Stock for this SKU
	Status     string  `gorm:"size:20;default:'ACTIVE'" json:"status"` // ACTIVE, OUT_OF_STOCK, DISABLED
}

// TableName specifies the table name for GORM
func (ProductItem) TableName() string {
	return "product_item"
}

// ProductItemRepository defines the interface for product item (SKU) data access
type ProductItemRepository interface {
	Create(item *ProductItem) error
	Update(item *ProductItem) error
	GetByID(id uint) (*ProductItem, error)
	GetBySKUCode(skuCode string) (*ProductItem, error)
	GetByProductID(productID uint) ([]*ProductItem, error)
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error // Atomic stock update
}

