package product_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ProductClient handles communication with Product Service
type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProductClient creates a new product client
func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Product represents product information from Product Service
type Product struct {
	ID     uint `json:"id"`
	ShopID uint `json:"shop_id"` // Required for marketplace
	Name   string `json:"name"`
	BasePrice float64 `json:"base_price"`
}

// ProductItem represents SKU information from Product Service
type ProductItem struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	SKUCode   string  `json:"sku_code"`
	Price     float64 `json:"price"`
	QtyInStock int    `json:"qty_in_stock"`
	Status    string  `json:"status"`
}

// GetProductByID retrieves product information by ID
func (c *ProductClient) GetProductByID(productID uint) (*Product, error) {
	return c.GetProductByIDInternal(productID)
}

// GetProductByIDInternal is the internal implementation
func (c *ProductClient) GetProductByIDInternal(productID uint) (*Product, error) {
	url := fmt.Sprintf("%s/api/v1/products/%d", c.baseURL, productID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product service returned error: %d - %s", resp.StatusCode, string(body))
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %w", err)
	}

	return &product, nil
}

// GetProductItemByID retrieves product item (SKU) information by ID
func (c *ProductClient) GetProductItemByID(productItemID uint) (*ProductItem, error) {
	// Note: Product Service doesn't have direct GET /product-items/:id endpoint
	// We need to get it through product_id first, then find the item
	// For now, return error - will need to implement proper endpoint
	return nil, fmt.Errorf("get product item by ID not yet implemented - need product_id")
}

// GetProductItemByProductID retrieves product items for a product
func (c *ProductClient) GetProductItemByProductID(productID uint) ([]*ProductItem, error) {
	url := fmt.Sprintf("%s/api/v1/products/%d/items", c.baseURL, productID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product service returned error: %d - %s", resp.StatusCode, string(body))
	}

	var response struct {
		Items []*ProductItem `json:"items"`
		Count int            `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode product items response: %w", err)
	}

	return response.Items, nil
}

// GetProductItemBySKUCode retrieves product item by SKU code
func (c *ProductClient) GetProductItemBySKUCode(skuCode string) (*ProductItem, error) {
	url := fmt.Sprintf("%s/api/v1/product-items/%s", c.baseURL, skuCode)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("product service returned error: %d - %s", resp.StatusCode, string(body))
	}

	var item ProductItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("failed to decode product item response: %w", err)
	}

	return &item, nil
}

// GetProductsByIDs retrieves multiple products by IDs (batch)
func (c *ProductClient) GetProductsByIDs(productIDs []uint) (map[uint]*Product, error) {
	products := make(map[uint]*Product)
	
	// For now, fetch sequentially (can optimize with batch endpoint later)
	for _, id := range productIDs {
		product, err := c.GetProductByID(id)
		if err != nil {
			// Log error but continue with other products
			continue
		}
		products[id] = product
	}
	
	return products, nil
}

