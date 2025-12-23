package service

import (
	"errors"
	"fmt"
	"identity-service/internal/domain"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService contains the business logic for user operations
type UserService struct {
	userRepo domain.UserRepository
	logger   *zap.Logger
}

// NewUserService creates a new user service
func NewUserService(
	userRepo domain.UserRepository,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetProfile retrieves a user's profile by ID
func (s *UserService) GetProfile(userID uint) (*domain.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

// UpdateProfile updates a user's profile
type UpdateProfileRequest struct {
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	AvatarURL   string `json:"avatar_url"`
}

func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*domain.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// Save updates
	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error("failed to update user profile", zap.Error(err))
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	s.logger.Info("user profile updated", zap.Uint("user_id", userID))

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

// ChangePassword changes a user's password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error("failed to update password", zap.Error(err))
		return fmt.Errorf("failed to update password: %w", err)
	}

	s.logger.Info("password changed", zap.Uint("user_id", userID))
	return nil
}


