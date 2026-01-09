package handler

import (
	"product-service/internal/domain"
)

// CategoryResponse is the DTO for API responses (prevents domain leak)
type CategoryResponse struct {
	ID          uint   `json:"id"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// ❌ NO Parent *CategoryResponse (prevents circular JSON)
	// ❌ NO Children []CategoryResponse (frontend should request separately)
}

// ToCategoryResponse converts domain.Category to CategoryResponse DTO
func ToCategoryResponse(category *domain.Category) *CategoryResponse {
	if category == nil {
		return nil
	}
	return &CategoryResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToCategoryResponses converts slice of domain.Category to slice of CategoryResponse
func ToCategoryResponses(categories []*domain.Category) []*CategoryResponse {
	responses := make([]*CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = ToCategoryResponse(cat)
	}
	return responses
}

// CategoryWithParentResponse includes parent info for breadcrumb use case
type CategoryWithParentResponse struct {
	CategoryResponse
	Parent *CategoryResponse `json:"parent,omitempty"`
}

// CategoryWithChildrenResponse includes direct children for sidebar use case
type CategoryWithChildrenResponse struct {
	CategoryResponse
	Children []*CategoryResponse `json:"children,omitempty"`
}
