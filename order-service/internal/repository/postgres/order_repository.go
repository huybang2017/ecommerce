package postgres

import (
	"order-service/internal/domain"

	"gorm.io/gorm"
)

// OrderRepository handles database operations for orders
// This is the infrastructure layer - it knows HOW to persist data
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order in the database
func (r *OrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetByOrderNumber(orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves all orders for a user
func (r *OrderRepository) GetByUserID(userID uint, limit, offset int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	// Count total
	if err := r.db.Model(&domain.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	err := r.db.Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetBySessionID retrieves all orders for a session (guest orders)
func (r *OrderRepository) GetBySessionID(sessionID string, limit, offset int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	// Count total
	if err := r.db.Model(&domain.Order{}).Where("session_id = ?", sessionID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	err := r.db.Preload("Items").
		Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateStatus updates the status of an order
func (r *OrderRepository) UpdateStatus(orderID uint, status domain.OrderStatus) error {
	return r.db.Model(&domain.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

