package domain

import (
	"time"
)

// User represents the core domain entity for user
// Following Clean Architecture: domain layer has no external dependencies
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"uniqueIndex;size:50" json:"username"`
	Email       string    `gorm:"uniqueIndex;size:100" json:"email"`
	PasswordHash string  `gorm:"column:password_hash;size:255" json:"-"`
	PhoneNumber string   `gorm:"size:20" json:"phone_number"`
	FullName    string    `gorm:"size:100" json:"full_name"`
	AvatarURL   string    `gorm:"column:avatar_url;size:255" json:"avatar_url"`
	Role        string    `gorm:"size:20;default:'BUYER'" json:"role"` // ADMIN, SELLER, BUYER
	Status      string    `gorm:"size:20;default:'ACTIVE'" json:"status"` // ACTIVE, BANNED, DELETED
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "user"
}

// UserRepository defines the interface for user data access
// This is part of the domain layer - it defines WHAT we need, not HOW
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Delete(id uint) error
}

