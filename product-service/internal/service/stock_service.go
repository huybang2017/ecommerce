package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"product-service/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// StockService handles stock management operations
// This service prevents overselling with Redis distributed locks
type StockService struct {
	productItemRepo domain.ProductItemRepository
	redisClient     *redis.Client
	logger          *zap.Logger
}

// NewStockService creates a new stock service
func NewStockService(
	productItemRepo domain.ProductItemRepository,
	redisClient *redis.Client,
	logger *zap.Logger,
) *StockService {
	return &StockService{
		productItemRepo: productItemRepo,
		redisClient:     redisClient,
		logger:          logger,
	}
}

// CheckStock checks if stock is available for given items
func (s *StockService) CheckStock(ctx context.Context, req *domain.StockCheckRequest) (*domain.StockCheckResponse, error) {
	unavailableItems := []domain.UnavailableStockItem{}

	for _, item := range req.Items {
		// Get product item
		productItem, err := s.productItemRepo.GetByID(item.ProductItemID)
		if err != nil {
			s.logger.Error("failed to get product item", zap.Uint("product_item_id", item.ProductItemID), zap.Error(err))
			unavailableItems = append(unavailableItems, domain.UnavailableStockItem{
				ProductItemID: item.ProductItemID,
				Requested:     item.Quantity,
				Available:     0,
			})
			continue
		}

		// Check if enough stock
		if productItem.QtyInStock < item.Quantity {
			unavailableItems = append(unavailableItems, domain.UnavailableStockItem{
				ProductItemID: item.ProductItemID,
				Requested:     item.Quantity,
				Available:     productItem.QtyInStock,
			})
		}
	}

	return &domain.StockCheckResponse{
		Available:        len(unavailableItems) == 0,
		UnavailableItems: unavailableItems,
	}, nil
}

// ReserveStock temporarily reserves stock for an order (stores in Redis)
// This prevents overselling during checkout flow
func (s *StockService) ReserveStock(ctx context.Context, req *domain.StockReserveRequest) error {
	// Validate order_id
	if req.OrderID == "" {
		return errors.New("order_id is required")
	}

	// Check stock availability first
	checkReq := &domain.StockCheckRequest{Items: []domain.StockCheckItem{}}
	for _, item := range req.Items {
		checkReq.Items = append(checkReq.Items, domain.StockCheckItem{
			ProductItemID: item.ProductItemID,
			Quantity:      item.Quantity,
		})
	}

	checkResp, err := s.CheckStock(ctx, checkReq)
	if err != nil {
		return fmt.Errorf("failed to check stock: %w", err)
	}
	if !checkResp.Available {
		return fmt.Errorf("insufficient stock: %v", checkResp.UnavailableItems)
	}

	// Reserve each item in Redis (with TTL = 15 minutes)
	expiresAt := time.Now().Add(15 * time.Minute)
	for _, item := range req.Items {
		reservation := &domain.StockReservation{
			OrderID:       req.OrderID,
			ProductItemID: item.ProductItemID,
			Quantity:      item.Quantity,
			ExpiresAt:     expiresAt,
		}

		// Store in Redis
		key := fmt.Sprintf("stock:reservation:%s:%d", req.OrderID, item.ProductItemID)
		data, err := json.Marshal(reservation)
		if err != nil {
			s.logger.Error("failed to marshal reservation", zap.Error(err))
			continue
		}

		if err := s.redisClient.Set(ctx, key, data, 15*time.Minute).Err(); err != nil {
			s.logger.Error("failed to store reservation", zap.String("key", key), zap.Error(err))
			return fmt.Errorf("failed to reserve stock: %w", err)
		}

		s.logger.Info("stock reserved",
			zap.String("order_id", req.OrderID),
			zap.Uint("product_item_id", item.ProductItemID),
			zap.Int("quantity", item.Quantity),
		)
	}

	return nil
}

// DeductStock permanently deducts stock from product_item.qty_in_stock
// This should be called after payment is confirmed
func (s *StockService) DeductStock(ctx context.Context, req *domain.StockDeductRequest) error {
	// Validate order_id
	if req.OrderID == "" {
		return errors.New("order_id is required")
	}

	// Deduct each item with distributed lock
	for _, item := range req.Items {
		if err := s.deductStockWithLock(ctx, item.ProductItemID, item.Quantity); err != nil {
			s.logger.Error("failed to deduct stock",
				zap.Uint("product_item_id", item.ProductItemID),
				zap.Int("quantity", item.Quantity),
				zap.Error(err),
			)
			return fmt.Errorf("failed to deduct stock for product_item %d: %w", item.ProductItemID, err)
		}
	}

	// Release reservation from Redis after successful deduction
	if err := s.ReleaseStock(ctx, &domain.StockReleaseRequest{OrderID: req.OrderID}); err != nil {
		s.logger.Warn("failed to release reservation after deduction", zap.String("order_id", req.OrderID), zap.Error(err))
		// Continue even if release fails (deduction already succeeded)
	}

	return nil
}

