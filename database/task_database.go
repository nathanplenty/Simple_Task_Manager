package taskdatabase

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DBManager interface {
	OpenDatabase() (*sql.DB, error)
	InitializeDatabase() error
}

type SQLiteDB struct {
	DatabasePath string
}

func NewSQLiteDB(databasePath string) *SQLiteDB {
	return &SQLiteDB{DatabasePath: databasePath}
}

func (db *SQLiteDB) OpenDatabase() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", db.DatabasePath)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return nil, err
	}
	return database, nil
}

func (db *SQLiteDB) InitializeDatabase() error {
	database, err := db.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}
	defer database.Close()

	_, err = database.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            user_id INTEGER PRIMARY KEY,
            user_name TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Error creating 'users' table: %v", err)
		return err
	}

	_, err = database.Exec(`
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

	log.Println("Database created successfully")
	return nil
}
