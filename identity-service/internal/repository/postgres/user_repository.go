package postgres

import (
	"identity-service/internal/domain"

	"gorm.io/gorm"
)

// userRepository implements the UserRepository interface
// This is the infrastructure layer - it knows HOW to interact with PostgreSQL
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository
// Dependency injection: we inject the database connection
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// GetByID retrieves a user by its ID
func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete soft deletes a user (sets status to DELETED)
func (r *userRepository) Delete(id uint) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("status", "DELETED").Error
}


