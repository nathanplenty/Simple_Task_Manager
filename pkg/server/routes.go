package server

import (
	"Simple_Task_Manager/pkg/database"
	"Simple_Task_Manager/pkg/domain"
	"encoding/json"
	"errors"
	"net/http"
)

// Server represents the HTTP server
type Server struct {
	DB database.Database
}

// NewServer creates a new instance of Server
func NewServer(dbType, dbPath string) (*Server, error) {
	var db database.Database
	var err error

	switch dbType {
	case "sqlite":
		db = database.NewSQLiteDB()
	case "mongodb":
		db = database.NewMongoDB()
	default:
		return nil, errors.New("unsupported database type")
	}

	err = db.Connect(dbPath)
	if err != nil {
		return nil, err
	}

	return &Server{DB: db}, nil
}

// SetupRoutes sets up the server routes
func (s *Server) SetupRoutes() {
	http.HandleFunc("/users/create", s.createUser)
	http.HandleFunc("/tasks/create", s.createTask)
	http.HandleFunc("/tasks/read", s.readTask)
	http.HandleFunc("/tasks/update", s.updateTask)
	http.HandleFunc("/tasks/delete", s.deleteTask)
}

// createUser handles the creation of a new user.
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	userID, err := s.DB.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"userID": userID})
}

// createTask handles the creation of a new task.
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user from request body", http.StatusBadRequest)
		return
	}

	err = s.DB.CreateTask(&task, &user)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// readTask handles the retrieval of a task.
func (s *Server) readTask(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user from request body", http.StatusBadRequest)
		return
	}

	readTask, err := s.DB.ReadTask(task.TaskID, &user)
	if err != nil {
		http.Error(w, "Failed to read task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(readTask)
}

// updateTask handles the update of a task.
func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user from request body", http.StatusBadRequest)
		return
	}

	err = s.DB.UpdateTask(&task, &user)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteTask handles the deletion of a task.
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	var user domain.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode user from request body", http.StatusBadRequest)
		return
	}

	err = s.DB.DeleteTask(task.TaskID, &user)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
