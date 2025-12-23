package domain

import (
	"encoding/json"
	"time"
)

// ProductEvent represents a domain event for product changes
// Events are used for inter-service communication via Kafka
// Following Domain-Driven Design principles
type ProductEvent struct {
	EventType   string      `json:"event_type"`   // e.g., "product_created", "product_updated"
	ProductID   uint        `json:"product_id"`
	ProductData *Product    `json:"product_data"`
	Timestamp   time.Time   `json:"timestamp"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// ToJSON converts the event to JSON bytes for Kafka publishing
func (e *ProductEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// EventPublisher defines the interface for publishing domain events
// This abstraction allows us to swap Kafka for other message brokers if needed
type EventPublisher interface {
	PublishProductEvent(event *ProductEvent) error
	Close() error // Close releases resources (e.g., Kafka connections)
}

