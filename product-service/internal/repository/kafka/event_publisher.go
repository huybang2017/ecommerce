package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"product-service/internal/domain"
	"time"

	"github.com/segmentio/kafka-go"
)

// eventPublisher implements the EventPublisher interface
// This is the infrastructure layer - it knows HOW to publish events to Kafka
type eventPublisher struct {
	writer *kafka.Writer
	topic  string
}

// NewEventPublisher creates a new Kafka event publisher
// Dependency injection: we inject the Kafka writer
func NewEventPublisher(brokers []string, topic string, writeTimeout time.Duration, requiredAcks int) domain.EventPublisher {
	// Convert int to kafka.RequiredAcks
	var kafkaAcks kafka.RequiredAcks
	switch requiredAcks {
	case -1:
		kafkaAcks = kafka.RequireAll
	case 0:
		kafkaAcks = kafka.RequireNone
	case 1:
		kafkaAcks = kafka.RequireOne
	default:
		kafkaAcks = kafka.RequireOne
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: writeTimeout,
		RequiredAcks: kafkaAcks,
		Async:        false, // Synchronous writes for reliability
	}

	return &eventPublisher{
		writer: writer,
		topic:  topic,
	}
}

// PublishProductEvent publishes a product event to Kafka
// This enables event-driven architecture and inter-service communication
func (p *eventPublisher) PublishProductEvent(event *domain.ProductEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", event.ProductID)),
		Value: eventJSON,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.EventType)},
			{Key: "timestamp", Value: []byte(event.Timestamp.Format(time.RFC3339))},
		},
	}

	// Write message to Kafka
	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

// Close closes the Kafka writer connection
// This should be called during graceful shutdown
func (p *eventPublisher) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

