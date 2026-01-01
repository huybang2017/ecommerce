package domain

import "time"

type Session struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"`

	RefreshTokenHash string `json:"refresh_token_hash"`
	IsRevoked        bool   `json:"is_revoked"`

	DeviceID   string     `json:"device_id"`
	DeviceType string     `json:"device_type"`
	UserAgent  string     `json:"user_agent"`
	IPAddress  string     `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	LastUsedAt time.Time  `json:"last_used_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

type SessionRepository interface {
	CreateSession(session *Session) error
	GetSession(sessionID string) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(sessionID string) error
	GetUserSessions(userID int64) ([]*Session, error)
	DeleteUserSessions(userID int64) error
	RevokeUserSessions(userID int64) error
	GetDeviceSessions(deviceID string) ([]*Session, error)
	DeleteDeviceSession(deviceID string) error
	UpdateLastUsed(sessionID string) error
	RevokeSession(sessionID string) error
	CleanupExpiredSessions() (int, error)
}

type SessionService interface {
	CreateSession(userID int64, refreshTokenHash, deviceID, deviceType, userAgent, ipAddress string) (*Session, error)
	ValidateSession(sessionID string) (*Session, error)
	RefreshSession(sessionID string) error
	RevokeSession(sessionID string) error
	GetActiveSessions(userID int64) ([]*Session, error)
	RevokeAllSessions(userID int64) error
	RevokeOtherSessions(userID int64, currentSessionID string) error
	DetectAnomalousSession(session *Session) (bool, string)
	RotateSession(oldSessionID string, newRefreshTokenHash string) (*Session, error)
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsRevoked && !s.IsExpired()
}

func (s *Session) DaysUntilExpiry() int {
	if s.IsExpired() {
		return 0
	}
	duration := time.Until(s.ExpiresAt)
	return int(duration.Hours() / 24)
}

func (s *Session) IsInactive(duration time.Duration) bool {
	return time.Since(s.LastUsedAt) > duration
}

func (s *Session) GetDeviceInfo() string {
	if s.DeviceType != "" {
		return s.DeviceType
	}
	return "unknown"
}
