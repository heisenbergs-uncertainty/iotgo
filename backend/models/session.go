package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID       uint      `gorm:"index"`
	User         User      `gorm:"foreignKey:UserID"`
	Token        string    `gorm:"size:255;uniqueIndex"`
	ExpiresAt    time.Time `gorm:"index"`
	IP           string    `gorm:"size:50"`
	UserAgent    string    `gorm:"size:255"`
	LastActivity time.Time
	Active       bool `gorm:"default:true"`
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// UpdateLastActivity updates the last activity timestamp
func (s *Session) UpdateLastActivity(db *gorm.DB) error {
	now := time.Now()
	s.LastActivity = now
	return db.Model(s).Update("last_activity", now).Error
}

// Invalidate marks the session as inactive
func (s *Session) Invalidate(db *gorm.DB) error {
	s.Active = false
	return db.Model(s).Update("active", false).Error
}

// CreateSession creates a new session for a user
func CreateSession(db *gorm.DB, userID uint, token string, expiresAt time.Time, ip string, userAgent string) (*Session, error) {
	session := &Session{
		UserID:       userID,
		Token:        token,
		ExpiresAt:    expiresAt,
		IP:           ip,
		UserAgent:    userAgent,
		LastActivity: time.Now(),
		Active:       true,
	}

	err := db.Create(session).Error
	return session, err
}

// FindActiveSessionByToken finds an active session by token
func FindActiveSessionByToken(db *gorm.DB, token string) (*Session, error) {
	var session Session
	err := db.Where("token = ? AND active = ? AND expires_at > ?", token, true, time.Now()).First(&session).Error
	return &session, err
}

// InvalidateUserSessions invalidates all sessions for a user
func InvalidateUserSessions(db *gorm.DB, userID uint) error {
	return db.Model(&Session{}).
		Where("user_id = ? AND active = ?", userID, true).
		Update("active", false).Error
}

// CleanExpiredSessions removes all expired sessions from the database
func CleanExpiredSessions(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
}
