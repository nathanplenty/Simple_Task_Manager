package database

import "Simple_Task_Manager/pkg/domain"

// Database interface defines the methods for interacting with the database
type Database interface {
	Connect(dbPath string) error
	CreateUser(user *domain.User) error
	CheckUser(user *domain.User) error
	CreateTask(task *domain.Task, user *domain.User, session *domain.Session) error
	ReadTask(task *domain.Task, user *domain.User, session *domain.Session) error
	UpdateTask(task *domain.Task, user *domain.User, session *domain.Session) error
	DeleteTask(task *domain.Task, user *domain.User, session *domain.Session) error
	CreateSession(userID int) (string, error)
	UpdateSession(sessionID string) error
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*domain.Session, error)
	CheckPassword(user *domain.User) error
	GetUserIDByUsername(userName string, userID *int) error
}
