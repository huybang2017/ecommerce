package service

import (
	"errors"
	"fmt"
	"product-service/internal/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductItemService contains the business logic for product item (SKU) operations
type ProductItemService struct {
	productItemRepo domain.ProductItemRepository
	variationRepo   domain.VariationRepository
	variationOptRepo domain.VariationOptionRepository
	skuConfigRepo   domain.SKUConfigurationRepository
	productRepo     domain.ProductRepository
	logger          *zap.Logger
}

// NewProductItemService creates a new product item service
func NewProductItemService(
	productItemRepo domain.ProductItemRepository,
	variationRepo domain.VariationRepository,
	variationOptRepo domain.VariationOptionRepository,
	skuConfigRepo domain.SKUConfigurationRepository,
	productRepo domain.ProductRepository,
	logger *zap.Logger,
) *ProductItemService {
	return &ProductItemService{
		productItemRepo:  productItemRepo,
		variationRepo:    variationRepo,
		variationOptRepo: variationOptRepo,
		skuConfigRepo:    skuConfigRepo,
		productRepo:      productRepo,
		logger:           logger,
	}
}

// CreateProductItemRequest represents the request to create a new product item (SKU)
type CreateProductItemRequest struct {
	ProductID        uint     `json:"product_id" binding:"required"`
	SKUCode          string   `json:"sku_code" binding:"required"`
	ImageURL         string   `json:"image_url"`
	Price            float64  `json:"price" binding:"required,min=0"`
	QtyInStock       int      `json:"qty_in_stock"`
	VariationOptions []uint   `json:"variation_options"` // List of variation_option_ids (e.g. [1, 5] = Size M + Color Red)
}

// UpdateProductItemRequest represents the request to update a product item
type UpdateProductItemRequest struct {
	ImageURL   string  `json:"image_url"`
	Price      float64 `json:"price" binding:"omitempty,min=0"`
	QtyInStock int     `json:"qty_in_stock"`
	Status     string  `json:"status"`
}

// CreateProductItem creates a new product item (SKU) with variation options
// Business logic:
// 1. Validate product exists
// 2. Validate SKU code is unique
// 3. Validate variation options belong to product's variations
// 4. Check duplicate combination (same variation options already exist)
// 5. Create product item
// 6. Create SKU configurations (link SKU with variation options)
func (s *ProductItemService) CreateProductItem(req *CreateProductItemRequest) (*domain.ProductItem, error) {
	// 1. Validate product exists
	_, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// 2. Check if SKU code already exists
	existing, err := s.productItemRepo.GetBySKUCode(req.SKUCode)
	if err == nil && existing != nil {
		return nil, errors.New("SKU code already exists")
	}

	// 3. Validate variation options belong to product's variations
	if len(req.VariationOptions) > 0 {
		productVariations, err := s.variationRepo.GetByProductID(req.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to get product variations: %w", err)
		}

		// Create map of variation IDs for validation
		variationIDs := make(map[uint]bool)
		for _, v := range productVariations {
			variationIDs[v.ID] = true
		}

		// Validate each variation option belongs to product's variations
		for _, optionID := range req.VariationOptions {
			option, err := s.variationOptRepo.GetByID(optionID)
			if err != nil {
				return nil, fmt.Errorf("variation option %d not found", optionID)
			}
			if !variationIDs[option.VariationID] {
				return nil, fmt.Errorf("variation option %d does not belong to product %d", optionID, req.ProductID)
			}
		}

		// TODO: 4. Check duplicate combination (same variation options already exist)
		// This requires more complex logic to query existing SKU configurations
	}

	// 5. Create product item
	item := &domain.ProductItem{
		ProductID:  req.ProductID,
		SKUCode:    req.SKUCode,
		ImageURL:   req.ImageURL,
		Price:      req.Price,
		QtyInStock: req.QtyInStock,
		Status:     "ACTIVE",
	}

	if err := s.productItemRepo.Create(item); err != nil {
		s.logger.Error("failed to create product item", zap.Error(err))
		return nil, fmt.Errorf("failed to create product item: %w", err)
	}

	s.logger.Info("product item created", zap.Uint("product_item_id", item.ID), zap.String("sku_code", item.SKUCode))

	// 6. Create SKU configurations (link SKU with variation options)
	if len(req.VariationOptions) > 0 {
		var configs []*domain.SKUConfiguration
		for _, optionID := range req.VariationOptions {
			configs = append(configs, &domain.SKUConfiguration{
				ProductItemID:     item.ID,
				VariationOptionID: optionID,
			})
		}

		if err := s.skuConfigRepo.CreateBatch(configs); err != nil {
			// Rollback: delete the product item if SKU configuration fails
			s.productItemRepo.Delete(item.ID)
			s.logger.Error("failed to create SKU configurations", zap.Error(err))
			return nil, fmt.Errorf("failed to create SKU configurations: %w", err)
		}

		s.logger.Info("SKU configurations created", zap.Uint("product_item_id", item.ID), zap.Int("count", len(configs)))
	}

	return item, nil
}

// UpdateProductItem updates an existing product item
func (s *ProductItemService) UpdateProductItem(id uint, req *UpdateProductItemRequest) (*domain.ProductItem, error) {
	// Get existing item
	item, err := s.productItemRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product item not found")
		}
		return nil, fmt.Errorf("failed to get product item: %w", err)
	}

	// Update fields
	if req.ImageURL != "" {
		item.ImageURL = req.ImageURL
	}
	if req.Price > 0 {
		item.Price = req.Price
	}
	if req.QtyInStock >= 0 {
		item.QtyInStock = req.QtyInStock
	}
	if req.Status != "" {
		// Validate status
		if req.Status != "ACTIVE" && req.Status != "OUT_OF_STOCK" && req.Status != "DISABLED" {
			return nil, errors.New("invalid status")
		}
		item.Status = req.Status
	}

	if err := s.productItemRepo.Update(item); err != nil {
		s.logger.Error("failed to update product item", zap.Error(err))
		return nil, fmt.Errorf("failed to update product item: %w", err)
	}

	s.logger.Info("product item updated", zap.Uint("product_item_id", item.ID))

	return item, nil
}

