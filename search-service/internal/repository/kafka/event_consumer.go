package kafka

import (
	"context"
	"encoding/json"
	"log"
	"search-service/internal/domain"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// EventConsumer handles consuming product events from Kafka
// This is the infrastructure layer - it knows HOW to consume from Kafka
type EventConsumer struct {
	reader      *kafka.Reader
	searchRepo  domain.SearchRepository
	logger      *zap.Logger
}

// NewEventConsumer creates a new Kafka event consumer
func NewEventConsumer(
	brokers []string,
	topic string,
	consumerGroup string,
	readTimeout time.Duration,
	minBytes int,
	maxBytes int,
	searchRepo domain.SearchRepository,
	logger *zap.Logger,
) *EventConsumer {
	// Validate inputs
	if len(brokers) == 0 {
		logger.Error("Kafka brokers list is empty")
		panic("Kafka brokers list is empty")
	}
	if topic == "" {
		logger.Error("Kafka topic is empty")
		panic("Kafka topic is empty")
	}
	if consumerGroup == "" {
		logger.Error("Kafka consumer group is empty")
		panic("Kafka consumer group is empty")
	}

	logger.Info("Creating Kafka reader",
		zap.Strings("brokers", brokers),
		zap.String("topic", topic),
		zap.String("consumer_group", consumerGroup),
	)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        consumerGroup,
		MinBytes:       minBytes,
		MaxBytes:       maxBytes,
		ReadBackoffMin: 100 * time.Millisecond,
		ReadBackoffMax: 1 * time.Second,
	})

	logger.Info("Kafka reader created successfully")

	return &EventConsumer{
		reader:     reader,
		searchRepo: searchRepo,
		logger:     logger,
	}
}

// Start starts consuming messages from Kafka
// This runs in a goroutine and processes events asynchronously
func (c *EventConsumer) Start(ctx context.Context) error {
	// Use both logger and log for maximum visibility
	log.Printf("🚀🚀🚀 Kafka consumer Start() method called! 🚀🚀🚀\n")
	c.logger.Info("🚀 Starting Kafka consumer",
		zap.String("topic", c.reader.Config().Topic),
		zap.String("consumer_group", c.reader.Config().GroupID),
		zap.Strings("brokers", c.reader.Config().Brokers),
	)

	log.Printf("✅ Kafka consumer entering main loop - ready to receive messages\n")
	c.logger.Info("✅ Kafka consumer entering main loop - ready to receive messages")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping Kafka consumer")
			return ctx.Err()
		default:
			// Read message with timeout
			msgCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			message, err := c.reader.ReadMessage(msgCtx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded || err == context.Canceled {
					// Timeout is normal when no messages - continue waiting
					// Log every 10th timeout to show consumer is alive
					if time.Now().Unix()%10 == 0 {
						c.logger.Debug("Waiting for messages... (timeout is normal)", zap.Error(err))
					}
					continue
				}
				c.logger.Error("❌ Failed to read message from Kafka", zap.Error(err))
				time.Sleep(1 * time.Second) // Backoff on error
				continue
			}

			c.logger.Info("🎉 RECEIVED MESSAGE FROM KAFKA!",
				zap.String("topic", message.Topic),
				zap.Int("partition", message.Partition),
				zap.Int64("offset", message.Offset),
				zap.Int("message_size", len(message.Value)),
			)

			// Process message in goroutine to avoid blocking
			go c.processMessage(message)
		}
	}
}

// processMessage processes a single Kafka message
func (c *EventConsumer) processMessage(message kafka.Message) {
	c.logger.Debug("Received message",
		zap.String("topic", message.Topic),
		zap.Int("partition", message.Partition),
		zap.Int64("offset", message.Offset),
	)

	// Parse event
	var event domain.ProductEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		c.logger.Error("Failed to unmarshal event", zap.Error(err))
		return
	}

	// Handle event based on type
	switch event.EventType {
	case "product_created", "product_updated":
		if event.ProductData == nil {
			c.logger.Warn("Product data is nil in event", zap.String("event_type", event.EventType))
			return
		}

		// Index or update product in Elasticsearch
		log.Printf("📤 Indexing product to Elasticsearch: ID=%d, Name=%s\n", event.ProductID, event.ProductData.Name)
		if err := c.searchRepo.IndexProduct(event.ProductData); err != nil {
			log.Printf("❌ Failed to index product: %v\n", err)
			c.logger.Error("Failed to index product",
				zap.Uint("product_id", event.ProductID),
				zap.String("event_type", event.EventType),
				zap.Error(err),
			)
			return
		}

		log.Printf("✅✅✅ Product indexed successfully: ID=%d, Name=%s\n", event.ProductID, event.ProductData.Name)
		c.logger.Info("Product indexed successfully",
			zap.Uint("product_id", event.ProductID),
			zap.String("event_type", event.EventType),
		)

	case "product_deleted":
		// Delete product from Elasticsearch
		if err := c.searchRepo.DeleteProduct(event.ProductID); err != nil {
			c.logger.Error("Failed to delete product from index",
				zap.Uint("product_id", event.ProductID),
				zap.Error(err),
			)
			return
		}

		c.logger.Info("Product deleted from index",
			zap.Uint("product_id", event.ProductID),
		)

	default:
		c.logger.Warn("Unknown event type", zap.String("event_type", event.EventType))
	}
}

// Close closes the Kafka reader connection
func (c *EventConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

