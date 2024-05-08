package domain

// TaskManager defines the interface for managing tasks
type TaskManager interface {
	CreateTask(task *Task, taskid, username, password string) error
	ReadTask(taskid, username, password string) (*Task, error)
	UpdateTask(task *Task, taskid, username, password string) error
	DeleteTask(taskid, username, password string) error
}
