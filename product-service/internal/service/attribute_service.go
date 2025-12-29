package service

import (
	"errors"
	"fmt"
	"product-service/internal/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AttributeService contains the business logic for EAV attributes
type AttributeService struct {
	categoryAttrRepo domain.CategoryAttributeRepository
	productAttrRepo  domain.ProductAttributeValueRepository
	categoryRepo     domain.CategoryRepository
	productRepo      domain.ProductRepository
	logger           *zap.Logger
}

// NewAttributeService creates a new attribute service
func NewAttributeService(
	categoryAttrRepo domain.CategoryAttributeRepository,
	productAttrRepo domain.ProductAttributeValueRepository,
	categoryRepo domain.CategoryRepository,
	productRepo domain.ProductRepository,
	logger *zap.Logger,
) *AttributeService {
	return &AttributeService{
		categoryAttrRepo: categoryAttrRepo,
		productAttrRepo:  productAttrRepo,
		categoryRepo:     categoryRepo,
		productRepo:      productRepo,
		logger:           logger,
	}
}

// CreateCategoryAttributeRequest represents the request to create a category attribute
type CreateCategoryAttributeRequest struct {
	CategoryID    uint   `json:"category_id" binding:"required"`
	AttributeName string `json:"attribute_name" binding:"required,min=2,max=50"`
	InputType     string `json:"input_type" binding:"required"` // text, number, select, checkbox
	IsMandatory   bool   `json:"is_mandatory"`
	IsFilterable  bool   `json:"is_filterable"`
}

// SetProductAttributesRequest represents the request to set attributes for a product
type SetProductAttributesRequest struct {
	Attributes map[uint]string `json:"attributes" binding:"required"` // map[attribute_id]value
}

// CreateCategoryAttribute creates a new attribute for a category
func (s *AttributeService) CreateCategoryAttribute(req *CreateCategoryAttributeRequest) (*domain.CategoryAttribute, error) {
	// Validate category exists
	_, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Validate input type
	validInputTypes := map[string]bool{
		"text": true, "number": true, "select": true, "checkbox": true,
	}
	if !validInputTypes[req.InputType] {
		return nil, errors.New("invalid input_type: must be text, number, select, or checkbox")
	}

	attr := &domain.CategoryAttribute{
		CategoryID:    req.CategoryID,
		AttributeName: req.AttributeName,
		InputType:     req.InputType,
		IsMandatory:   req.IsMandatory,
		IsFilterable:  req.IsFilterable,
	}

	if err := s.categoryAttrRepo.Create(attr); err != nil {
		s.logger.Error("failed to create category attribute", zap.Error(err))
		return nil, fmt.Errorf("failed to create category attribute: %w", err)
	}

	s.logger.Info("category attribute created", zap.Uint("attr_id", attr.ID), zap.String("name", attr.AttributeName))

	return attr, nil
}

// GetCategoryAttributes retrieves all attributes for a category
func (s *AttributeService) GetCategoryAttributes(categoryID uint) ([]*domain.CategoryAttribute, error) {
	attrs, err := s.categoryAttrRepo.GetByCategoryID(categoryID)
	if err != nil {
		s.logger.Error("failed to get category attributes", zap.Error(err))
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}
	return attrs, nil
}

// SetProductAttributes sets attributes for a product
// Business logic:
// 1. Validate product exists and get its category
// 2. Validate all attribute_ids belong to the product's category
// 3. Check mandatory attributes are provided
// 4. Delete old attribute values
// 5. Create new attribute values
func (s *AttributeService) SetProductAttributes(productID uint, req *SetProductAttributesRequest) error {
	// 1. Get product and its category
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return fmt.Errorf("failed to get product: %w", err)
	}

	if product.CategoryID == nil {
		return errors.New("product must have a category to set attributes")
	}

	// 2. Get category attributes
	categoryAttrs, err := s.categoryAttrRepo.GetByCategoryID(*product.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to get category attributes: %w", err)
	}

	// Create map of valid attribute IDs for this category
	validAttrIDs := make(map[uint]*domain.CategoryAttribute)
	mandatoryAttrIDs := make(map[uint]bool)
	for _, attr := range categoryAttrs {
		validAttrIDs[attr.ID] = attr
		if attr.IsMandatory {
			mandatoryAttrIDs[attr.ID] = true
		}
	}

	// 3. Validate provided attributes
	for attrID := range req.Attributes {
		if _, exists := validAttrIDs[attrID]; !exists {
			return fmt.Errorf("attribute_id %d does not belong to product's category", attrID)
		}
	}

	// 4. Check mandatory attributes are provided
	for attrID := range mandatoryAttrIDs {
		if _, provided := req.Attributes[attrID]; !provided {
			return fmt.Errorf("mandatory attribute_id %d is missing", attrID)
		}
	}

	// 5. Delete old attribute values
	if err := s.productAttrRepo.DeleteByProductID(productID); err != nil {
		s.logger.Error("failed to delete old product attributes", zap.Error(err))
		return fmt.Errorf("failed to delete old attributes: %w", err)
	}

	// 6. Create new attribute values
	var values []*domain.ProductAttributeValue
	for attrID, value := range req.Attributes {
		values = append(values, &domain.ProductAttributeValue{
			ProductID:   productID,
			AttributeID: attrID,
			Value:       value,
		})
	}

	if len(values) > 0 {
		if err := s.productAttrRepo.CreateBatch(values); err != nil {
			s.logger.Error("failed to create product attributes", zap.Error(err))
			return fmt.Errorf("failed to create product attributes: %w", err)
		}
	}

	s.logger.Info("product attributes set", zap.Uint("product_id", productID), zap.Int("count", len(values)))

	return nil
}

