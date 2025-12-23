package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"product-service/internal/domain"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// productSearchRepository implements the ProductSearchRepository interface
// This is the infrastructure layer - it knows HOW to interact with Elasticsearch
type productSearchRepository struct {
	client    *elasticsearch.Client
	indexName string
}

// NewProductSearchRepository creates a new Elasticsearch product search repository
// Dependency injection: we inject the Elasticsearch client
func NewProductSearchRepository(client *elasticsearch.Client, indexName string) domain.ProductSearchRepository {
	return &productSearchRepository{
		client:    client,
		indexName: indexName,
	}
}

// IndexProduct indexes a product document in Elasticsearch
// This enables fast full-text search and filtering
func (r *productSearchRepository) IndexProduct(product *domain.Product) error {
	ctx := context.Background()

	// Convert product to JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	// Create index request
	req := esapi.IndexRequest{
		Index:      r.indexName,
		DocumentID: fmt.Sprintf("%d", product.ID),
		Body:       bytes.NewReader(productJSON),
		Refresh:    "true", // Make the document immediately searchable
	}

	// Execute the request
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("failed to index product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

// SearchProducts performs a search query with filters
// This is a simplified implementation - in production, you'd want more sophisticated queries
func (r *productSearchRepository) SearchProducts(query string, filters map[string]interface{}) ([]*domain.Product, error) {
	ctx := context.Background()

	// Build the search query
	// In production, you'd use a more sophisticated query builder
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
	}

	// Add text search if query is provided
	if query != "" {
		searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  query,
					"fields": []string{"name^2", "description", "category"},
					"type":   "best_fields",
				},
			},
		)
	}

	// Add filters
	if len(filters) > 0 {
		for key, value := range filters {
			searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
				searchQuery["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{}),
				map[string]interface{}{
					"term": map[string]interface{}{
						key: value,
					},
				},
			)
		}
	}

	// Convert to JSON
	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	// Execute search
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.indexName),
		r.client.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Extract products from hits
	products := make([]*domain.Product, 0)
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		// Convert to Product struct
		productJSON, _ := json.Marshal(source)
		var product domain.Product
		if err := json.Unmarshal(productJSON, &product); err == nil {
			products = append(products, &product)
		}
	}

	return products, nil
}

// DeleteFromIndex removes a product from the Elasticsearch index
func (r *productSearchRepository) DeleteFromIndex(id uint) error {
	ctx := context.Background()

	req := esapi.DeleteRequest{
		Index:      r.indexName,
		DocumentID: fmt.Sprintf("%d", id),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("failed to delete from index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

