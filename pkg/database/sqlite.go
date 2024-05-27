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
	log.Println("Start Function sqlite/NewSQLiteDB")

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
	log.Println("Start Function sqlite/createDatabaseTables")

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
	log.Println("Start Function sqlite/Connect")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}
	s.db = db
	return nil
}

// CreateUser creates a new user in the SQLite database and returns the new user_id
func (s *SQLiteDB) CreateUser(user *domain.User) error {
	log.Println("Start Function sqlite/CreateUser")

	err := s.CheckUser(user)
	if err == nil {
		log.Println("User already exists")
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

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return err
	}

	userIDInt := int(lastInsertID)

	log.Printf("New user created with ID: %d", userIDInt)
	return nil
}

// LoginUser authenticates a user by checking username and password
func (s *SQLiteDB) LoginUser(user *domain.User) error {
	log.Println("Start Function sqlite/LoginUser")

	err := s.CheckUser(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	err = s.CheckPassword(user)
	if err != nil {
		if errors.Is(err, errors.New("incorrect password")) {
			return errors.New("wrong password")
		}
		return err
	}

	log.Println("User logged in")
	return nil
}

// CheckUser checks if a user exists in the SQLite database
func (s *SQLiteDB) CheckUser(user *domain.User) error {
	log.Println("Start Function sqlite/CheckUser")

	var userID int
	err := s.db.QueryRow("SELECT user_id FROM users WHERE user_name=?", user.UserName).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

// CheckPassword checks if a user's password is correct in the database
func (s *SQLiteDB) CheckPassword(user *domain.User) error {
	log.Println("Start Function sqlite/CheckPassword")

	var storedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE user_name=?", user.UserName).Scan(&storedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}
	if storedPassword != user.Password {
		return errors.New("incorrect password")
	}
	return nil
}

// CreateTask creates a new task in the SQLite database
func (s *SQLiteDB) CreateTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function sqlite/CreateTask")

	if err := s.CheckUser(user); err != nil {
		return err
	}

	var userID int
	err := s.db.QueryRow("SELECT user_id FROM users WHERE user_name=?", user.UserName).Scan(&userID)
	if err != nil {
		log.Println("Error querying the database for user ID:", err)
		return err
	}

	if err = s.CheckPassword(user); err != nil {
		return err
	}

	var existingTaskID int
	err = s.db.QueryRow("SELECT task_id FROM tasks WHERE task_name=? AND due_date=? AND user_id=?",
		task.TaskName, task.DueDate, userID).Scan(&existingTaskID)
	if err == nil {
		log.Println("Task already exists")
		return errors.New("task already exists")
	} else if !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error querying the database for existing task:", err)
		return err
	}

	result, err := s.db.Exec("INSERT INTO tasks (task_name, due_date, completed, user_id) VALUES (?, ?, ?, ?)",
		task.TaskName, task.DueDate, false, userID)
	if err != nil {
		return err
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	log.Println("New task created with ID:", taskID, ", with userID:", userID)
	return nil
}

// ReadTask reads a task from the SQLite database
func (s *SQLiteDB) ReadTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function sqlite/ReadTask")

	if err := s.CheckUser(user); err != nil {
		return err
	}

	if err := s.CheckPassword(user); err != nil {
		return err
	}

	err := s.db.QueryRow("SELECT task_id, task_name, due_date, completed FROM tasks WHERE task_id=? AND user_id=(SELECT user_id FROM users WHERE user_name = ?)",
		task.TaskID, user.UserName).Scan(&task.TaskID, &task.TaskName, &task.DueDate, &task.Completed)
	return err
}

// UpdateTask updates a task in the SQLite database
func (s *SQLiteDB) UpdateTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function sqlite/UpdateTask")

	if err := s.CheckUser(user); err != nil {
		return err
	}

	var userID int
	err := s.db.QueryRow("SELECT user_id FROM users WHERE user_name=? AND password=?", user.UserName, user.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No matching user found!")
			return err
		}
		log.Println("Error querying the database for user ID:", err)
		return err
	}

	var taskID int
	err = s.db.QueryRow("SELECT task_id FROM tasks WHERE task_name=? AND user_id=?", task.TaskName, userID).Scan(&taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No matching task found for the user!")
			return err
		}
		log.Println("Error querying the database for task ID:", err)
		return err
	}

	log.Printf("Task before update: %+v\n", task)
	log.Printf("User before update: %+v\n", user)

	completed, err := strconv.ParseBool(task.Completed)
	if err != nil {
		log.Println("Failed to parse 'completed' as boolean:", err)
		return err
	}

	result, err := s.db.Exec("UPDATE tasks SET task_name=?, due_date=?, completed=? WHERE task_id=? AND user_id=?", task.TaskName, task.DueDate, completed, taskID, userID)
	if err != nil {
		log.Println("Error updating task:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected > 0 {
		log.Printf("Task updated for User ID: %d, Task ID: %d, Rows affected: %d", userID, taskID, rowsAffected)
	} else {
		log.Printf("No task updated for User ID: %d, Task ID: %d", userID, taskID)
	}
	return nil
}

// DeleteTask anonymizes a task in the SQLite database instead of deleting it
func (s *SQLiteDB) DeleteTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function sqlite/DeleteTask")

	if err := s.CheckUser(user); err != nil {
		return err
	}

	var userID int
	err := s.db.QueryRow("SELECT user_id FROM users WHERE user_name=? AND password=?", user.UserName, user.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No matching user found!")
			return err
		}
		log.Println("Error querying the database for user ID:", err)
		return err
	}

	result, err := s.db.Exec("UPDATE tasks SET task_name='(deleted)', due_date='0001-01-01', completed=false WHERE task_name=? AND user_id=?", task.TaskName, userID)
	if err != nil {
		log.Println("Error updating task:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected > 0 {
		log.Printf("Task anonymized for User ID: %d, Task Name: %s, Rows affected: %d", userID, task.TaskName, rowsAffected)
	} else {
		log.Printf("No task anonymized for User ID: %d, Task Name: %s", userID, task.TaskName)
	}
	return nil
}
