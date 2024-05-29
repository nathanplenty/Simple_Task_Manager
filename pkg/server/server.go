package server

import (
	"Simple_Task_Manager/pkg/database"
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Server represents the HTTP server
type Server struct {
	DB database.Database
}

func NewServer(db database.Database) *Server {
	return &Server{DB: db}
}

// SetupRoutes sets up the server routes
func (s *Server) SetupRoutes() {
	// Public routes
	http.HandleFunc("/users/login", s.loginUser)

	// Protected routes
	protectedRoutes := http.NewServeMux()
	protectedRoutes.HandleFunc("/users/create", s.createUser)
	protectedRoutes.HandleFunc("/tasks/create", s.createTask)
	protectedRoutes.HandleFunc("/tasks/read", s.readTask)
	protectedRoutes.HandleFunc("/tasks/update", s.updateTask)
	protectedRoutes.HandleFunc("/tasks/delete", s.deleteTask)

	http.Handle("/", s.JWTMiddleware(protectedRoutes))
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

// JWTMiddleware checks the JWT token.
func (s *Server) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token missing", http.StatusUnauthorized)
			return
		}

		token, err := s.validateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user information in the context
			ctx := context.WithValue(r.Context(), "username", claims["username"])
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
	})
}

func (s *Server) validateJWT(tokenString string) (*jwt.Token, error) {
	signingKey := []byte("SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	return token, err
}
