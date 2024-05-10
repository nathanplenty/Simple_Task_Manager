package server

import (
	"Simple_Task_Manager/pkg/domain"
	"encoding/json"
	"net/http"
)

// createUser handles the creation of a new user.
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = s.DB.CreateUser(&user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		return
	}
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
	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task to JSON", http.StatusInternalServerError)
		return
	}
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

	err = s.DB.ReadTask(&task, &user)
	if err != nil {
		http.Error(w, "Failed to read task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task to JSON", http.StatusInternalServerError)
		return
	}
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
	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task to JSON", http.StatusInternalServerError)
		return
	}
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

	err = s.DB.DeleteTask(&task, &user)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task to JSON", http.StatusInternalServerError)
		return
	}
}
