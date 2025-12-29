package postgres

import (
	"identity-service/internal/domain"

	"gorm.io/gorm"
)

// shopRepository implements the ShopRepository interface
// This is the infrastructure layer - it knows HOW to interact with PostgreSQL
type shopRepository struct {
	db *gorm.DB
}

// NewShopRepository creates a new PostgreSQL shop repository
// Dependency injection: we inject the database connection
func NewShopRepository(db *gorm.DB) domain.ShopRepository {
	return &shopRepository{db: db}
}

// Create inserts a new shop into the database
func (r *shopRepository) Create(shop *domain.Shop) error {
	return r.db.Create(shop).Error
}

// Update updates an existing shop
func (r *shopRepository) Update(shop *domain.Shop) error {
	return r.db.Save(shop).Error
}

// GetByID retrieves a shop by its ID
func (r *shopRepository) GetByID(id uint) (*domain.Shop, error) {
	var shop domain.Shop
	err := r.db.First(&shop, id).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

// GetByOwnerUserID retrieves a shop by owner user ID (1 User = 1 Shop)
func (r *shopRepository) GetByOwnerUserID(ownerUserID uint) (*domain.Shop, error) {
	var shop domain.Shop
	err := r.db.Where("owner_user_id = ?", ownerUserID).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

// GetAll retrieves all shops with pagination
func (r *shopRepository) GetAll(page, limit int) ([]*domain.Shop, int64, error) {
	var shops []*domain.Shop
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&domain.Shop{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&shops).Error; err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}

// GetByStatus retrieves shops by status with pagination
func (r *shopRepository) GetByStatus(status string, page, limit int) ([]*domain.Shop, int64, error) {
	var shops []*domain.Shop
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&domain.Shop{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Where("status = ?", status).Offset(offset).Limit(limit).Find(&shops).Error; err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}

// Delete soft deletes a shop (sets status to SUSPENDED)
func (r *shopRepository) Delete(id uint) error {
	return r.db.Model(&domain.Shop{}).Where("id = ?", id).Update("status", "SUSPENDED").Error
}

// UpdateStatus updates the status of a shop
func (r *shopRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.Shop{}).Where("id = ?", id).Update("status", status).Error
}

