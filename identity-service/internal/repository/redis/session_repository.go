package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"identity-service/internal/domain"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type SessionRedisRepository struct {
	client *redis.Client
	logger *zap.Logger
	ctx    context.Context
}

func NewSessionRedisRepository(client *redis.Client, logger *zap.Logger) *SessionRedisRepository {
	return &SessionRedisRepository{
		client: client,
		logger: logger,
		ctx:    context.Background(),
	}
}

// Redis key patterns
const (
	sessionKeyPrefix       = "session:"        // session:{session_id}
	userSessionsKeyPrefix  = "user_sessions:"  // user_sessions:{user_id} -> Set of session_ids
	deviceSessionKeyPrefix = "device_session:" // device_session:{device_id} -> session_id
)

// CreateSession stores a new session in Redis
func (r *SessionRedisRepository) CreateSession(session *domain.Session) error {
	sessionKey := fmt.Sprintf("%s%s", sessionKeyPrefix, session.ID)

	// Marshal session to JSON
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Calculate TTL
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	// Store session with expiration
	if err := r.client.Set(r.ctx, sessionKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Add session to user's sessions set
	userSessionsKey := fmt.Sprintf("%s%d", userSessionsKeyPrefix, session.UserID)
	if err := r.client.SAdd(r.ctx, userSessionsKey, session.ID).Err(); err != nil {
		r.logger.Warn("failed to add session to user set",
			zap.Error(err),
			zap.Int64("user_id", session.UserID),
		)
	}
	// Set TTL for user sessions set
	r.client.Expire(r.ctx, userSessionsKey, ttl)

	// Map device to session
	if session.DeviceID != "" {
		deviceKey := fmt.Sprintf("%s%s", deviceSessionKeyPrefix, session.DeviceID)
		if err := r.client.Set(r.ctx, deviceKey, session.ID, ttl).Err(); err != nil {
			r.logger.Warn("failed to map device to session",
				zap.Error(err),
				zap.String("device_id", session.DeviceID),
			)
		}
	}

	r.logger.Info("session created",
		zap.String("session_id", session.ID),
		zap.Int64("user_id", session.UserID),
		zap.String("device_id", session.DeviceID),
		zap.Duration("ttl", ttl),
	)

	return nil
}

// GetSession retrieves a session by ID
func (r *SessionRedisRepository) GetSession(sessionID string) (*domain.Session, error) {
	sessionKey := fmt.Sprintf("%s%s", sessionKeyPrefix, sessionID)

	data, err := r.client.Get(r.ctx, sessionKey).Bytes()
	if err == redis.Nil {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session domain.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Validate session
	if session.IsExpired() {
		r.DeleteSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	if session.IsRevoked {
		return nil, fmt.Errorf("session revoked")
	}

	return &session, nil
}

// UpdateSession updates an existing session
func (r *SessionRedisRepository) UpdateSession(session *domain.Session) error {
	// Simply overwrite by creating again
	return r.CreateSession(session)
}

// DeleteSession removes a session from Redis
func (r *SessionRedisRepository) DeleteSession(sessionID string) error {
	sessionKey := fmt.Sprintf("%s%s", sessionKeyPrefix, sessionID)

	// Get session first to cleanup related keys
	session, _ := r.GetSession(sessionID)

	// Delete session
	if err := r.client.Del(r.ctx, sessionKey).Err(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Cleanup user sessions set
	if session != nil {
		userSessionsKey := fmt.Sprintf("%s%d", userSessionsKeyPrefix, session.UserID)
		r.client.SRem(r.ctx, userSessionsKey, sessionID)

		// Cleanup device mapping
		if session.DeviceID != "" {
			deviceKey := fmt.Sprintf("%s%s", deviceSessionKeyPrefix, session.DeviceID)
			r.client.Del(r.ctx, deviceKey)
		}
	}

	r.logger.Info("session deleted",
		zap.String("session_id", sessionID),
	)

	return nil
}

// GetUserSessions retrieves all sessions for a user
func (r *SessionRedisRepository) GetUserSessions(userID int64) ([]*domain.Session, error) {
	userSessionsKey := fmt.Sprintf("%s%d", userSessionsKeyPrefix, userID)

	sessionIDs, err := r.client.SMembers(r.ctx, userSessionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	sessions := make([]*domain.Session, 0, len(sessionIDs))
	for _, sessionID := range sessionIDs {
		session, err := r.GetSession(sessionID)
		if err == nil && session.IsValid() {
			sessions = append(sessions, session)
		} else {
			// Cleanup invalid session from set
			r.client.SRem(r.ctx, userSessionsKey, sessionID)
		}
	}

	return sessions, nil
}

// DeleteUserSessions deletes all sessions for a user
func (r *SessionRedisRepository) DeleteUserSessions(userID int64) error {
	sessions, err := r.GetUserSessions(userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := r.DeleteSession(session.ID); err != nil {
			r.logger.Error("failed to delete session",
				zap.Error(err),
				zap.String("session_id", session.ID),
			)
		}
	}

	// Clear user sessions set
	userSessionsKey := fmt.Sprintf("%s%d", userSessionsKeyPrefix, userID)
	r.client.Del(r.ctx, userSessionsKey)

	r.logger.Info("all user sessions deleted",
		zap.Int64("user_id", userID),
		zap.Int("count", len(sessions)),
	)

	return nil
}

// RevokeUserSessions revokes all sessions for a user (soft delete)
func (r *SessionRedisRepository) RevokeUserSessions(userID int64) error {
	sessions, err := r.GetUserSessions(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, session := range sessions {
		session.IsRevoked = true
		session.RevokedAt = &now
		if err := r.UpdateSession(session); err != nil {
			r.logger.Error("failed to revoke session",
				zap.Error(err),
				zap.String("session_id", session.ID),
			)
		}
	}

	r.logger.Info("all user sessions revoked",
		zap.Int64("user_id", userID),
		zap.Int("count", len(sessions)),
	)

	return nil
}

// GetDeviceSessions retrieves sessions for a specific device
func (r *SessionRedisRepository) GetDeviceSessions(deviceID string) ([]*domain.Session, error) {
	deviceKey := fmt.Sprintf("%s%s", deviceSessionKeyPrefix, deviceID)

	sessionID, err := r.client.Get(r.ctx, deviceKey).Result()
	if err == redis.Nil {
		return []*domain.Session{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get device session: %w", err)
	}

	session, err := r.GetSession(sessionID)
	if err != nil {
		return []*domain.Session{}, nil
	}

	return []*domain.Session{session}, nil
}

// DeleteDeviceSession deletes session associated with a device
func (r *SessionRedisRepository) DeleteDeviceSession(deviceID string) error {
	sessions, err := r.GetDeviceSessions(deviceID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if err := r.DeleteSession(session.ID); err != nil {
			return err
		}
	}

	return nil
}

// UpdateLastUsed updates the last used timestamp of a session
func (r *SessionRedisRepository) UpdateLastUsed(sessionID string) error {
	session, err := r.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.LastUsedAt = time.Now()
	return r.UpdateSession(session)
}

// RevokeSession revokes a specific session
func (r *SessionRedisRepository) RevokeSession(sessionID string) error {
	session, err := r.GetSession(sessionID)
	if err != nil {
		return err
	}

	now := time.Now()
	session.IsRevoked = true
	session.RevokedAt = &now

	if err := r.UpdateSession(session); err != nil {
		return err
	}

	r.logger.Info("session revoked",
		zap.String("session_id", sessionID),
		zap.Int64("user_id", session.UserID),
	)

	return nil
}

// CleanupExpiredSessions removes all expired sessions (maintenance task)
func (r *SessionRedisRepository) CleanupExpiredSessions() (int, error) {
	// In Redis, expired keys are auto-deleted, but we need to cleanup sets
	// This is a best-effort cleanup for orphaned references

	count := 0

	// Scan for all user_sessions sets
	iter := r.client.Scan(r.ctx, 0, userSessionsKeyPrefix+"*", 0).Iterator()
	for iter.Next(r.ctx) {
		userSessionsKey := iter.Val()

		sessionIDs, err := r.client.SMembers(r.ctx, userSessionsKey).Result()
		if err != nil {
			continue
		}

		for _, sessionID := range sessionIDs {
			sessionKey := fmt.Sprintf("%s%s", sessionKeyPrefix, sessionID)
			exists, err := r.client.Exists(r.ctx, sessionKey).Result()
			if err != nil || exists == 0 {
				// Session expired but still in set, remove it
				r.client.SRem(r.ctx, userSessionsKey, sessionID)
				count++
			}
		}

		// If set is empty, delete it
		size, _ := r.client.SCard(r.ctx, userSessionsKey).Result()
		if size == 0 {
			r.client.Del(r.ctx, userSessionsKey)
		}
	}

	if err := iter.Err(); err != nil {
		return count, fmt.Errorf("scan error: %w", err)
	}

	r.logger.Info("expired sessions cleaned up",
		zap.Int("count", count),
	)

	return count, nil
}
