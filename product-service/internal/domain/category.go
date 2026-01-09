package domain

import (
	"time"
)

// Category represents the category domain entity
// Schema: db-diagram.db (SOURCE OF TRUTH)
// NOTE: NO Parent/Children to avoid circular reference and N+1 queries
type Category struct {
	ID          uint      `json:"id"`
	ParentID    *uint     `json:"parent_id,omitempty"` // Nullable for root categories
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`        // Backward compatibility
	Description string    `json:"description"` // Backward compatibility
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// ❌ Removed: Parent *Category (circular reference)
	// ❌ Removed: Children []Category (N+1 problem)
	// ✅ Use repository methods to get parent/children when needed
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
