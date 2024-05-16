package domain

type Session struct {
	SessionID    string `json:"session_id"`
	UserID       string `json:"user_id"`
	CreatedAt    string `json:"created_at"`
	LastAccessed string `json:"last_accessed"`
}

// NewSession creates and returns a new Session instance
func NewSession(sessionid, userid, createdat, lastaccessed string) *Session {
	return &Session{
		SessionID:    sessionid,
		UserID:       userid,
		CreatedAt:    createdat,
		LastAccessed: lastaccessed,
	}
}

// SessionManager defines the interface for managing users
type SessionManager interface {
	CreateSession(session *Session) error
	DeleteSession(session *Session) error
	GetByToken(session *Session) error
	GetByID(session *Session) error
}
