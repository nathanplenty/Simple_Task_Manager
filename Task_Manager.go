package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

var tasks []Task

func main() {
	http.HandleFunc("/tasks", handleTasks)
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	task := Task{
		ID:          len(tasks),
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
	if err != nil || id < 0 || id >= len(tasks) {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	description := r.Header.Get("INFO")
	if description == "" {
		http.Error(w, "Missing description header", http.StatusBadRequest)
		return
	}

	tasks[id].Description = description

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
	if err != nil || id < 0 || id >= len(tasks) {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tasks[id].Description = ""

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Task deleted successfully")
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

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
