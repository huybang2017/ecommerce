package domain

// VariationOption represents a value for a variation (e.g. "M", "L", "Red", "Blue")
// Following db-diagram.db schema (SOURCE OF TRUTH)
type VariationOption struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	VariationID uint   `gorm:"index;not null" json:"variation_id"`
	Value       string `gorm:"size:50;not null" json:"value"` // "M", "L", "XL", "Red", "Blue"
}

// TableName specifies the table name for GORM
func (VariationOption) TableName() string {
	return "variation_option"
}

// VariationOptionRepository defines the interface for variation option data access
type VariationOptionRepository interface {
	Create(option *VariationOption) error
	Update(option *VariationOption) error
	GetByID(id uint) (*VariationOption, error)
	GetByVariationID(variationID uint) ([]*VariationOption, error)
	Delete(id uint) error
}

