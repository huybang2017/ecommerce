package domain

import "time"

// Shop represents a shop in the marketplace
// Business rule: 1 User = 1 Shop (unique constraint on owner_user_id)
// Following Clean Architecture: domain layer has no external dependencies
type Shop struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OwnerUserID  uint      `gorm:"column:owner_user_id;uniqueIndex;not null" json:"owner_user_id"` // 1 User = 1 Shop
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	LogoURL      string    `gorm:"column:logo_url;size:255" json:"logo_url"`
	CoverURL     string    `gorm:"column:cover_url;size:255" json:"cover_url"`
	IsOfficial   bool      `gorm:"column:is_official;default:false" json:"is_official"`
	Rating       float64   `gorm:"type:decimal(2,1);default:0" json:"rating"`
	ResponseRate int       `gorm:"column:response_rate;default:0" json:"response_rate"`
	Status       string    `gorm:"size:20;default:'ACTIVE'" json:"status"` // ACTIVE, SUSPENDED
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Shop) TableName() string {
	return "shop"
}

// ShopRepository defines the interface for shop data access
// This is part of the domain layer - it defines WHAT we need, not HOW
type ShopRepository interface {
	Create(shop *Shop) error
	Update(shop *Shop) error
	GetByID(id uint) (*Shop, error)
	GetByOwnerUserID(ownerUserID uint) (*Shop, error)
	GetAll(page, limit int) ([]*Shop, int64, error)
	GetByStatus(status string, page, limit int) ([]*Shop, int64, error)
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
}

