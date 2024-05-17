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
	http.HandleFunc("/createUser", s.createUser)
	http.HandleFunc("/createTask", s.createTask)
	http.HandleFunc("/readTask", s.readTask)
	http.HandleFunc("/updateTask", s.updateTask)
	http.HandleFunc("/deleteTask", s.deleteTask)
	http.HandleFunc("/loginUser", s.loginUser)
}

// StartServer starts the HTTP server
func (s *Server) StartServer(port int) error {
	addr := ":" + strconv.Itoa(port)
	log.Printf("Server listening on port %d...\n", port)
	return http.ListenAndServe(addr, nil)
}
