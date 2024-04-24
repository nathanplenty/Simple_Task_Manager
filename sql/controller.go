package sql

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Task struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	DueDate   string `json:"due_date"`
	Completed bool   `json:"completed"`
}

type Controller struct {
	db *Database
}

func NewController(db *Database) *Controller {
	return &Controller{db}
}

func (c *Controller) GetAllTasks(w http.ResponseWriter) {
	rows, err := c.db.db.Query("SELECT * FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Task, &task.DueDate, &task.Completed); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (c *Controller) AddTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := c.db.db.Exec("INSERT INTO tasks (task, due_date, completed) VALUES (?, ?, ?)",
		task.Task, task.DueDate, task.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := result.LastInsertId()
	task.ID = int(lastID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (c *Controller) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	_, err = c.db.db.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	_, err = c.db.db.Exec("UPDATE tasks SET completed=? WHERE id=?", true, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
