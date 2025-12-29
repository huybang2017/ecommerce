package domain

// ProductAttributeValue stores the value of an attribute for a specific product
// Example: Product iPhone 15 has RAM = "8GB", Màn hình = "6.1 inch"
// Following db-diagram.db schema (SOURCE OF TRUTH)
// NOTE: Cần compound index (attribute_id, value) cho tìm kiếm nhanh
type ProductAttributeValue struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ProductID   uint   `gorm:"index;not null" json:"product_id"` // Index for product queries
	AttributeID uint   `gorm:"column:attribute_id;index;not null" json:"attribute_id"` // Index for attribute queries
	Value       string `gorm:"size:255;not null" json:"value"` // "8GB", "6.1 inch", "Xanh"
}

// TableName specifies the table name for GORM
func (ProductAttributeValue) TableName() string {
	return "product_attribute_value"
}

// ProductAttributeValueRepository defines the interface for product attribute value data access
type ProductAttributeValueRepository interface {
	Create(value *ProductAttributeValue) error
	CreateBatch(values []*ProductAttributeValue) error // Bulk insert
	Update(value *ProductAttributeValue) error
	GetByID(id uint) (*ProductAttributeValue, error)
	GetByProductID(productID uint) ([]*ProductAttributeValue, error)
	GetByAttributeID(attributeID uint) ([]*ProductAttributeValue, error)
	SearchByAttributeValue(attributeID uint, value string) ([]*ProductAttributeValue, error) // Search products by attribute
	Delete(id uint) error
	DeleteByProductID(productID uint) error // Delete all attributes for a product
}

