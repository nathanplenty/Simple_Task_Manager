package sql

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Data struct {
	TaskID    int    `json:"task_id"`
	TaskName  string `json:"task_name"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	DueDate   string `json:"due_date"`
	Completed bool   `json:"completed"`
}

type Controller struct {
	db *Database
}

func NewController(db *Database) *Controller {
	return &Controller{db}
}

func (c *Controller) AddTask(w http.ResponseWriter, r *http.Request) {
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := c.db.db.Exec("INSERT INTO tasks (task, due_date, completed, user_id) VALUES (?, ?, ?, ?)",
		data.TaskName, data.DueDate, data.Completed, data.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, _ := result.LastInsertId()
	data.TaskID = int(lastID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func (c *Controller) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) GetTaskAll(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	query := "SELECT * FROM tasks"
	if userID != "" {
		query += " WHERE user_id = " + userID
	}

	rows, err := c.db.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Data
	for rows.Next() {
		var data Data
		if err := rows.Scan(&data.TaskID, &data.TaskName, &data.DueDate, &data.Completed, &data.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (c *Controller) GetTaskByID(w http.ResponseWriter, r *http.Request) {
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

	var data Data
	err = c.db.db.QueryRow("SELECT id, task, due_date, completed, user_id FROM tasks WHERE id=?", id).Scan(
		&data.TaskID, &data.TaskName, &data.DueDate, &data.Completed, &data.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *Controller) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
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
