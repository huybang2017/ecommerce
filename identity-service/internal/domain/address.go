package domain

// Address represents the core domain entity for user address
// Following Clean Architecture: domain layer has no external dependencies
type Address struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	UserID        uint   `gorm:"column:user_id;index" json:"user_id"`
	RecipientName string `gorm:"column:recipient_name;size:100" json:"recipient_name"`
	PhoneNumber   string `gorm:"column:phone_number;size:20" json:"phone_number"`
	AddressLine   string `gorm:"column:address_line;size:255" json:"address_line"`
	City          string `gorm:"size:100" json:"city"`
	District      string `gorm:"size:100" json:"district"`
	Ward          string `gorm:"size:100" json:"ward"`
	IsDefault     bool   `gorm:"column:is_default;default:false" json:"is_default"`
	Label         string `gorm:"size:20" json:"label"` // HOME, WORK, etc.
}

// TableName specifies the table name for GORM
func (Address) TableName() string {
	return "address"
}

// AddressRepository defines the interface for address data access
// This is part of the domain layer - it defines WHAT we need, not HOW
type AddressRepository interface {
	Create(address *Address) error
	Update(address *Address) error
	GetByID(id uint) (*Address, error)
	GetByUserID(userID uint) ([]*Address, error)
	GetDefaultByUserID(userID uint) (*Address, error)
	Delete(id uint) error
	SetDefault(userID uint, addressID uint) error
}

