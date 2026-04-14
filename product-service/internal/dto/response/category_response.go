package response

import "time"

type CategoryResponse struct {
	ID          uint                `json:"id"`
	ParentID    *uint               `json:"parent_id,omitempty"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	ImageURL    string              `json:"image_url"`
	IsActive    bool                `json:"is_active"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Children    []*CategoryResponse `json:"children,omitempty"`
}

type AdminCategoryListResponse struct {
	Data       []*CategoryResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

type CategoryListResponse struct {
	Data []*CategoryResponse `json:"data"`
}

type CategoryTreeResponse struct {
	Data []*CategoryResponse `json:"data"`
}
