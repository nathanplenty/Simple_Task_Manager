package domain

// Task struct represents the task model
type Task struct {
	TaskID    string `json:"task_id"`
	TaskName  string `json:"task_name"`
	DueDate   string `json:"due_date"`
	Completed string `json:"completed"`
}

// NewTask creates and returns a new Task instance
func NewTask(taskid, taskname, duedate, completed string) *Task {
	return &Task{
		TaskID:    taskid,
		TaskName:  taskname,
		DueDate:   duedate,
		Completed: completed,
	}
}

// TaskManager defines the interface for managing tasks
type TaskManager interface {
	CreateTask(task *Task, user *User) error
	ReadTask(task *Task, user *User) error
	UpdateTask(task *Task, user *User) error
	DeleteTask(task *Task, user *User) error
}
