package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			due_date DATE,
			completed BOOLEAN
		);
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	log.Println("Database schema created successfully")
}
