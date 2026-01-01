package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"identity-service/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService contains the business logic for authentication
type AuthService struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	sessionRepo      domain.SessionRepository
	logger           *zap.Logger
	jwtSecret        string
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	sessionRepo domain.SessionRepository,
	logger *zap.Logger,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		sessionRepo:      sessionRepo,
		logger:           logger,
		jwtSecret:        jwtSecret,
	}
}

// RegisterRequest represents the request to register a new user
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

// LoginRequest represents the request to login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
// NOTE: Token should NOT be in response body for production
// Instead, it should be set as HttpOnly cookie by the handler
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"` // Deprecated: Use SessionID instead
	SessionID    string       `json:"session_id"`    // Session ID for session-based auth
	User         *domain.User `json:"user"`
	ExpiresIn    int64        `json:"expires_in"` // seconds until access token expires
}

// Register creates a new user account
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	existing, _ := s.userRepo.GetByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	existing, _ = s.userRepo.GetByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		PhoneNumber:  req.PhoneNumber,
		Role:         "BUYER",
		Status:       "ACTIVE",
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("user registered", zap.Uint("user_id", user.ID), zap.String("email", user.Email))

	// Generate Access Token (short-lived: 15 minutes)
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token (long-lived: 7 days)
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session in Redis with refresh token hash
	session := &domain.Session{
		ID:               uuid.New().String(),
		UserID:           int64(user.ID),
		RefreshTokenHash: hashToken(refreshToken),
		IsRevoked:        false,
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 7), // 7 days
		CreatedAt:        time.Now(),
		LastUsedAt:       time.Now(),
		// Device info will be set by handler from request headers
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		s.logger.Error("failed to create session", zap.Error(err))
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // For backward compatibility
		SessionID:    session.ID,
		User:         user,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

// Login authenticates a user and returns a JWT token with session
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check user status
	if user.Status != "ACTIVE" {
		return nil, errors.New("account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	s.logger.Info("user logged in", zap.Uint("user_id", user.ID), zap.String("email", user.Email))

	// Generate Access Token (short-lived: 15 minutes)
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate Refresh Token (long-lived: 7 days)
	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session in Redis with refresh token hash
	session := &domain.Session{
		ID:               uuid.New().String(),
		UserID:           int64(user.ID),
		RefreshTokenHash: hashToken(refreshToken),
		IsRevoked:        false,
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 7), // 7 days
		CreatedAt:        time.Now(),
		LastUsedAt:       time.Now(),
		// Device info will be set by handler from request headers
	}

	if err := s.sessionRepo.CreateSession(session); err != nil {
		s.logger.Error("failed to create session", zap.Error(err))
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // For backward compatibility
		SessionID:    session.ID,
		User:         user,
		ExpiresIn:    900, // 15 minutes in seconds
	}, nil
}

// generateAccessToken generates a short-lived JWT access token (15 minutes)
func (s *AuthService) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"type":    "access",                                // Token type identifier
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken generates a long-lived refresh token (7 days) and stores it in database
func (s *AuthService) generateRefreshToken(user *domain.User) (string, error) {
	// Generate random token string
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	tokenString := base64.URLEncoding.EncodeToString(tokenBytes)

	// Create refresh token record
	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
		IsRevoked: false,
	}

	// Save to database
	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		s.logger.Error("failed to save refresh token", zap.Error(err))
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *AuthService) ValidateToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)
		return userID, role, nil
	}

	return 0, "", errors.New("invalid token")
}

// RefreshRequest represents the request to refresh access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshAccessTokenBySession validates session_id and issues a new access token (session-based)
func (s *AuthService) RefreshAccessTokenBySession(sessionID string) (*AuthResponse, error) {
	// Get session from Redis
	session, err := s.sessionRepo.GetSession(sessionID)
	if err != nil {
		s.logger.Warn("session not found", zap.String("session_id", sessionID), zap.Error(err))
		return nil, errors.New("invalid session")
	}

	// Validate session
	if !session.IsValid() {
		s.logger.Warn("session is invalid or revoked", zap.String("session_id", sessionID))
		return nil, errors.New("session expired or revoked")
	}

	// Get user
	user, err := s.userRepo.GetByID(uint(session.UserID))
	if err != nil {
		s.logger.Error("user not found for session", zap.Int64("user_id", session.UserID))
		return nil, errors.New("user not found")
	}

	// Check user status
	if user.Status != "ACTIVE" {
		return nil, errors.New("account is not active")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Update last used timestamp in session
	if err := s.sessionRepo.UpdateLastUsed(sessionID); err != nil {
		s.logger.Warn("failed to update session last used", zap.String("session_id", sessionID), zap.Error(err))
		// Don't fail the request, just log the warning
	}

	s.logger.Info("access token refreshed via session", zap.Uint("user_id", user.ID), zap.String("session_id", sessionID))

	return &AuthResponse{
		AccessToken: accessToken,
		SessionID:   sessionID,
		User:        user,
		ExpiresIn:   900, // 15 minutes
	}, nil
}

// RefreshAccessToken validates refresh token and issues a new access token (legacy method)
func (s *AuthService) RefreshAccessToken(refreshTokenString string) (*AuthResponse, error) {
	// Get refresh token from database
	refreshToken, err := s.refreshTokenRepo.GetByToken(refreshTokenString)
	if err != nil {
		s.logger.Warn("refresh token not found", zap.Error(err))
		return nil, errors.New("invalid refresh token")
	}

	// Validate refresh token
	if !refreshToken.IsValid() {
		s.logger.Warn("refresh token is invalid or revoked", zap.Uint("user_id", refreshToken.UserID))
		return nil, errors.New("refresh token expired or revoked")
	}

	// Get user
	user, err := s.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		s.logger.Error("user not found for refresh token", zap.Uint("user_id", refreshToken.UserID))
		return nil, errors.New("user not found")
	}

	// Check user status
	if user.Status != "ACTIVE" {
		return nil, errors.New("account is not active")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	s.logger.Info("access token refreshed", zap.Uint("user_id", user.ID))

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString, // Return same refresh token
		User:         user,
		ExpiresIn:    900, // 15 minutes
	}, nil
}

// LogoutBySession revokes a specific session by session ID
func (s *AuthService) LogoutBySession(sessionID string) error {
	err := s.sessionRepo.DeleteSession(sessionID)
	if err != nil {
		s.logger.Error("failed to delete session", zap.String("session_id", sessionID), zap.Error(err))
		return fmt.Errorf("failed to logout: %w", err)
	}

	s.logger.Info("user logged out via session", zap.String("session_id", sessionID))
	return nil
}

// Logout revokes all refresh tokens for a user (legacy method)
func (s *AuthService) Logout(userID uint) error {
	err := s.refreshTokenRepo.RevokeAllByUserID(userID)
	if err != nil {
		s.logger.Error("failed to revoke refresh tokens", zap.Uint("user_id", userID), zap.Error(err))
		return fmt.Errorf("failed to logout: %w", err)
	}

	s.logger.Info("user logged out", zap.Uint("user_id", userID))
	return nil
}

// hashToken creates SHA256 hash of a token for secure storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
