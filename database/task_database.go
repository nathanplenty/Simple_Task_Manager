package taskdatabase

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func CreateTables() error {
	db, err := sql.Open("sqlite3", "./database/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			user_name TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating 'users' table: %v", err)
		return err
	}

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
		return err
	}

	log.Println("Database created successfully")
	return nil
}
