package taskmanager

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	DB *sql.DB
}

type Task struct {
	TaskID    int    `json:"task_id"`
	TaskName  string `json:"task_name"`
	DueDate   string `json:"due_date"`
	Completed bool   `json:"completed"`
}

type User struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

func (app *App) CheckDatabase() error {
	var err error
	app.DB, err = sql.Open("sqlite3", "./database/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}
	rows, err := app.DB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name IN ('tasks', 'users')")
	if err != nil {
		log.Fatalf("Error querying database tables: %v", err)
		return err
	}
	defer rows.Close()

	log.Println("Database initialized successfully")
	return nil
}

func (app *App) GetTasks(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("task_id")
	if taskIDStr != "" {
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			log.Println("Invalid task ID:", err)
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}
		app.GetTaskByID(w, taskID)
		return
	}
	var tasks []Task
	rows, err := app.DB.Query("SELECT t.task_id, t.task_name, t.due_date, t.completed, u.user_id, u.user_name FROM tasks t INNER JOIN users u ON t.user_id = u.user_id")
	if err != nil {
		log.Printf("Error querying tasks from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var user User
		err = rows.Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed, &user.UserID, &user.UserName)
		if err != nil {
			log.Printf("Error scanning task row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Error iterating over task rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Error encoding tasks to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		log.Printf("Error writing response body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Tasks gathered successfully")
}

func (app *App) CreateTask(w http.ResponseWriter, userName, taskName, dueDate string) {
	var requestBody = struct {
		UserName string `json:"user_name"`
		TaskName string `json:"task_name"`
		DueDate  string `json:"due_date"`
	}{
		UserName: userName,
		TaskName: taskName,
		DueDate:  dueDate,
	}

	if requestBody.UserName == "" {
		log.Println("Missing user name in request body")
		http.Error(w, "Missing user name in request body", http.StatusBadRequest)
		return
	}
	if requestBody.TaskName == "" {
		log.Println("Missing task name in request body")
		http.Error(w, "Missing task name in request body", http.StatusBadRequest)
		return
	}
	if requestBody.DueDate == "" {
		log.Println("Missing due date in request body")
		http.Error(w, "Missing due date in request body", http.StatusBadRequest)
		return
	}

	var userID int
	err := app.DB.QueryRow("SELECT user_id FROM users WHERE user_name = ?", requestBody.UserName).Scan(&userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		result, err := app.DB.Exec("INSERT INTO users(user_name) VALUES(?)", requestBody.UserName)
		if err != nil {
			log.Printf("Error creating new user: %v", err)
			http.Error(w, "Error creating new user", http.StatusInternalServerError)
			return
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Printf("Error getting last inserted ID: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		userID = int(lastInsertID)
	case err != nil:
		log.Printf("Error checking user existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result, err := app.DB.Exec("INSERT INTO tasks(task_name, due_date, completed, user_id) VALUES(?, ?, ?, ?)", requestBody.TaskName, requestBody.DueDate, false, userID)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Task created successfully")
	_, _ = w.Write([]byte("Task created successfully"))
}

func (app *App) UpdateTask(w http.ResponseWriter, taskID, userID int) {
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id=? AND user_id=?)", taskID, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking task assignment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Println("Task assignment not found or user does not have permission")
		http.Error(w, "Task assignment not found or user does not have permission", http.StatusNotFound)
		return
	}

	_, err = app.DB.Exec("UPDATE tasks SET completed=true WHERE task_id=?", taskID)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Task updated successfully")
	_, _ = w.Write([]byte("Task updated successfully"))
}

func (app *App) DeleteTask(w http.ResponseWriter, taskID, userID int) {
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id=? AND user_id=?)", taskID, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking task assignment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Println("Task assignment not found or user does not have permission")
		http.Error(w, "Task assignment not found or user does not have permission", http.StatusNotFound)
		return
	}

	_, err = app.DB.Exec("UPDATE tasks SET task_name='X', due_date='0001-01-01', completed=false WHERE task_id=?", taskID)
	if err != nil {
		log.Printf("Error anonymizing task: %v", err)
		http.Error(w, "Error anonymizing task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Task content anonymized successfully")
	_, _ = w.Write([]byte("Task anonymized successfully"))
}

func (app *App) GetTaskByID(w http.ResponseWriter, TaskID int) {
	var task Task
	err := app.DB.QueryRow("SELECT task_id, task_name, due_date, completed FROM tasks WHERE task_id=?", TaskID).Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("Error retrieving task: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("Error encoding task to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Task retrieved successfully")
}
