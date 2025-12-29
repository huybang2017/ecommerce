package domain

import (
	"encoding/json"
	"time"
)

// OrderEvent represents a domain event for order changes
// Events are used for inter-service communication via Kafka
// Following Domain-Driven Design principles
type OrderEvent struct {
	EventType   string      `json:"event_type"`   // e.g., "order_created", "order_updated"
	OrderID     uint        `json:"order_id"`
	OrderData   *Order      `json:"order_data"`
	Timestamp   time.Time   `json:"timestamp"`
	Metadata    interface{} `json:"metadata,omitempty"`
}

// ToJSON converts the event to JSON bytes for Kafka publishing
func (e *OrderEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// EventPublisher defines the interface for publishing domain events
// This abstraction allows us to swap Kafka for other message brokers if needed
type OrderEventPublisher interface {
	PublishOrderEvent(event *OrderEvent) error
	Close() error // Close releases resources (e.g., Kafka connections)
}