// GetProductAttributes retrieves all attributes for a product
// Returns map[attribute_name]value for easy display
func (s *AttributeService) GetProductAttributes(productID uint) (map[string]string, error) {
	// Get product attribute values
	values, err := s.productAttrRepo.GetByProductID(productID)
	if err != nil {
		s.logger.Error("failed to get product attributes", zap.Error(err))
		return nil, fmt.Errorf("failed to get product attributes: %w", err)
	}

	// Get attribute names
	result := make(map[string]string)
	for _, val := range values {
		attr, err := s.categoryAttrRepo.GetByID(val.AttributeID)
		if err != nil {
			s.logger.Warn("failed to get attribute name", zap.Uint("attr_id", val.AttributeID))
			continue
		}
		result[attr.AttributeName] = val.Value
	}

	return result, nil
}

// UpdateCategoryAttribute updates a category attribute
func (s *AttributeService) UpdateCategoryAttribute(id uint, name, inputType string, isMandatory, isFilterable bool) (*domain.CategoryAttribute, error) {
	attr, err := s.categoryAttrRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category attribute not found")
		}
		return nil, fmt.Errorf("failed to get category attribute: %w", err)
	}

	// Update fields
	if name != "" {
		attr.AttributeName = name
	}
	if inputType != "" {
		// Validate input type
		validInputTypes := map[string]bool{
			"text": true, "number": true, "select": true, "checkbox": true,
		}
		if !validInputTypes[inputType] {
			return nil, errors.New("invalid input_type")
		}
		attr.InputType = inputType
	}
	attr.IsMandatory = isMandatory
	attr.IsFilterable = isFilterable

	if err := s.categoryAttrRepo.Update(attr); err != nil {
		s.logger.Error("failed to update category attribute", zap.Error(err))
		return nil, fmt.Errorf("failed to update category attribute: %w", err)
	}

	s.logger.Info("category attribute updated", zap.Uint("attr_id", attr.ID))

	return attr, nil
}

// DeleteCategoryAttribute deletes a category attribute
func (s *AttributeService) DeleteCategoryAttribute(id uint) error {
	if err := s.categoryAttrRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete category attribute", zap.Error(err))
		return fmt.Errorf("failed to delete category attribute: %w", err)
	}

	s.logger.Info("category attribute deleted", zap.Uint("attr_id", id))

	return nil
}

