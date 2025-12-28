package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"search-service/internal/domain"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// searchRepository implements the SearchRepository interface
// This is the infrastructure layer - it knows HOW to interact with Elasticsearch
type searchRepository struct {
	client    *elasticsearch.Client
	indexName string
}

// NewSearchRepository creates a new Elasticsearch search repository
// Dependency injection: we inject the Elasticsearch client
func NewSearchRepository(client *elasticsearch.Client, indexName string) domain.SearchRepository {
	return &searchRepository{
		client:    client,
		indexName: indexName,
	}
}

// IndexProduct indexes a product document in Elasticsearch
func (r *searchRepository) IndexProduct(product *domain.Product) error {
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

// UpdateProduct updates a product document in Elasticsearch (same as IndexProduct)
func (r *searchRepository) UpdateProduct(product *domain.Product) error {
	return r.IndexProduct(product)
}

// DeleteProduct removes a product from the Elasticsearch index
func (r *searchRepository) DeleteProduct(id uint) error {
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

// SearchProducts performs a search query with filters, sort, and pagination
func (r *searchRepository) SearchProducts(req *domain.SearchRequest) (*domain.SearchResult, error) {
	ctx := context.Background()

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	// Build the search query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []map[string]interface{}{},
				"filter": []map[string]interface{}{},
			},
		},
		"from": (req.Page - 1) * req.Limit,
		"size": req.Limit,
	}

	boolQuery := query["query"].(map[string]interface{})["bool"].(map[string]interface{})
	mustClauses := boolQuery["must"].([]map[string]interface{})
	filterClauses := boolQuery["filter"].([]map[string]interface{})

	// Add text search if query is provided
	if strings.TrimSpace(req.Query) != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^3", "description^2", "sku"},
				"type":   "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// Add filters
	if req.Filters != nil {
		if req.Filters.CategoryID != nil {
			filterClauses = append(filterClauses, map[string]interface{}{
				"term": map[string]interface{}{
					"category_id": *req.Filters.CategoryID,
				},
			})
		}

		if req.Filters.MinPrice != nil || req.Filters.MaxPrice != nil {
			rangeQuery := map[string]interface{}{}
			if req.Filters.MinPrice != nil {
				rangeQuery["gte"] = *req.Filters.MinPrice
			}
			if req.Filters.MaxPrice != nil {
				rangeQuery["lte"] = *req.Filters.MaxPrice
			}
			filterClauses = append(filterClauses, map[string]interface{}{
				"range": map[string]interface{}{
					"price": rangeQuery,
				},
			})
		}

		if req.Filters.Status != nil {
			filterClauses = append(filterClauses, map[string]interface{}{
				"term": map[string]interface{}{
					"status": *req.Filters.Status,
				},
			})
		}
	}

	// Update clauses
	boolQuery["must"] = mustClauses
	boolQuery["filter"] = filterClauses

	// Add sort
	if req.Sort != nil {
		sortField := req.Sort.Field
		if sortField == "" {
			sortField = "_score" // Default to relevance
		}

		sortOrder := "asc"
		if req.Sort.Order == "desc" {
			sortOrder = "desc"
		}

		query["sort"] = []map[string]interface{}{
			{
				sortField: map[string]interface{}{
					"order": sortOrder,
				},
			},
		}
	} else {
		// Default sort by relevance
		query["sort"] = []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
		}

		// If no query, sort by created_at desc
		if strings.TrimSpace(req.Query) == "" {
			query["sort"] = []map[string]interface{}{
				{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			}
		}
	}

	// Convert to JSON
	queryJSON, err := json.Marshal(query)
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

	// Extract total
	total := int64(0)
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if totalValue, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := totalValue["value"].(float64); ok {
				total = int64(value)
			}
		} else if totalValue, ok := hits["total"].(float64); ok {
			total = int64(totalValue)
		}
	}

	// Extract products from hits
	products := make([]*domain.Product, 0)
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				hitMap := hit.(map[string]interface{})
				source := hitMap["_source"].(map[string]interface{})

				// Convert to Product struct
				productJSON, _ := json.Marshal(source)
				var product domain.Product
				if err := json.Unmarshal(productJSON, &product); err == nil {
					products = append(products, &product)
				}
			}
		}
	}

	return &domain.SearchResult{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

