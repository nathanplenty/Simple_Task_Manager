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
