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
	productItemRepo  domain.ProductItemRepository
	variationRepo    domain.VariationRepository
	variationOptRepo domain.VariationOptionRepository
	skuConfigRepo    domain.SKUConfigurationRepository
	productRepo      domain.ProductRepository
	logger           *zap.Logger
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
	ProductID        uint    `json:"product_id" binding:"required"`
	SKUCode          string  `json:"sku_code" binding:"required"`
	ImageURL         string  `json:"image_url"`
	Price            float64 `json:"price" binding:"required,min=0"`
	QtyInStock       int     `json:"qty_in_stock"`
	VariationOptions []uint  `json:"variation_options"` // List of variation_option_ids (e.g. [1, 5] = Size M + Color Red)
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

// ProductItemWithVariations includes variation option IDs for UI matching
type ProductItemWithVariations struct {
	ID                 uint    `json:"id"`
	ProductID          uint    `json:"product_id"`
	SKUCode            string  `json:"sku_code"`
	ImageURL           string  `json:"image_url"`
	Price              float64 `json:"price"`
	QtyInStock         int     `json:"qty_in_stock"`
	Status             string  `json:"status"`
	VariationOptionIDs []uint  `json:"variation_option_ids"`
}

// GetProductItemsWithVariations retrieves product items with their variation option IDs
// This is used for variation selector UI (Shopee-style)
func (s *ProductItemService) GetProductItemsWithVariations(productID uint) ([]*ProductItemWithVariations, error) {
	// Get all product items
	items, err := s.productItemRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	result := make([]*ProductItemWithVariations, 0, len(items))
	for _, item := range items {
		// Get SKU configurations (variation options)
		configs, err := s.skuConfigRepo.GetByProductItemID(item.ID)
		if err != nil {
			s.logger.Warn("Failed to get SKU configurations",
				zap.Uint("product_item_id", item.ID),
				zap.Error(err))
			configs = []*domain.SKUConfiguration{} // Continue with empty
		}

		// Extract variation option IDs
		optionIDs := make([]uint, len(configs))
		for i, config := range configs {
			optionIDs[i] = config.VariationOptionID
		}

		result = append(result, &ProductItemWithVariations{
			ID:                 item.ID,
			ProductID:          item.ProductID,
			SKUCode:            item.SKUCode,
			ImageURL:           item.ImageURL,
			Price:              item.Price,
			QtyInStock:         item.QtyInStock,
			Status:             item.Status,
			VariationOptionIDs: optionIDs,
		})
	}

	return result, nil
}

// ProductItemWithProduct represents a product item with nested product info
type ProductItemWithProduct struct {
	ID         uint    `json:"id"`
	ProductID  uint    `json:"product_id"`
	SKUCode    string  `json:"sku_code"`
	ImageURL   string  `json:"image_url"`
	Price      float64 `json:"price"`
	QtyInStock int     `json:"qty_in_stock"`
	Status     string  `json:"status"`
	Product    *struct {
		ID     uint   `json:"id"`
		ShopID uint   `json:"shop_id"`
		Name   string `json:"name"`
	} `json:"product"`
}

// GetProductItemsWithProduct retrieves multiple product items by IDs with product details
// Used by order-service/cart-service for batch fetching
func (s *ProductItemService) GetProductItemsWithProduct(ids []uint) ([]*ProductItemWithProduct, error) {
	if len(ids) == 0 {
		return []*ProductItemWithProduct{}, nil
	}

	result := make([]*ProductItemWithProduct, 0, len(ids))

	for _, id := range ids {
		// Get product item
		item, err := s.productItemRepo.GetByID(id)
		if err != nil {
			s.logger.Warn("product item not found", zap.Uint("id", id), zap.Error(err))
			continue // Skip missing items instead of failing entire batch
		}

		// Get product info
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			s.logger.Warn("product not found for item", zap.Uint("product_id", item.ProductID), zap.Error(err))
			continue
		}

		// Build response
		itemWithProduct := &ProductItemWithProduct{
			ID:         item.ID,
			ProductID:  item.ProductID,
			SKUCode:    item.SKUCode,
			ImageURL:   item.ImageURL,
			Price:      item.Price,
			QtyInStock: item.QtyInStock,
			Status:     item.Status,
			Product: &struct {
				ID     uint   `json:"id"`
				ShopID uint   `json:"shop_id"`
				Name   string `json:"name"`
			}{
				ID:     product.ID,
				ShopID: product.ShopID,
				Name:   product.Name,
			},
		}

		result = append(result, itemWithProduct)
	}

	return result, nil
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
