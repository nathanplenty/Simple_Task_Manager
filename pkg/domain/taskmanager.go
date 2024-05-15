package domain

// TaskManager defines the interface for managing tasks
type TaskManager interface {
	CreateUser(user *User) error
	CheckUser(user *User) error
	CreateTask(task *Task, user *User) error
	ReadTask(task *Task, user *User) error
	UpdateTask(task *Task, user *User) error
	DeleteTask(task *Task, user *User) error
}

// TaskManager für Tasks
// UserManager für User
