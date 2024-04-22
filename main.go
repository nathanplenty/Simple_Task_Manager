package main

import (
	databaseSqlite "Simple_Task_Manager/database"
	routerHandler "Simple_Task_Manager/router"
	taskManager "Simple_Task_Manager/task_manager"
	"log"
	"net/http"
)

func main() {
	dbManager := databaseSqlite.NewSQLiteDB("./database/tasks.db")

	database, err := dbManager.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	defer database.Close()

	if err := dbManager.InitializeDatabase(); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	app := &taskManager.App{DB: database}

	taskApp := &routerHandler.App{TaskManager: app}

	http.HandleFunc("/tasks", taskApp.HandleTasks)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
