package domain

import "time"

// Session represents a session in the system
type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}
