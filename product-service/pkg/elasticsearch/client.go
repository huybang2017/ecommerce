package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"product-service/config"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	// clientInstance is the singleton Elasticsearch client
	clientInstance *elasticsearch.Client
	// once ensures the client is created only once
	once sync.Once
)

// GetClient returns the singleton Elasticsearch client
// This implements the Singleton pattern to ensure only one ES connection pool exists
func GetClient(cfg *config.ElasticsearchConfig) (*elasticsearch.Client, error) {
	var err error

	once.Do(func() {
		// Configure Elasticsearch client
		esConfig := elasticsearch.Config{
			Addresses: cfg.Addresses,
			Username:  cfg.Username,
			Password:  cfg.Password,
		}

		clientInstance, err = elasticsearch.NewClient(esConfig)
		if err != nil {
			log.Printf("Failed to create Elasticsearch client: %v", err)
			return
		}

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		res, err := clientInstance.Info(clientInstance.Info.WithContext(ctx))
		if err != nil {
			log.Printf("Failed to connect to Elasticsearch: %v", err)
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			err = fmt.Errorf("elasticsearch error: %s", res.String())
			log.Printf("Elasticsearch connection error: %v", err)
			return
		}

		log.Println("Elasticsearch connection established successfully")
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Elasticsearch client: %w", err)
	}

	return clientInstance, nil
}

// EnsureIndex creates the Elasticsearch index if it doesn't exist
// This should be called at application startup
func EnsureIndex(client *elasticsearch.Client, indexName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if index exists
	exists, err := client.Indices.Exists([]string{indexName}, client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}
	defer exists.Body.Close()

	if exists.StatusCode == 200 {
		log.Printf("Index '%s' already exists", indexName)
		return nil
	}

	// Create index with mapping
	// In production, you'd want a more sophisticated mapping
	mapping := `{
		"mappings": {
			"properties": {
				"id": { "type": "long" },
				"name": { "type": "text", "analyzer": "standard" },
				"description": { "type": "text", "analyzer": "standard" },
				"price": { "type": "float" },
				"sku": { "type": "keyword" },
				"category": { "type": "keyword" },
				"stock": { "type": "integer" },
				"is_active": { "type": "boolean" },
				"created_at": { "type": "date" },
				"updated_at": { "type": "date" }
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error creating index: %s", res.String())
	}

	log.Printf("Index '%s' created successfully", indexName)
	return nil
}

