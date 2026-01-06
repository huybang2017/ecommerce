package service

import (
	"order-service/pkg/product_client"
)

// ==================== CartProductClientAdapter for CartService ====================

type CartProductClientAdapter struct {
	Client *product_client.ProductClient
}

// GetProductItem fetches single product item details (SKU-level) - for CartService display
// Returns display-only DTO without validation fields
func (a *CartProductClientAdapter) GetProductItem(productItemID uint) (*ProductItemDTO, error) {
	items, err := a.Client.GetProductItems([]uint{productItemID})
	if err != nil {
		return nil, err
	}

	item, exists := items[productItemID]
	if !exists {
		return nil, nil
	}

	var productName string
	var shopID uint
	if item.Product != nil {
		productName = item.Product.Name
		shopID = item.Product.ShopID
	}

	return &ProductItemDTO{
		ID:          item.ID,
		SKUCode:     item.SKUCode,
		QtyInStock:  item.QtyInStock,
		ProductName: productName,
		Price:       item.Price,
		ImageURL:    item.ImageURL,
		Status:      item.Status,
		ShopID:      shopID,
	}, nil
}

// GetProductItems fetches multiple product items in batch - for CartService display
// Returns display-only DTOs without validation fields
func (a *CartProductClientAdapter) GetProductItems(productItemIDs []uint) (map[uint]*ProductItemDTO, error) {
	items, err := a.Client.GetProductItems(productItemIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*ProductItemDTO)
	for id, item := range items {
		var productName string
		var shopID uint
		if item.Product != nil {
			productName = item.Product.Name
			shopID = item.Product.ShopID
		}

		result[id] = &ProductItemDTO{
			ID:          item.ID,
			SKUCode:     item.SKUCode,
			QtyInStock:  item.QtyInStock,
			ProductName: productName,
			Price:       item.Price,
			ImageURL:    item.ImageURL,
			Status:      item.Status,
			ShopID:      shopID,
		}
	}

	return result, nil
}

// ==================== OrderProductClientAdapter for OrderService ====================

type OrderProductClientAdapter struct {
	Client *product_client.ProductClient
}

// GetProductItem fetches single product item details (SKU-level) - for OrderService validation
// Returns full DTO with validation fields (Stock, IsActive)
func (a *OrderProductClientAdapter) GetProductItem(productItemID uint) (*OrderProductItemDTO, error) {
	items, err := a.Client.GetProductItems([]uint{productItemID})
	if err != nil {
		return nil, err
	}

	item, exists := items[productItemID]
	if !exists {
		return nil, nil
	}

	var productName string
	var shopID uint
	if item.Product != nil {
		productName = item.Product.Name
		shopID = item.Product.ShopID
	}

	return &OrderProductItemDTO{
		ID:          item.ID,
		ProductID:   item.ProductID,
		ShopID:      shopID,
		ProductName: productName,
		SKU:         item.SKUCode,
		Price:       item.Price,
		Stock:       item.QtyInStock,
		ImageURL:    item.ImageURL,
		IsActive:    item.Status == "active",
	}, nil
}

// GetProductItems fetches multiple product items in batch - for OrderService validation
// Returns full DTOs with validation fields
func (a *OrderProductClientAdapter) GetProductItems(productItemIDs []uint) (map[uint]*OrderProductItemDTO, error) {
	items, err := a.Client.GetProductItems(productItemIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*OrderProductItemDTO)
	for id, item := range items {
		var productName string
		var shopID uint
		if item.Product != nil {
			productName = item.Product.Name
			shopID = item.Product.ShopID
		}

		result[id] = &OrderProductItemDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ShopID:      shopID,
			ProductName: productName,
			SKU:         item.SKUCode,
			Price:       item.Price,
			Stock:       item.QtyInStock,
			ImageURL:    item.ImageURL,
			IsActive:    item.Status == "active",
		}
	}

	return result, nil
}
