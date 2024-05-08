package server

import (
	"Simple_Task_Manager/pkg/database"
	"log"
	"net/http"
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

// Run starts the server
func (s *Server) Run() {
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
