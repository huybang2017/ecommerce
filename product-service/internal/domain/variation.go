package domain

// Variation represents a type of product variation (e.g. Size, Color, Storage)
// Following db-diagram.db schema (SOURCE OF TRUTH)
type Variation struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `gorm:"index;not null" json:"product_id"`
	Name      string `gorm:"size:50;not null" json:"name"` // "Size", "Color", "Storage"
}

// TableName specifies the table name for GORM
func (Variation) TableName() string {
	return "variation"
}

// VariationRepository defines the interface for variation data access
type VariationRepository interface {
	Create(variation *Variation) error
	Update(variation *Variation) error
	GetByID(id uint) (*Variation, error)
	GetByProductID(productID uint) ([]*Variation, error)
	Delete(id uint) error
}

