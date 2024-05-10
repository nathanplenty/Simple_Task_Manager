package database

import (
	"Simple_Task_Manager/pkg/domain"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB represents the SQLite database
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new instance of SQLiteDB
func NewSQLiteDB() *SQLiteDB {
	return &SQLiteDB{}
}

// Connect connects to the SQLite database
func (s *SQLiteDB) Connect(database *Database) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *SQLiteDB) CreateTask(task *domain.Task, user *domain.User) error {
	// Implement your logic to create a task in SQLite database here
	log.Println("Text for CreateTask")
	return nil
}

func (s *SQLiteDB) ReadTask(task *domain.Task, user *domain.User) error {
	// Implement your logic to create a task in SQLite database here
	return nil
}

func (s *SQLiteDB) UpdateTask(task *domain.Task, user *domain.User) error {
	// Implement your logic to create a task in SQLite database here
	return nil
}

func (s *SQLiteDB) DeleteTask(task *domain.Task, user *domain.User) error {
	// Implement your logic to create a task in SQLite database here
	return nil
}

func (s *SQLiteDB) CreateUser(user *domain.User) error {
	// Implement your logic to create a task in SQLite database here
	return nil
}

func (s *SQLiteDB) CheckUser(user *domain.User) error {
	// Implement your logic to create a task in SQLite database here
	return nil
}

//// CreateTask creates a new task in the SQLite database
//func (s *SQLiteDB) CreateTask(task *domain.Task, user *domain.User) (int, error) {
//	result, err := s.db.Exec("INSERT INTO tasks (task_name, due_date, completed, user_id) VALUES (?, ?, ?, ?)",
//		task.TaskName, task.DueDate, false, user.UserID)
//	if err != nil {
//		return 0, err
//	}
//	lastID, err := result.LastInsertId()
//	if err != nil {
//		return 0, err
//	}
//	return int(lastID), nil
//}
//
//// ReadTask reads a task from the SQLite database
//func (s *SQLiteDB) ReadTask(task *domain.Task, user *domain.User) (*domain.Task, error) {
//	var resultTask domain.Task
//	err := s.db.QueryRow("SELECT task_id, task_name, due_date, completed FROM tasks WHERE task_id=? AND user_id=?",
//		task.TaskID, user.UserID).Scan(&resultTask.TaskID, &resultTask.TaskName, &resultTask.DueDate, &resultTask.Completed)
//	if err != nil {
//		return nil, err
//	}
//	return &resultTask, nil
//}
//
//// UpdateTask updates a task in the SQLite database
//func (s *SQLiteDB) UpdateTask(task *domain.Task, user *domain.User) error {
//	_, err := s.db.Exec("UPDATE tasks SET task_name=?, due_date=?, completed=? WHERE task_id=? AND user_id=?",
//		task.TaskName, task.DueDate, task.Completed, task.TaskID, user.UserID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// DeleteTask deletes a task from the SQLite database
//func (s *SQLiteDB) DeleteTask(task *domain.Task, user *domain.User) error {
//	_, err := s.db.Exec("DELETE FROM tasks WHERE task_id=? AND user_id=?", task.TaskID, user.UserID)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// CreateUser creates a new user in the SQLite database
//func (s *SQLiteDB) CreateUser(user *domain.User) (int, error) {
//	result, err := s.db.Exec("INSERT INTO users (user_name, password) VALUES (?, ?)", user.UserName, user.Password)
//	if err != nil {
//		return 0, err
//	}
//	lastID, err := result.LastInsertId()
//	if err != nil {
//		return 0, err
//	}
//	return int(lastID), nil
//}
//
//// CheckUser checks if a user exists in the SQLite database
//func (s *SQLiteDB) CheckUser(user *domain.User) (int, error) {
//	var userID int
//	err := s.db.QueryRow("SELECT user_id FROM users WHERE user_name=? AND password=?", user.UserName, user.Password).Scan(&userID)
//	if errors.Is(err, sql.ErrNoRows) {
//		return 0, nil
//	} else if err != nil {
//		return 0, err
//	}
//	return userID, nil
//}
