package server

import (
	"Simple_Task_Manager/pkg/domain"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// logBody parses a JSON-encoded request body.
func logBody(w http.ResponseWriter, r *http.Request) (task domain.Task, user domain.User, err error) {
	log.Println("Start Function router/logBody")
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Failed to read request body: %v", err)
		return domain.Task{}, domain.User{}, err
	}

	var bodyData map[string]interface{}
	if err = json.Unmarshal(bodyBytes, &bodyData); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return domain.Task{}, domain.User{}, err
	}

	user = domain.User{}
	task = domain.Task{}

	if val, ok := bodyData["user_id"].(string); ok {
		user.UserID = val
	}
	if val, ok := bodyData["user_name"].(string); ok {
		user.UserName = val
	}
	if val, ok := bodyData["password"].(string); ok {
		user.Password = val
	}
	if val, ok := bodyData["task_id"].(string); ok {
		task.TaskID = val
	}
	if val, ok := bodyData["task_name"].(string); ok {
		task.TaskName = val
	}
	if val, ok := bodyData["due_date"].(string); ok {
		task.DueDate = val
	}
	if val, ok := bodyData["completed"].(string); ok {
		task.Completed = val
	}

	return task, user, nil
}

// encodeJSON encodes a value to JSON format and writes it to the response writer.
func encodeJSON(w http.ResponseWriter, v interface{}) {
	log.Println("Start Function router/endcodeJSON")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Failed to encode data to JSON", http.StatusInternalServerError)
	}
}

// handleRequestWithResult handles related requests and returns a result.
func (s *Server) handleRequestWithResult(w http.ResponseWriter, r *http.Request, handler func(task domain.Task, user domain.User) (interface{}, error)) {
	log.Println("Start Function router/handleRequestWithResult")

	task, user, err := logBody(w, r)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	result, err := handler(task, user)
	if err != nil {
		http.Error(w, "Failed to handle request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encodeJSON(w, result)
}

// handleRequest handles related requests and does not return a result.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request, handler func(task domain.Task, user domain.User) error) {
	log.Println("Start Function router/handleRequest")

	task, user, err := logBody(w, r)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = handler(task, user)
	if err != nil {
		http.Error(w, "Failed to handle request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encodeJSON(w, "Request successfully executed")
}

// createUser handles the creation of a new user.
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function router/createUser")
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		newUser := domain.NewUser(user.UserID, user.UserName, user.Password)
		return s.DB.CreateUser(newUser)
	})
}

// loginUser handles the login of a user.
func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function router/loginUser")
	s.handleRequestWithResult(w, r, func(task domain.Task, user domain.User) (interface{}, error) {
		getUser := domain.NewUser(user.UserID, user.UserName, user.Password)
		token, err := s.DB.LoginUser(getUser)
		if err != nil {
			return nil, err
		}
		return map[string]string{"token": token}, nil
	})
}

// createTask handles the creation of a new task.
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function router/createTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		newTask := domain.NewTask(task.TaskID, task.TaskName, task.DueDate, task.Completed)
		return s.DB.CreateTask(newTask, &user)
	})
}

// readTask handles the retrieval of a task.
func (s *Server) readTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function router/readTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.ReadTask(&task, &user)
	})
}

// updateTask handles the update of a task.
func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function router/updateTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.UpdateTask(&task, &user)
	})
}

// deleteTask handles the deletion of a task.
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function router/deleteTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User) error {
		return s.DB.DeleteTask(&task, &user)
	})
}
