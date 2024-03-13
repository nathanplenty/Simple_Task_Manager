package main

import (
	taskdatabase "Simple_Task_Manager/database"
	taskmanager "Simple_Task_Manager/task_manager"
	"log"
	"net/http"
)

func main() {
	if err := taskdatabase.CreateTables(); err != true {
		log.Fatalf("Fehler beim Erstellen der Datenbanktabellen: %v", err)
	}
	app := &taskmanager.App{}
	if err := app.InitDatabase(); err != true {
		log.Fatalf("Fehler beim Initialisieren der Datenbankverbindung: %v", err)
	}
	http.HandleFunc("/tasks", app.HandleTasks)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