// GetProductItem retrieves a product item by ID
func (s *ProductItemService) GetProductItem(id uint) (*domain.ProductItem, error) {
	item, err := s.productItemRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product item not found")
		}
		return nil, fmt.Errorf("failed to get product item: %w", err)
	}
	return item, nil
}

// GetProductItemBySKU retrieves a product item by SKU code
func (s *ProductItemService) GetProductItemBySKU(skuCode string) (*domain.ProductItem, error) {
	item, err := s.productItemRepo.GetBySKUCode(skuCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product item not found")
		}
		return nil, fmt.Errorf("failed to get product item: %w", err)
	}
	return item, nil
}

// GetProductItems retrieves all product items (SKUs) for a product
func (s *ProductItemService) GetProductItems(productID uint) ([]*domain.ProductItem, error) {
	items, err := s.productItemRepo.GetByProductID(productID)
	if err != nil {
		s.logger.Error("failed to get product items", zap.Error(err))
		return nil, fmt.Errorf("failed to get product items: %w", err)
	}
	return items, nil
}

// DeleteProductItem deletes a product item and its SKU configurations
func (s *ProductItemService) DeleteProductItem(id uint) error {
	// Delete SKU configurations first (foreign key constraint)
	if err := s.skuConfigRepo.DeleteByProductItemID(id); err != nil {
		s.logger.Error("failed to delete SKU configurations", zap.Error(err))
		return fmt.Errorf("failed to delete SKU configurations: %w", err)
	}

	// Delete product item
	if err := s.productItemRepo.Delete(id); err != nil {
		s.logger.Error("failed to delete product item", zap.Error(err))
		return fmt.Errorf("failed to delete product item: %w", err)
	}

	s.logger.Info("product item deleted", zap.Uint("product_item_id", id))

	return nil
}

