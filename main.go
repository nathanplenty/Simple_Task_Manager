package main

import (
	"Simple_Task_Manager/pkg/database"
	"Simple_Task_Manager/pkg/server"
	"log"
	"net/http"
)

func main() {
	dbType := "sqlite"

	db, err := database.NewDatabase(dbType)
	if err != nil {
		log.Fatalf("Failed to create %s database instance: %v", dbType, err)
	}

	srv := server.NewServer(db)
	srv.SetupRoutes()

	log.Println("Server is running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
