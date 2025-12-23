package postgres

import (
	"identity-service/internal/domain"

	"gorm.io/gorm"
)

// addressRepository implements the AddressRepository interface
// This is the infrastructure layer - it knows HOW to interact with PostgreSQL
type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new PostgreSQL address repository
// Dependency injection: we inject the database connection
func NewAddressRepository(db *gorm.DB) domain.AddressRepository {
	return &addressRepository{db: db}
}

// Create inserts a new address into the database
func (r *addressRepository) Create(address *domain.Address) error {
	return r.db.Create(address).Error
}

// Update updates an existing address
func (r *addressRepository) Update(address *domain.Address) error {
	return r.db.Save(address).Error
}

// GetByID retrieves an address by its ID
func (r *addressRepository) GetByID(id uint) (*domain.Address, error) {
	var address domain.Address
	err := r.db.First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *addressRepository) GetByUserID(userID uint) ([]*domain.Address, error) {
	var addresses []*domain.Address
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetDefaultByUserID retrieves the default address for a user
func (r *addressRepository) GetDefaultByUserID(userID uint) (*domain.Address, error) {
	var address domain.Address
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Delete deletes an address
func (r *addressRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Address{}, id).Error
}

// SetDefault sets an address as default and unsets others for the same user
func (r *addressRepository) SetDefault(userID uint, addressID uint) error {
	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Unset all default addresses for this user
	if err := tx.Model(&domain.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set the specified address as default
	if err := tx.Model(&domain.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}


