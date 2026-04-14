package mapper

import (
	"product-service/internal/domain"
	"product-service/internal/dto/request"
	"product-service/internal/dto/response"
)

func CategoryToResponse(category *domain.Category) *response.CategoryResponse {
	if category == nil {
		return nil
	}

	return &response.CategoryResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func CategoriesToResponse(categories []*domain.Category) []*response.CategoryResponse {
	responses := make([]*response.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = CategoryToResponse(category)
	}
	return responses
}

func CreateParentCategoryRequestToDomain(req *request.CreateParentCategoryRequest) *domain.Category {
	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    nil,
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	} else {
		category.IsActive = true
	}

	return category
}

func CreateCategoryRequestToDomain(req *request.CreateCategoryRequest) *domain.Category {
	return &domain.Category{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		IsActive:    req.IsActive,
	}
}

func CreateChildCategoryRequestToDomain(req *request.CreateChildCategoryRequest, parentID uint) *domain.Category {
	category := &domain.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    &parentID,
	}

	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	} else {
		category.IsActive = true
	}

	return category
}

func UpdateCategoryRequestToDomain(existing *domain.Category, req *request.UpdateCategoryRequest) *domain.Category {
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Slug != "" {
		existing.Slug = req.Slug
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	return existing
}

func CategoryQueryRequestToDomain(req *request.CategoryQueryRequest) domain.CategoryQueryParams {
	return domain.CategoryQueryParams{
		Page:      req.Page,
		Limit:     req.Limit,
		Search:    req.Search,
		ParentID:  req.ParentID,
		IsActive:  req.IsActive,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}
}

func BuildCategoryTree(categories []*domain.Category) []*response.CategoryResponse {
	categoryMap := make(map[uint]*response.CategoryResponse)
	var roots []*response.CategoryResponse

	for _, c := range categories {
		categoryMap[c.ID] = &response.CategoryResponse{
			ID:          c.ID,
			ParentID:    c.ParentID,
			Name:        c.Name,
			Slug:        c.Slug,
			Description: c.Description,
			ImageURL:    c.ImageURL,
			IsActive:    c.IsActive,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		}
	}

	for _, c := range categories {
		node := categoryMap[c.ID]

		if c.ParentID == nil {
			roots = append(roots, node)
		} else {
			parent, ok := categoryMap[*c.ParentID]
			if ok {
				parent.Children = append(parent.Children, node)
			} else {
				// Keep orphaned categories visible instead of dropping them.
				roots = append(roots, node)
			}
		}
	}

	return roots
}
