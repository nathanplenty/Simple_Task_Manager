package server

import (
	"Simple_Task_Manager/pkg/domain"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// logBody parses a JSON-encoded request body and logs its contents.
func logBody(w http.ResponseWriter, r *http.Request) (task domain.Task, user domain.User, err error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Failed to read request body: %v", err)
		return domain.Task{}, domain.User{}, err
	}

	var temp struct {
		Task domain.Task `json:"task"`
		User domain.User `json:"user"`
	}

	err = json.Unmarshal(bodyBytes, &temp)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return domain.Task{}, domain.User{}, err
	}

	log.Printf("Raw request body: %s", bodyBytes)

	return temp.Task, temp.User, nil
}

// encodeJSON encodes a value to JSON format and writes it to the response writer.
func encodeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Failed to encode data to JSON", http.StatusInternalServerError)
	}
}

// handleRequest handles related requests.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request, handler func(task domain.Task, user domain.User) error) {
	task, user, err := logBody(w, r)
	if err != nil {
		return
	}

	err = handler(task, user)
	if err != nil {
		http.Error(w, "Failed to handle request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encodeJSON(w, task)
}

// createUser handles the creation of a new user.
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.CreateUser(&user)
	})
}

// createTask handles the creation of a new task.
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.CreateTask(&task, &user)
	})
}

// readTask handles the retrieval of a task.
func (s *Server) readTask(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.ReadTask(&task, &user)
	})
}

// updateTask handles the update of a task.
func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.UpdateTask(&task, &user)
	})
}

// deleteTask handles the deletion of a task.
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.DeleteTask(&task, &user)
	})
}
