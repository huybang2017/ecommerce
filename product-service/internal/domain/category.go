package domain

import (
	"time"
)

// Category represents the category domain entity
// Supports nested categories via parent_id
// NOTE: Following db-diagram.db schema (SOURCE OF TRUTH)
type Category struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	ParentID    *uint      `gorm:"index" json:"parent_id,omitempty"` // Nullable for root categories
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Name        string     `gorm:"not null" json:"name"`
	ImageURL    string     `gorm:"column:image_url;size:255" json:"image_url"` // THÊM MỚI từ db-diagram.db
	IsActive    bool       `gorm:"column:is_active;default:true" json:"is_active"` // THÊM MỚI từ db-diagram.db
	Slug        string     `gorm:"uniqueIndex;not null" json:"slug"` // GIỮ LẠI để backward compatibility
	Description string     `json:"description"` // GIỮ LẠI để backward compatibility
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "categories"
}

// CategoryRepository defines the interface for category data access
// This is part of the domain layer - it defines WHAT we need, not HOW
type CategoryRepository interface {
	Create(category *Category) error
	Update(category *Category) error
	GetByID(id uint) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	GetAll() ([]*Category, error)
	GetChildren(parentID uint) ([]*Category, error)
	Delete(id uint) error
}

