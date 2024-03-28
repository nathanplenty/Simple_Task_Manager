package main

import (
	taskdatabase "Simple_Task_Manager/database"
	taskhandler "Simple_Task_Manager/router"
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

	taskApp := &taskhandler.App{DB: app.DB, TaskManager: app}
	http.HandleFunc("/tasks", taskApp.HandleTasks)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
