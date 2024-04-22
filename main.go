package main

import (
	databaseMongodb "Simple_Task_Manager/database/mongodb"
	databaseSqlite "Simple_Task_Manager/database/sqlite"
	routerMongodb "Simple_Task_Manager/router/mongodb"
	routerSqlite "Simple_Task_Manager/router/sqlite"
	taskManagerMongodb "Simple_Task_Manager/task_manager/sqlite"
	taskManagerSqlite "Simple_Task_Manager/task_manager/sqlite"
	"log"
	"net/http"
)

func main() {
	const databaseType = 2

	switch databaseType {
	case 1:
		sqlite()
	case 2:
		mongodb()
	default:
		log.Println("Invalid database type. Using SQLite as default.")
		sqlite()
	}
}

func sqlite() {
	dbManager := databaseSqlite.NewSQLiteDB("./database/sqlite/sqlite.db")

	database, err := dbManager.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	defer database.Close()

	if err = dbManager.InitializeDatabase(); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	app := &taskManagerSqlite.App{DB: database}

	taskApp := &routerSqlite.App{TaskManager: app}

	http.HandleFunc("/tasks", taskApp.HandleTasks)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mongodb() {
	dbManager := databaseMongodb.NewMongoDB("mongodb://localhost:27017")

	database, err := dbManager.OpenDatabase()
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	if err = dbManager.InitializeDatabase(database); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	app := &taskManagerMongodb.App{DB: nil}

	taskApp := &routerMongodb.App{TaskManager: app}

	http.HandleFunc("/tasks", taskApp.HandleTasks)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
