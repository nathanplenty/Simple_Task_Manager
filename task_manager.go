package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
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
		log.Fatal(err)
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
		log.Fatal(err)
	}
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	if name == "" {
		http.Error(w, "Missing name header", http.StatusBadRequest)
		return
	}

	dueDate := r.Header.Get("due_date")
	if dueDate == "" {
		http.Error(w, "Missing due date header", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO tasks(name, due_date, completed) VALUES(?, ?, ?)", name, dueDate, false)
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting last inserted ID", http.StatusInternalServerError)
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
	fmt.Fprintf(w, "Task created successfully")
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idHeader)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	completedHeader := r.Header.Get("completed")
	completed, err := strconv.ParseBool(completedHeader)
	if err != nil {
		http.Error(w, "Invalid completed header", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE tasks SET completed=? WHERE id=?", completed, id)
	if err != nil {
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
	fmt.Fprintf(w, "Task updated successfully")
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idHeader)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE tasks SET name='', due_date=NULL, completed=false WHERE id=?", id)
	if err != nil {
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
	fmt.Fprintf(w, "Task cleared successfully")
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	idHeader := r.Header.Get("ID")
	if idHeader == "" {
		http.Error(w, "Missing ID header", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idHeader)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
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
}
