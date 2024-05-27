package database

import (
	"Simple_Task_Manager/pkg/domain"
)

// Database interface defines the methods for interacting with the database
type Database interface {
	Connect(dbPath string) error
	CreateUser(user *domain.User) error
	CheckUser(user *domain.User) error
	LoginUser(user *domain.User) error
	CreateTask(task *domain.Task, user *domain.User) error
	ReadTask(task *domain.Task, user *domain.User) error
	UpdateTask(task *domain.Task, user *domain.User) error
	DeleteTask(task *domain.Task, user *domain.User) error
}
