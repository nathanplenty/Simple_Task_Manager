package server

import (
	"Simple_Task_Manager/pkg/database"
	"log"
	"net/http"
	"strconv"
)

// Server represents the HTTP server
type Server struct {
	DB database.Database
}

// NewServer creates a new instance of Server
func NewServer(db database.Database) *Server {
	return &Server{DB: db}
}

// SetupRoutes sets up the server routes
func (s *Server) SetupRoutes() {
	http.HandleFunc("/users/create", s.createUser)
	http.HandleFunc("/tasks/create", s.createTask)
	http.HandleFunc("/tasks/read", s.readTask)
	http.HandleFunc("/tasks/update", s.updateTask)
	http.HandleFunc("/tasks/delete", s.deleteTask)
}

// StartServer starts the HTTP server on the specified port
func (s *Server) StartServer(port int) error {
	addr := ":" + strconv.Itoa(port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server on port %d: %v", port, err)
	}
	return nil
}
