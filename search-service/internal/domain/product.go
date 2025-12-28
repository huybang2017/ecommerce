package domain

import (
	"time"
)

// Product represents the core domain entity for search
// This is the business object that exists independently of infrastructure
// Following Clean Architecture: domain layer has no external dependencies
type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SKU         string    `json:"sku"`
	CategoryID  *uint     `json:"category_id,omitempty"`
	Status      string    `json:"status"` // ACTIVE, INACTIVE
	Stock       int       `json:"stock"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductEvent represents a domain event for product changes from Kafka
// Events are used for inter-service communication
type ProductEvent struct {
	EventType   string    `json:"event_type"`   // e.g., "product_created", "product_updated", "product_deleted"
	ProductID   uint      `json:"product_id"`
	ProductData *Product  `json:"product_data"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// SearchFilters represents search filters
type SearchFilters struct {
	CategoryID *uint    `json:"category_id,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	Status     *string  `json:"status,omitempty"`
}

// SearchSort represents sort options
type SearchSort struct {
	Field string `json:"field"` // "price", "name", "created_at"
	Order string `json:"order"` // "asc", "desc"
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query   string         `json:"query"`
	Filters *SearchFilters `json:"filters,omitempty"`
	Sort    *SearchSort    `json:"sort,omitempty"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
}

// SearchResult represents search results with pagination
type SearchResult struct {
	Products []*Product `json:"products"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	Limit    int        `json:"limit"`
}

// SearchRepository defines the interface for search operations
// This is part of the domain layer - it defines WHAT we need, not HOW
type SearchRepository interface {
	IndexProduct(product *Product) error
	UpdateProduct(product *Product) error
	DeleteProduct(id uint) error
	SearchProducts(req *SearchRequest) (*SearchResult, error)
}



