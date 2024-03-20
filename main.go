package main

import (
	taskdatabase "Simple_Task_Manager/database"
	taskmanager "Simple_Task_Manager/task_manager"
	"log"
	"net/http"
)

func main() {
	err := taskdatabase.CreateTables()
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	app := &taskmanager.App{}
	err = app.CheckDatabase()
	if err != nil {
		log.Fatalf("Error validating the database connection: %v", err)
	}

	http.HandleFunc("/tasks", app.HandleTasks)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// no changes for MongoDB
