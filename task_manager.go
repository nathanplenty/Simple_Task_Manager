package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	DueDate   string `json:"due_date"`
	Completed bool   `json:"completed"`
}

var tasks []Task

func main() {
	initDatabase()
	http.HandleFunc("/tasks", handleTasks)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDatabase() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			name TEXT,
			due_date DATE,
			completed BOOLEAN
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Database initiated successfully")
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.Header.Get("ID") != "" {
			getTaskByID(w, r)
		} else {
			getTasks(w, r)
		}
	case http.MethodPost:
		createTask(w, r)
	case http.MethodPatch:
		updateTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		log.Printf("Method %s not allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Printf("Error encoding tasks to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Tasks gathered successfully")
	_, _ = w.Write([]byte("Tasks gathered successfully"))
}

func createTask(w http.ResponseWriter, r *http.Request) {
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
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Printf("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	result, err := db.Exec("INSERT INTO tasks(name, due_date, completed) VALUES(?, ?, ?)", name, dueDate, false)
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
	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	log.Println("Task created successfully")
	_, _ = w.Write([]byte("Task created successfully"))
}

func updateTask(w http.ResponseWriter, r *http.Request) {
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
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Printf("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var existingName string
	err = db.QueryRow("SELECT name FROM tasks WHERE id=?", id).Scan(&existingName)
	if err != nil {
		log.Printf("Error retrieving task to user name: %v", err)
		http.Error(w, "Error retrieving task to user name", http.StatusInternalServerError)
		return
	}
	name := r.Header.Get("name")
	if name != existingName {
		log.Println("Provided name does not match existing user name")
		http.Error(w, "Provided name does not match existing user name", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE tasks SET completed=? WHERE id=?", completed, id)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Completed = completed
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Task updated successfully")
	_, _ = w.Write([]byte("Task updated successfully"))
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		log.Println("Missing ID header")
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idHeader)
	if err != nil {
		log.Println("Invalid ID")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var existingName string
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Printf("Internal server error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	err = db.QueryRow("SELECT name FROM tasks WHERE id=?", id).Scan(&existingName)
	if err != nil {
		log.Printf("Error retrieving task to user name: %v", err)
		http.Error(w, "Error retrieving task to user name", http.StatusInternalServerError)
		return
	}
	name := r.Header.Get("name")
	if name != existingName {
		log.Println("Provided name does not match existing user name")
		http.Error(w, "Provided name does not match existing user name", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE tasks SET name='', due_date=NULL, completed=false WHERE id=?", id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Name = ""
			tasks[i].DueDate = ""
			tasks[i].Completed = false
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	log.Println("Task cleared successfully")
	_, _ = w.Write([]byte("Task cleared successfully"))
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		log.Println("Missing ID header")
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idHeader)
	if err != nil {
		log.Println("Invalid ID")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Println("Internal server error")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var (
		name      string
		dueDate   string
		completed bool
	)
	err = db.QueryRow("SELECT name, due_date, completed FROM tasks WHERE id=?", id).Scan(&name, &dueDate, &completed)
	if err != nil {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	task := Task{
		ID:        id,
		Name:      name,
		DueDate:   dueDate,
		Completed: completed,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
	log.Println("Task gathered successfully")
	_, _ = w.Write([]byte("Task gathered successfully"))
}
