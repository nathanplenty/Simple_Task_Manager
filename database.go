package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func createTables() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			task_id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			due_date DATE,
			completed BOOLEAN
		);
	`)
	if err != nil {
		log.Fatalf("Error creating 'tasks' table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			username TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating 'users' table: %v", err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS task_assignments (
        task_id INTEGER PRIMARY KEY,
        user_id INTEGER,
        FOREIGN KEY (task_id) REFERENCES tasks(task_id),
        FOREIGN KEY (user_id) REFERENCES users(user_id)
    );
`)
	if err != nil {
		log.Fatalf("Error creating 'task_assignments' table: %v", err)
	}

	log.Println("Database schema created successfully")
}

func main() {
	createTables()
}
