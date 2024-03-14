package taskdatabase

// Adjust import for MongoDB
import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func CreateTables() bool {
	// Adjust connection for MongoDB, instead of a lite database we want a client
	db, err := sql.Open("sqlite3", "./database/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return false
	}
	defer db.Close()
	// Adjust table for MongoDB, instead of create table we want a new collection
	// Adjust columns for MongoDB, instead of creating columns we want models
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			user_name TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating 'users' table: %v", err)
		return false
	}
	// Adjust table for MongoDB, instead of create table we want a new collection
	// Adjust columns for MongoDB, instead of creating columns we want models
	_, err = db.Exec(`
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
		return false
	}
	log.Println("Database schema created successfully")
	return true
}
