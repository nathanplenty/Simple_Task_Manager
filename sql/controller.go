package sql

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	userID, err := c.CheckUser(data.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := c.db.db.Exec("INSERT INTO tasks (task_name, due_date, completed, user_id) VALUES (?, ?, ?, ?)",
		data.TaskName, data.DueDate, false, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if lastID == 0 {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	data.TaskID = int(lastID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func (c *Controller) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("task_id")
	if idStr == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	_, err = c.db.db.Exec("DELETE FROM tasks WHERE task_id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetTaskAll(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	query := "SELECT task_id, task_name, due_date, completed, user_id FROM tasks"
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
	idStr := r.URL.Query().Get("task_id")
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

	err = c.db.db.QueryRow("SELECT task_id, task_name, due_date, completed, user_id FROM tasks WHERE task_id=?", id).Scan(
		&data.TaskID, &data.TaskName, &data.DueDate, &data.Completed, &data.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *Controller) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("task_id")
	if idStr == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	_, err = c.db.db.Exec("UPDATE tasks SET completed=? WHERE task_id=?", true, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) CheckUser(userName string) (int, error) {
	var userID int
	err := c.db.db.QueryRow("SELECT user_id FROM users WHERE user_name=?", userName).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		result, err := c.db.db.Exec("INSERT INTO users (user_name) VALUES (?)", userName)
		if err != nil {
			return 0, err
		}
		lastID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		userID = int(lastID)
	} else if err != nil {
		return 0, err
	}
	return userID, nil
}
