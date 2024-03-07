package main

//goland:noinspection Annotator,Annotator,Annotator
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
	ID          int    `json:"id"`
	Description string `json:"description"`
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT
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
	description := r.Header.Get("INFO")
	if description == "" {
		http.Error(w, "Missing description header", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO tasks(description) VALUES(?)", description)
	if err != nil {
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	task := Task{
		ID:          int(id),
		Description: description,
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

	description := r.Header.Get("INFO")
	if description == "" {
		http.Error(w, "Missing description header", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE tasks SET description=? WHERE id=?", description, id)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
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

	_, err = db.Exec("UPDATE tasks SET description='' WHERE id=?", id)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = ""
			break
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task description cleared successfully")
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

	var description string
	err = db.QueryRow("SELECT description FROM tasks WHERE id=?", id).Scan(&description)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task := Task{
		ID:          id,
		Description: description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}
