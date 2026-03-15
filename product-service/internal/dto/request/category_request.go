package request

type CreateParentCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Slug        string `json:"slug" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	IsActive    *bool  `json:"is_active"`
}

type CreateCategoryRequest struct {
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Slug        string `json:"slug" binding:"omitempty,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	IsActive    bool   `json:"is_active"`
}

type CreateChildCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Slug        string `json:"slug" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Slug        string `json:"slug" binding:"omitempty,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	ParentID    *uint  `json:"parent_id"`
	IsActive    *bool  `json:"is_active"`
}

type CategoryQueryRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search    string `form:"search" binding:"omitempty,max=100"`
	ParentID  *uint  `form:"parent_id"`
	IsActive  *bool  `form:"is_active"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=id name created_at updated_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc ASC DESC"`
}
