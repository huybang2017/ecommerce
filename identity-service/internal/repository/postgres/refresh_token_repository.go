package postgres

import (
	"identity-service/internal/domain"
	"time"

	"gorm.io/gorm"
)

// refreshTokenRepository implements the RefreshTokenRepository interface
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new PostgreSQL refresh token repository
func NewRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create inserts a new refresh token into the database
func (r *refreshTokenRepository) Create(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetByToken retrieves a refresh token by its token string
func (r *refreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *refreshTokenRepository) GetByUserID(userID uint) ([]*domain.RefreshToken, error) {
	var tokens []*domain.RefreshToken
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Update updates an existing refresh token
func (r *refreshTokenRepository) Update(token *domain.RefreshToken) error {
	return r.db.Save(token).Error
}

// Delete deletes a refresh token by ID
func (r *refreshTokenRepository) Delete(id uint) error {
	return r.db.Delete(&domain.RefreshToken{}, id).Error
}

// RevokeAllByUserID revokes all tokens for a user (used during logout)
func (r *refreshTokenRepository) RevokeAllByUserID(userID uint) error {
	now := time.Now()
	return r.db.Model(&domain.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		}).Error
}

// CleanupExpired removes expired tokens (can be called periodically)
func (r *refreshTokenRepository) CleanupExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&domain.RefreshToken{}).Error
}
