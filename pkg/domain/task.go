package domain

// Task represents a task in the system
type Task struct {
	TaskID    string `json:"task_id"`
	TaskName  string `json:"task_name"`
	DueDate   string `json:"due_date"`
	Completed string `json:"completed"`
}

// NewTask creates a new Task instance
func NewTask(taskID, taskName, dueDate, completed string) *Task {
	return &Task{
		TaskID:    taskID,
		TaskName:  taskName,
		DueDate:   dueDate,
		Completed: completed,
	}
}
