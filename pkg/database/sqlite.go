package database

import (
	"Simple_Task_Manager/pkg/domain"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

// SQLiteDB represents the SQLite database
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new instance of SQLiteDB
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	log.Println("Start Function NewSQLiteDB")

	db := &SQLiteDB{}
	if err := db.Connect(dbPath); err != nil {
		return nil, err
	}
	if err := db.InitializeDatabase(dbPath); err != nil {
		return nil, err
	}
	return db, nil
}

// InitializeDatabase initializes the SQLite database
func (s *SQLiteDB) InitializeDatabase(dbPath string) error {

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		if err = s.createDatabaseTables(); err != nil {
			return err
		}
		log.Println("Database created successfully")
	} else {
		log.Println("Database reused successfully")
	}

	return nil
}

// createDatabaseTables creates the SQLite database
func (s *SQLiteDB) createDatabaseTables() error {
	log.Println("Start Function createDatabaseTables")

	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            user_id INTEGER PRIMARY KEY,
            user_name TEXT NOT NULL,
            password TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Error creating 'users' table: %v", err)
		return err
	}

	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS tasks (
            task_id INTEGER PRIMARY KEY,
            task_name TEXT NOT NULL,
            due_date TEXT NOT NULL,
            completed TEXT NOT NULL,
            user_id INTEGER,
            FOREIGN KEY (user_id) REFERENCES users (user_id)
        );
    `)
	if err != nil {
		log.Fatalf("Error creating 'tasks' table: %v", err)
		return err
	}

	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS sessions (
            session_id TEXT PRIMARY KEY,
            user_id INTEGER,
            expires_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (user_id)
        );
    `)
	if err != nil {
		log.Fatalf("Error creating 'sessions' table: %v", err)
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

// Connect establishes a connection to the SQLite database
func (s *SQLiteDB) Connect(dbPath string) error {
	log.Println("Start Function Connect")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

// CreateUser creates a new user in the SQLite database
func (s *SQLiteDB) CreateUser(user *domain.User) error {
	log.Println("Start Function CreateUser")

	_, err := s.db.Exec("INSERT INTO users (user_name, password) VALUES (?, ?)", user.UserName, user.Password)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
		return err
	}

	return nil
}

// CheckUser checks if the user exists in the SQLite database
func (s *SQLiteDB) CheckUser(user *domain.User) error {
	log.Println("Start Function CheckUser")

	row := s.db.QueryRow("SELECT user_name FROM users WHERE user_name = ?", user.UserName)
	err := row.Scan(&user.UserName)

	if err != nil {
		return errors.New("user does not exist")
	}

	return nil
}

// CreateTask creates a new task in the SQLite database
func (s *SQLiteDB) CreateTask(task *domain.Task, user *domain.User, session *domain.Session) error {
	log.Println("Start Function CreateTask")

	session, err := s.GetSession(session.SessionID)
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
		return err
	}

	if session == nil {
		log.Println("Session not found")
		return errors.New("session not found")
	}

	if strconv.Itoa(session.UserID) != user.UserID {
		log.Println("Invalid session")
		return errors.New("invalid session")
	}

	_, err = s.db.Exec("INSERT INTO tasks (task_name, due_date, completed, user_id) VALUES (?, ?, ?, ?)", task.TaskName, task.DueDate, false, session.UserID)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
		return err
	}

	return nil
}

// ReadTask reads a task from the SQLite database
func (s *SQLiteDB) ReadTask(task *domain.Task, user *domain.User, session *domain.Session) error {
	log.Println("Start Function ReadTask")

	session, err := s.GetSession(session.SessionID)
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
		return err
	}

	if session == nil {
		log.Println("Session not found")
		return errors.New("session not found")
	}

	if strconv.Itoa(session.UserID) != user.UserID {
		log.Println("Invalid session")
		return errors.New("invalid session")
	}

	row := s.db.QueryRow("SELECT task_id, task_name, due_date, completed FROM tasks WHERE task_id = ? AND user_id = ?", task.TaskID, session.UserID)
	err = row.Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed)
	if err != nil {
		log.Fatalf("Failed to read task: %v", err)
		return err
	}

	return nil
}

