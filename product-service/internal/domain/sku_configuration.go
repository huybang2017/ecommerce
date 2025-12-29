package domain

// SKUConfiguration links a ProductItem (SKU) with VariationOptions
// Example: SKU "TS-M-RED-001" = Size M (option_id=1) + Color Red (option_id=5)
// This is a many-to-many relationship with composite primary key
// Following db-diagram.db schema (SOURCE OF TRUTH)
type SKUConfiguration struct {
	ProductItemID     uint `gorm:"primaryKey;autoIncrement:false" json:"product_item_id"`
	VariationOptionID uint `gorm:"primaryKey;autoIncrement:false" json:"variation_option_id"`
}

// TableName specifies the table name for GORM
func (SKUConfiguration) TableName() string {
	return "sku_configuration"
}

// SKUConfigurationRepository defines the interface for SKU configuration data access
type SKUConfigurationRepository interface {
	Create(config *SKUConfiguration) error
	CreateBatch(configs []*SKUConfiguration) error // Bulk insert for multiple options
	GetByProductItemID(productItemID uint) ([]*SKUConfiguration, error)
	GetByVariationOptionID(optionID uint) ([]*SKUConfiguration, error)
	Delete(productItemID uint, variationOptionID uint) error
	DeleteByProductItemID(productItemID uint) error // Delete all configs for a SKU
}

