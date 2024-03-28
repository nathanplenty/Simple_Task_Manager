package main

import (
	taskdatabase "Simple_Task_Manager/database"
	taskhandler "Simple_Task_Manager/router"
	taskmanager "Simple_Task_Manager/task_manager"
	"log"
	"net/http"
)

func main() {
	dbManager := taskdatabase.NewSQLiteDB("./database/tasks.db")

	database, err := dbManager.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	defer database.Close()

	if err := dbManager.InitializeDatabase(); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	app := &taskmanager.App{DB: database}

	taskApp := &taskhandler.App{TaskManager: app}

	http.HandleFunc("/tasks", taskApp.HandleTasks)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
