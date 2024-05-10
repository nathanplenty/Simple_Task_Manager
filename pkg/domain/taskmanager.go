package domain

// TaskManager defines the interface for managing tasks
type TaskManager interface {
	CreateTask(task *Task, user *User) error
	ReadTask(task *Task, user *User) error
	UpdateTask(task *Task, user *User) error
	DeleteTask(task *Task, user *User) error
}