// UpdateTask updates a task in the SQLite database
func (s *SQLiteDB) UpdateTask(task *domain.Task, user *domain.User, session *domain.Session) error {
	log.Println("Start Function UpdateTask")

	session, err := s.GetSession(session.SessionID)
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
		return err
	}

	if session == nil {
		log.Println("Session not found")
		return errors.New("session not found")
	}

	if strconv.Itoa(session.UserID) != user.UserID {
		log.Println("Invalid session")
		return errors.New("invalid session")
	}

	_, err = s.db.Exec("UPDATE tasks SET task_name = ?, due_date = ?, completed = ? WHERE task_id = ? AND user_id = ?", task.TaskName, task.DueDate, task.Completed, task.TaskID, session.UserID)
	if err != nil {
		log.Fatalf("Failed to update task: %v", err)
		return err
	}

	return nil
}

// DeleteTask deletes a task from the SQLite database
func (s *SQLiteDB) DeleteTask(task *domain.Task, user *domain.User, session *domain.Session) error {
	log.Println("Start Function DeleteTask")

	session, err := s.GetSession(session.SessionID)
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
		return err
	}

	if session == nil {
		log.Println("Session not found")
		return errors.New("session not found")
	}

	if strconv.Itoa(session.UserID) != user.UserID {
		log.Println("Invalid session")
		return errors.New("invalid session")
	}

	_, err = s.db.Exec("DELETE FROM tasks WHERE task_id = ? AND user_id = ?", task.TaskID, session.UserID)
	if err != nil {
		log.Fatalf("Failed to delete task: %v", err)
		return err
	}

	return nil
}

// CreateSession creates a new session in the SQLite database
func (s *SQLiteDB) CreateSession(userID int) (string, error) {
	log.Println("Start Function CreateSession")

	sessionID := strconv.Itoa(userID) + "_session"
	_, err := s.db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?, ?)", sessionID, userID)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
		return "", err
	}

	return sessionID, nil
}

// UpdateSession updates an existing session in the SQLite database
func (s *SQLiteDB) UpdateSession(sessionID string) error {
	log.Println("Start Function UpdateSession")

	_, err := s.db.Exec("UPDATE sessions SET expires_at = CURRENT_TIMESTAMP WHERE session_id = ?", sessionID)
	if err != nil {
		log.Fatalf("Failed to update session: %v", err)
		return err
	}

	return nil
}

// DeleteSession deletes a session from the SQLite database
func (s *SQLiteDB) DeleteSession(sessionID string) error {
	log.Println("Start Function DeleteSession")

	_, err := s.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	if err != nil {
		log.Fatalf("Failed to delete session: %v", err)
		return err
	}

	return nil
}

// GetSession retrieves a session from the SQLite database
func (s *SQLiteDB) GetSession(sessionID string) (*domain.Session, error) {
	log.Println("Start Function GetSession")

	row := s.db.QueryRow("SELECT session_id, user_id, expires_at FROM sessions WHERE session_id = ?", sessionID)
	session := &domain.Session{}
	err := row.Scan(&session.SessionID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		log.Fatalf("Failed to get session: %v", err)
		return nil, err
	}

	return session, nil
}

// CheckPassword checks if the password for the user is correct
func (s *SQLiteDB) CheckPassword(user *domain.User) error {
	log.Println("Start Function CheckPassword")

	row := s.db.QueryRow("SELECT password FROM users WHERE user_name = ?", user.UserName)
	var storedPassword string
	err := row.Scan(&storedPassword)
	if err != nil {
		log.Fatalf("Failed to get password: %v", err)
		return err
	}

	if user.Password != storedPassword {
		return errors.New("invalid password")
	}

	return nil
}

// GetUserIDByUsername retrieves the user ID by username from the SQLite database
func (s *SQLiteDB) GetUserIDByUsername(userName string, userID *int) error {
	log.Println("Start Function GetUserIDByUsername")

	row := s.db.QueryRow("SELECT user_id FROM users WHERE user_name = ?", userName)
	err := row.Scan(userID)
	if err != nil {
		log.Fatalf("Failed to get user ID: %v", err)
		return err
	}

	return nil
}
