package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	db *sql.DB
}

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
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
		if r.Header.Get("ID") != "" {
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
	rows, err := app.db.Query("SELECT id, name, due_date, completed FROM tasks")
	if err != nil {
		log.Printf("Error querying tasks from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Name, &task.DueDate, &task.Completed)
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
	name := r.Header.Get("name")
	if name == "" {
		log.Println("Missing name header")
		http.Error(w, "Missing name header", http.StatusBadRequest)
		return
	}
	dueDate := r.Header.Get("due_date")
	if dueDate == "" {
		log.Println("Missing due date header")
		http.Error(w, "Missing due date header", http.StatusBadRequest)
		return
	}
	result, err := app.db.Exec("INSERT INTO tasks(name, due_date, completed) VALUES(?, ?, ?)", name, dueDate, false)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last inserted ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	task := Task{
		ID:        int(id),
		Name:      name,
		DueDate:   dueDate,
		Completed: false,
	}
	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Error encoding task to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println("Task created successfully")
}

func (app *App) updateTask(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		log.Println("Missing ID header")
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idHeader)
	if err != nil {
		log.Println("Invalid ID header")
		http.Error(w, "Invalid ID header", http.StatusBadRequest)
		return
	}
	completedHeader := r.Header.Get("completed")
	completed, err := strconv.ParseBool(completedHeader)
	if err != nil {
		log.Println("Invalid completed header")
		http.Error(w, "Invalid completed header", http.StatusBadRequest)
		return
	}
	result, err := app.db.Exec("UPDATE tasks SET completed=? WHERE id=?", completed, id)
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
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		log.Println("Missing ID header")
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idHeader)
	if err != nil {
		log.Println("Invalid ID header")
		http.Error(w, "Invalid ID header", http.StatusBadRequest)
		return
	}
	result, err := app.db.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
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
	log.Println("Task deleted successfully")
	_, _ = w.Write([]byte("Task deleted successfully"))
}

func (app *App) getTaskByID(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		log.Println("Missing ID header")
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idHeader)
	if err != nil {
		log.Println("Invalid ID header")
		http.Error(w, "Invalid ID header", http.StatusBadRequest)
		return
	}
	var task Task
	err = app.db.QueryRow("SELECT id, name, due_date, completed FROM tasks WHERE id=?", id).Scan(&task.ID, &task.Name, &task.DueDate, &task.Completed)
	if err != nil {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Error encoding task to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Task retrieved successfully")
}