// deductStockWithLock deducts stock with Redis distributed lock to prevent race condition
func (s *StockService) deductStockWithLock(ctx context.Context, productItemID uint, quantity int) error {
	lockKey := fmt.Sprintf("stock:lock:%d", productItemID)
	lockValue := fmt.Sprintf("%d-%d", time.Now().UnixNano(), productItemID)
	lockTTL := 30 * time.Second

	// Acquire lock with retry (max 3 attempts)
	var locked bool
	for i := 0; i < 3; i++ {
		locked, err := s.redisClient.SetNX(ctx, lockKey, lockValue, lockTTL).Result()
		if err != nil {
			s.logger.Error("failed to acquire lock", zap.String("key", lockKey), zap.Error(err))
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if locked {
			break
		}
		// Lock already held by another process, wait and retry
		time.Sleep(100 * time.Millisecond)
	}

	if !locked {
		return errors.New("failed to acquire lock after retries")
	}

	// Ensure lock is released even if error occurs
	defer func() {
		// Release lock
		if err := s.redisClient.Del(ctx, lockKey).Err(); err != nil {
			s.logger.Warn("failed to release lock", zap.String("key", lockKey), zap.Error(err))
		}
	}()

	// Get current stock
	productItem, err := s.productItemRepo.GetByID(productItemID)
	if err != nil {
		return fmt.Errorf("product item not found: %w", err)
	}

	// Check if enough stock
	if productItem.QtyInStock < quantity {
		return fmt.Errorf("insufficient stock: requested %d, available %d", quantity, productItem.QtyInStock)
	}

	// Deduct stock (atomic operation)
	newStock := productItem.QtyInStock - quantity
	if err := s.productItemRepo.UpdateStock(productItemID, newStock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Update status if out of stock
	if newStock == 0 {
		productItem.Status = "OUT_OF_STOCK"
		if err := s.productItemRepo.Update(productItem); err != nil {
			s.logger.Warn("failed to update status to OUT_OF_STOCK", zap.Uint("product_item_id", productItemID), zap.Error(err))
		}
	}

	s.logger.Info("stock deducted",
		zap.Uint("product_item_id", productItemID),
		zap.Int("quantity", quantity),
		zap.Int("new_stock", newStock),
	)

	return nil
}

// ReleaseStock releases reserved stock from Redis
// This should be called when order is cancelled or payment failed
func (s *StockService) ReleaseStock(ctx context.Context, req *domain.StockReleaseRequest) error {
	// Validate order_id
	if req.OrderID == "" {
		return errors.New("order_id is required")
	}

	// Find and delete all reservations for this order
	pattern := fmt.Sprintf("stock:reservation:%s:*", req.OrderID)
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		s.logger.Error("failed to find reservations", zap.String("order_id", req.OrderID), zap.Error(err))
		return fmt.Errorf("failed to find reservations: %w", err)
	}

	if len(keys) == 0 {
		s.logger.Warn("no reservations found for order", zap.String("order_id", req.OrderID))
		return nil // No reservations to release
	}

	// Delete all reservation keys
	if err := s.redisClient.Del(ctx, keys...).Err(); err != nil {
		s.logger.Error("failed to delete reservations", zap.String("order_id", req.OrderID), zap.Error(err))
		return fmt.Errorf("failed to release reservations: %w", err)
	}

	s.logger.Info("stock reservations released",
		zap.String("order_id", req.OrderID),
		zap.Int("count", len(keys)),
	)

	return nil
}

// GetStock retrieves current stock for a product item
func (s *StockService) GetStock(ctx context.Context, productItemID uint) (int, error) {
	productItem, err := s.productItemRepo.GetByID(productItemID)
	if err != nil {
		return 0, fmt.Errorf("product item not found: %w", err)
	}

	return productItem.QtyInStock, nil
}

// UpdateStock updates the stock quantity for a product item
// This is for shop owners to update their stock
func (s *StockService) UpdateStock(ctx context.Context, productItemID uint, newStock int) error {
	if newStock < 0 {
		return errors.New("stock cannot be negative")
	}

	productItem, err := s.productItemRepo.GetByID(productItemID)
	if err != nil {
		return fmt.Errorf("product item not found: %w", err)
	}

	// Update stock with lock
	lockKey := fmt.Sprintf("stock:lock:%d", productItemID)
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
	locked, err := s.redisClient.SetNX(ctx, lockKey, lockValue, 10*time.Second).Result()
	if err != nil || !locked {
		return errors.New("failed to acquire lock for stock update")
	}
	defer s.redisClient.Del(ctx, lockKey)

	// Update stock
	if err := s.productItemRepo.UpdateStock(productItemID, newStock); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// Update status based on stock
	if newStock == 0 && productItem.Status != "OUT_OF_STOCK" {
		productItem.Status = "OUT_OF_STOCK"
		if err := s.productItemRepo.Update(productItem); err != nil {
			s.logger.Warn("failed to update status", zap.Error(err))
		}
	} else if newStock > 0 && productItem.Status == "OUT_OF_STOCK" {
		productItem.Status = "ACTIVE"
		if err := s.productItemRepo.Update(productItem); err != nil {
			s.logger.Warn("failed to update status", zap.Error(err))
		}
	}

	s.logger.Info("stock updated",
		zap.Uint("product_item_id", productItemID),
		zap.Int("new_stock", newStock),
	)

	return nil
}

