package database

import (
	"Simple_Task_Manager/pkg/domain"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

// SQLiteDB represents the SQLite database
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new instance of SQLiteDB
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
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
            user_id INTEGER,
            task_name TEXT NOT NULL,
            due_date DATE,
            completed BOOLEAN,
            FOREIGN KEY (user_id) REFERENCES users(user_id)
        );
    `)
	if err != nil {
		log.Fatalf("Error creating 'tasks' table: %v", err)
		return err
	}

	return nil
}

// Connect connects to the SQLite database
func (s *SQLiteDB) Connect(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}
	s.db = db
	return nil
}

// CreateUser creates a new user in the SQLite database if the user does not already exist
func (s *SQLiteDB) CreateUser(user *domain.User) error {
	err := s.CheckUser(user)
	if err == nil {
		return errors.New("user already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Error checking user existence: %v", err)
		return err
	}

	result, err := s.db.Exec("INSERT INTO users (user_name, password) VALUES (?, ?)", user.UserName, user.Password)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Rows affected after creating user: %d", rowsAffected)
	return nil
}

// CheckUser checks if a user exists in the SQLite database
func (s *SQLiteDB) CheckUser(user *domain.User) error {
	var userID int
	err := s.db.QueryRow("SELECT user_id FROM users WHERE user_name=? AND password=?", user.UserName, user.Password).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

// CreateTask creates a new task in the SQLite database
func (s *SQLiteDB) CreateTask(task *domain.Task, user *domain.User) error {
	_, err := s.db.Exec("INSERT INTO tasks (task_name, due_date, completed, user_id) VALUES (?, ?, ?, ?)",
		task.TaskName, task.DueDate, false, user.UserID)
	return err
}

// ReadTask liest eine Aufgabe aus der SQLite-Datenbank
func (s *SQLiteDB) ReadTask(task *domain.Task, user *domain.User) error {
	err := s.db.QueryRow("SELECT task_id, task_name, due_date, completed FROM tasks WHERE task_id=? AND user_id=?",
		task.TaskID, user.UserID).Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed)
	return err
}

// UpdateTask aktualisiert eine Aufgabe in der SQLite-Datenbank
func (s *SQLiteDB) UpdateTask(task *domain.Task, user *domain.User) error {
	_, err := s.db.Exec("UPDATE tasks SET task_name=?, due_date=?, completed=? WHERE task_id=? AND user_id=?",
		task.TaskName, task.DueDate, task.Completed, task.TaskID, user.UserID)
	return err
}

// DeleteTask löscht eine Aufgabe aus der SQLite-Datenbank
func (s *SQLiteDB) DeleteTask(task *domain.Task, user *domain.User) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE task_id=? AND user_id=?", task.TaskID, user.UserID)
	return err
}
