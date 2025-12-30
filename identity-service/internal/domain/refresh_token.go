package domain

import (
	"time"
)

// RefreshToken represents a refresh token for maintaining user sessions
// Used for implementing secure token refresh mechanism
type RefreshToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Token     string     `gorm:"uniqueIndex;size:500;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	IsRevoked bool       `gorm:"default:false" json:"is_revoked"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`

	// Relationship
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_token"
}

// IsValid checks if the refresh token is still valid
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && time.Now().Before(rt.ExpiresAt)
}

// Revoke marks the token as revoked (used for logout/blacklist)
func (rt *RefreshToken) Revoke() {
	rt.IsRevoked = true
	now := time.Now()
	rt.RevokedAt = &now
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	GetByUserID(userID uint) ([]*RefreshToken, error)
	Update(token *RefreshToken) error
	Delete(id uint) error
	RevokeAllByUserID(userID uint) error
	CleanupExpired() error
}
