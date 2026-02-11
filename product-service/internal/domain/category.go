package domain

import (
	"context"
	"time"
)

// Category represents the category domain entity
// Schema: db-diagram.db (SOURCE OF TRUTH)
// NOTE: NO Parent/Children to avoid circular reference and N+1 queries
type Category struct {
	ID          uint      `json:"id"`
	ParentID    *uint     `json:"parent_id,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CategoryQueryParams struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	Status    string `json:"status"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type AdminCategoryResponse struct {
	Data       []*Category `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func (Category) TableName() string {
	return "categories"
}

type CategoryRepository interface {
	Create(category *Category) error
	Update(category *Category) error
	GetByID(id uint) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	GetAll() ([]*Category, error)
	GetAdminList(ctx context.Context, params CategoryQueryParams) ([]*Category, int64, error)
	GetChildren(parentID uint) ([]*Category, error)
	Delete(id uint) error
}
