package domain

// CategoryAttribute defines an attribute that products in a category must/can have
// Example: Category "Điện thoại" has attributes: "RAM", "Màn hình", "Pin"
// Following db-diagram.db schema (SOURCE OF TRUTH)
type CategoryAttribute struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	CategoryID    uint   `gorm:"column:category_id;index;not null" json:"category_id"`
	AttributeName string `gorm:"column:attribute_name;size:50;not null" json:"attribute_name"` // "RAM", "Màn hình"
	InputType     string `gorm:"column:input_type;size:20;not null" json:"input_type"` // text, number, select, checkbox
	IsMandatory   bool   `gorm:"column:is_mandatory;default:false" json:"is_mandatory"` // Bắt buộc điền?
	IsFilterable  bool   `gorm:"column:is_filterable;default:false" json:"is_filterable"` // Hiển thị ở bộ lọc?
}

// TableName specifies the table name for GORM
func (CategoryAttribute) TableName() string {
	return "category_attribute"
}

// CategoryAttributeRepository defines the interface for category attribute data access
type CategoryAttributeRepository interface {
	Create(attr *CategoryAttribute) error
	Update(attr *CategoryAttribute) error
	GetByID(id uint) (*CategoryAttribute, error)
	GetByCategoryID(categoryID uint) ([]*CategoryAttribute, error)
	GetFilterablesByCategoryID(categoryID uint) ([]*CategoryAttribute, error) // Chỉ lấy attributes có thể filter
	Delete(id uint) error
}

