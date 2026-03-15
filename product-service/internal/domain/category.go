package domain

import (
	"context"
	"time"
)

type Category struct {
	ID          uint
	ParentID    *uint
	Name        string
	Slug        string
	Description string
	ImageURL    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Category) TableName() string {
	return "categories"
}

type CategoryQueryParams struct {
	Page      int
	Limit     int
	Search    string
	ParentID  *uint
	IsActive  *bool
	SortBy    string
	SortOrder string
}
type CategoryRepository interface {
	Create(category *Category) error
	Update(category *Category) error
	GetByID(id uint) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	GetAll() ([]*Category, error)
	GetAdminList(ctx context.Context, params CategoryQueryParams) ([]*Category, int64, error)
	GetChildren(parentID uint) ([]*Category, error)
	GetAdminCategoryParent(ctx context.Context) ([]*Category, error)
	Delete(id uint) error
}
