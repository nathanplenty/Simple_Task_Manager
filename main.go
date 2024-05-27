package main

import (
	"Simple_Task_Manager/pkg/database"
	"Simple_Task_Manager/pkg/server"
	"log"
)

func main() {
	// Database path - Manually change
	dbPath := "./database.db"

	// Database type - Manually change
	dbType := "sqlite"

	// Create a new database instance based on the specified dbType
	db, err := database.NewDatabase(database.SelectDatabaseType(dbType), dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Create a new server with the database instance
	srv := server.NewServer(db)
	srv.SetupRoutes()

	// Start the server
	if err = srv.StartServer(8080); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
