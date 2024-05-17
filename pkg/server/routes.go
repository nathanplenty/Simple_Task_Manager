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
	log.Println("Start Function logBody")
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
	log.Println("Start Function encodeJSON")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		http.Error(w, "Failed to encode data to JSON", http.StatusInternalServerError)
	}
}

// handleRequest handles related requests.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request, handler func(task domain.Task, user domain.User, session *domain.Session) error) {
	log.Println("Start Function handleRequest")

	task, user, err := logBody(w, r)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	session := s.GetSession(r) //!!!STUCK!!!
	if session == nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	err = handler(task, user, session)
	if err != nil {
		log.Printf("Failed to handle request: %v\n", err)
		http.Error(w, "Failed to handle request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encodeJSON(w, "Request successfully executed")
}

// createUser handles the creation of a new user.
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function createUser")
	s.handleRequest(w, r, func(task domain.Task, user domain.User, session *domain.Session) error {
		newUser := domain.NewUser(user.UserID, user.UserName, user.Password)
		return s.DB.CreateUser(newUser)
	})
}

// createTask handles the creation of a new task.
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function createTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User, session *domain.Session) error {
		newTask := domain.NewTask(task.TaskID, task.TaskName, task.DueDate, task.Completed)
		return s.DB.CreateTask(newTask, &user, session)
	})
}

// readTask handles the retrieval of a task.
func (s *Server) readTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function readTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User, session *domain.Session) error {
		return s.DB.ReadTask(&task, &user, session)
	})
}

// updateTask handles the update of a task.
func (s *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function updateTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User, session *domain.Session) error {
		return s.DB.UpdateTask(&task, &user, session)
	})
}

// deleteTask handles the deletion of a task.
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function deleteTask")
	s.handleRequest(w, r, func(task domain.Task, user domain.User, session *domain.Session) error {
		return s.DB.DeleteTask(&task, &user, session)
	})
}

// loginUser handles the login of a user and creates a session.
func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Start Function loginUser")

	_, user, err := logBody(w, r)
	if err != nil {
		return
	}

	err = s.DB.CheckPassword(&user)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var userID int
	err = s.DB.GetUserIDByUsername(user.UserName, &userID)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusInternalServerError)
		return
	}

	sessionID, err := s.DB.CreateSession(userID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encodeJSON(w, map[string]string{"session_id": sessionID})
}
