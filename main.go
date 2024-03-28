package main

import (
	taskdatabase "Simple_Task_Manager/database"
	taskhandler "Simple_Task_Manager/router"
	taskmanager "Simple_Task_Manager/task_manager"
	"database/sql"
	"log"
	"net/http"
)

func main() {
	err := taskdatabase.NewSQLiteDB()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	db, err := sql.Open("sqlite3", "./database/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	app := &taskmanager.App{DB: db}
	if err := app.CheckDatabase(); err != nil {
		log.Fatalf("Error validating the database connection: %v", err)
	}

	taskApp := &taskhandler.App{TaskManager: app}

	http.HandleFunc("/tasks", taskApp.HandleTasks)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
