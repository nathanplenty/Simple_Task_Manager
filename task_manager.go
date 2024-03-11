package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	db *sql.DB
}

type Task struct {
	TaskID    int    `json:"task_id"`
	TaskName  string `json:"task_name"`
	UserName  string `json:"user_name"`
	DueDate   string `json:"due_date"`
	Completed bool   `json:"completed"`
}

func (app *App) initDatabase() {
	var err error
	app.db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	_, err = app.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			due_date DATE,
			completed BOOLEAN
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Database initiated successfully")
}

func main() {
	app := &App{}
	app.initDatabase()
	http.HandleFunc("/tasks", app.handleTasks)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (app *App) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.Header.Get("task_id") != "" {
			app.getTaskByID(w, r)
		} else {
			app.getTasks(w, r)
		}
	case http.MethodPost:
		app.createTask(w, r)
	case http.MethodPatch:
		app.updateTask(w, r)
	case http.MethodDelete:
		app.deleteTask(w, r)
	default:
		log.Printf("Method %s not allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) getTasks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tasks []Task
	rows, err := app.db.Query("SELECT t.task_id, t.name, t.due_date, t.completed, u.username FROM tasks AS t INNER JOIN task_assignments AS ta ON t.task_id = ta.task_id INNER JOIN users AS u ON ta.user_id = u.user_id")
	if err != nil {
		log.Printf("Error querying tasks from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed, &task.UserName)
		if err != nil {
			log.Printf("Error scanning task row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over task rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("Error encoding tasks to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Tasks gathered successfully")
}

func (app *App) createTask(w http.ResponseWriter, r *http.Request) {
	userName := r.Header.Get("user_name")
	if userName == "" {
		log.Println("Missing user name header")
		http.Error(w, "Missing user name header", http.StatusBadRequest)
		return
	}
	taskName := r.Header.Get("task_name")
	if taskName == "" {
		log.Println("Missing task name header")
		http.Error(w, "Missing task name header", http.StatusBadRequest)
		return
	}
	dueDate := r.Header.Get("due_date")
	if dueDate == "" {
		log.Println("Missing due date header")
		http.Error(w, "Missing due date header", http.StatusBadRequest)
		return
	}

	// Insert task into tasks table
	result, err := app.db.Exec("INSERT INTO tasks(name, due_date, completed) VALUES(?, ?, ?)", taskName, dueDate, false)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}
	taskID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get or insert user ID
	var userID int
	err = app.db.QueryRow("SELECT user_id FROM users WHERE username = ?", userName).Scan(&userID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		result, err := app.db.Exec("INSERT INTO users(username) VALUES(?)", userName)
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

	// Assign task to user
	_, err = app.db.Exec("INSERT INTO task_assignments(task_id, user_id) VALUES(?, ?)", taskID, userID)
	if err != nil {
		log.Printf("Error assigning task to user: %v", err)
		http.Error(w, "Error assigning task to user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Task created successfully")
}

func (app *App) updateTask(w http.ResponseWriter, r *http.Request) {
	userIDHeader := r.Header.Get("user_id")
	if userIDHeader == "" {
		log.Println("Missing user ID header")
		http.Error(w, "Missing user ID header", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		log.Println("Invalid user ID header")
		http.Error(w, "Invalid user ID header", http.StatusBadRequest)
		return
	}
	taskIDHeader := r.Header.Get("task_id")
	if taskIDHeader == "" {
		log.Println("Missing task ID header")
		http.Error(w, "Missing task ID header", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(taskIDHeader)
	if err != nil {
		log.Println("Invalid task ID header")
		http.Error(w, "Invalid task ID header", http.StatusBadRequest)
		return
	}
	var exists bool
	err = app.db.QueryRow("SELECT EXISTS(SELECT 1 FROM task_assignments WHERE task_id=? AND user_id=?)", taskID, userID).Scan(&exists)
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
	result, err := app.db.Exec("UPDATE tasks SET completed=true WHERE task_id=?", taskID)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Task updated successfully")
	_, _ = w.Write([]byte("Task updated successfully"))
}

func (app *App) deleteTask(w http.ResponseWriter, r *http.Request) {
	userIDHeader := r.Header.Get("user_id")
	if userIDHeader == "" {
		log.Println("Missing user ID header")
		http.Error(w, "Missing user ID header", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		log.Println("Invalid user ID header")
		http.Error(w, "Invalid user ID header", http.StatusBadRequest)
		return
	}
	taskIDHeader := r.Header.Get("task_id")
	if taskIDHeader == "" {
		log.Println("Missing task ID header")
		http.Error(w, "Missing task ID header", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(taskIDHeader)
	if err != nil {
		log.Println("Invalid task ID header")
		http.Error(w, "Invalid task ID header", http.StatusBadRequest)
		return
	}
	var (
		taskExists       bool
		userExists       bool
		assignmentExists bool
	)
	err = app.db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE task_id=?)", taskID).Scan(&taskExists)
	if err != nil {
		log.Printf("Error checking task existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = app.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE user_id=?)", userID).Scan(&userExists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = app.db.QueryRow("SELECT EXISTS(SELECT 1 FROM task_assignments WHERE task_id=? AND user_id=?)", taskID, userID).Scan(&assignmentExists)
	if err != nil {
		log.Printf("Error checking task assignment: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !taskExists || !userExists || !assignmentExists {
		log.Println("Task, user, or assignment not found")
		http.Error(w, "Task, user, or assignment not found", http.StatusNotFound)
		return
	}
	_, err = app.db.Exec("UPDATE tasks SET name='X', due_date='0001-01-01', completed=false WHERE task_id=?", taskID)
	if err != nil {
		log.Printf("Error anonymizing task: %v", err)
		http.Error(w, "Error anonymizing task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Task content anonymized successfully")
}

func (app *App) getTaskByID(w http.ResponseWriter, r *http.Request) {
	taskIDHeader := r.Header.Get("task_id")
	if taskIDHeader == "" {
		log.Println("Missing task ID header")
		http.Error(w, "Missing task ID header", http.StatusBadRequest)
		return
	}
	taskID, err := strconv.Atoi(taskIDHeader)
	if err != nil {
		log.Println("Invalid task ID header")
		http.Error(w, "Invalid task ID header", http.StatusBadRequest)
		return
	}
	var task Task
	var userName string
	err = app.db.QueryRow("SELECT t.task_id, t.name, t.due_date, t.completed, u.username FROM tasks AS t INNER JOIN task_assignments AS ta ON t.task_id = ta.task_id INNER JOIN users AS u ON ta.user_id = u.user_id WHERE t.task_id=?", taskID).Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed, &userName)
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
	task.UserName = userName
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Error encoding task to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Task retrieved successfully")
}
