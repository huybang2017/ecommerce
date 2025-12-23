package service

import (
	"errors"
	"fmt"
	"identity-service/internal/domain"

	"go.uber.org/zap"
)

// AddressService contains the business logic for address operations
type AddressService struct {
	addressRepo domain.AddressRepository
	logger      *zap.Logger
}

// NewAddressService creates a new address service
func NewAddressService(
	addressRepo domain.AddressRepository,
	logger *zap.Logger,
) *AddressService {
	return &AddressService{
		addressRepo: addressRepo,
		logger:      logger,
	}
}

// CreateAddressRequest represents the request to create an address
type CreateAddressRequest struct {
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	AddressLine   string `json:"address_line" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district" binding:"required"`
	Ward          string `json:"ward"`
	IsDefault     bool   `json:"is_default"`
	Label         string `json:"label"`
}

// CreateAddress creates a new address for a user
func (s *AddressService) CreateAddress(userID uint, req *CreateAddressRequest) (*domain.Address, error) {
	address := &domain.Address{
		UserID:        userID,
		RecipientName: req.RecipientName,
		PhoneNumber:   req.PhoneNumber,
		AddressLine:   req.AddressLine,
		City:          req.City,
		District:      req.District,
		Ward:          req.Ward,
		IsDefault:     req.IsDefault,
		Label:         req.Label,
	}

	// If this is set as default, unset other defaults
	if req.IsDefault {
		if err := s.addressRepo.SetDefault(userID, 0); err != nil {
			// If no default exists, this is fine - continue
			s.logger.Debug("no existing default address to unset")
		}
	}

	if err := s.addressRepo.Create(address); err != nil {
		s.logger.Error("failed to create address", zap.Error(err))
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	// If this is set as default, update it
	if req.IsDefault {
		if err := s.addressRepo.SetDefault(userID, address.ID); err != nil {
			s.logger.Warn("failed to set address as default", zap.Error(err))
		}
	}

	s.logger.Info("address created", zap.Uint("address_id", address.ID), zap.Uint("user_id", userID))
	return address, nil
}

// UpdateAddressRequest represents the request to update an address
type UpdateAddressRequest struct {
	RecipientName string `json:"recipient_name"`
	PhoneNumber   string `json:"phone_number"`
	AddressLine   string `json:"address_line"`
	City          string `json:"city"`
	District      string `json:"district"`
	Ward          string `json:"ward"`
	IsDefault     *bool  `json:"is_default"`
	Label         string `json:"label"`
}

// UpdateAddress updates an existing address
func (s *AddressService) UpdateAddress(userID uint, addressID uint, req *UpdateAddressRequest) (*domain.Address, error) {
	// Get address
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	if req.RecipientName != "" {
		address.RecipientName = req.RecipientName
	}
	if req.PhoneNumber != "" {
		address.PhoneNumber = req.PhoneNumber
	}
	if req.AddressLine != "" {
		address.AddressLine = req.AddressLine
	}
	if req.City != "" {
		address.City = req.City
	}
	if req.District != "" {
		address.District = req.District
	}
	if req.Ward != "" {
		address.Ward = req.Ward
	}
	if req.Label != "" {
		address.Label = req.Label
	}

	// Handle is_default
	if req.IsDefault != nil && *req.IsDefault {
		if err := s.addressRepo.SetDefault(userID, addressID); err != nil {
			return nil, fmt.Errorf("failed to set default: %w", err)
		}
		address.IsDefault = true
	}

	// Save updates
	if err := s.addressRepo.Update(address); err != nil {
		s.logger.Error("failed to update address", zap.Error(err))
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	s.logger.Info("address updated", zap.Uint("address_id", addressID), zap.Uint("user_id", userID))
	return address, nil
}

// GetAddresses retrieves all addresses for a user
func (s *AddressService) GetAddresses(userID uint) ([]*domain.Address, error) {
	addresses, err := s.addressRepo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("failed to get addresses", zap.Error(err))
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	return addresses, nil
}

// GetAddress retrieves a specific address
func (s *AddressService) GetAddress(userID uint, addressID uint) (*domain.Address, error) {
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	// Verify ownership
	if address.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return address, nil
}

// DeleteAddress deletes an address
func (s *AddressService) DeleteAddress(userID uint, addressID uint) error {
	// Get address to verify ownership
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return errors.New("address not found")
	}

	// Verify ownership
	if address.UserID != userID {
		return errors.New("unauthorized")
	}

	// Delete address
	if err := s.addressRepo.Delete(addressID); err != nil {
		s.logger.Error("failed to delete address", zap.Error(err))
		return fmt.Errorf("failed to delete address: %w", err)
	}

	s.logger.Info("address deleted", zap.Uint("address_id", addressID), zap.Uint("user_id", userID))
	return nil
}

// SetDefaultAddress sets an address as default
func (s *AddressService) SetDefaultAddress(userID uint, addressID uint) error {
	// Verify ownership
	address, err := s.addressRepo.GetByID(addressID)
	if err != nil {
		return errors.New("address not found")
	}

	if address.UserID != userID {
		return errors.New("unauthorized")
	}

	// Set as default
	if err := s.addressRepo.SetDefault(userID, addressID); err != nil {
		s.logger.Error("failed to set default address", zap.Error(err))
		return fmt.Errorf("failed to set default address: %w", err)
	}

	s.logger.Info("default address set", zap.Uint("address_id", addressID), zap.Uint("user_id", userID))
	return nil
}


