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

// flexibel den Code anpassen > Ziel für weitere Entwicklung
// repetitive Arbeit soll vermieden werden
// Abhängigkeiten reduzieren innerhalb vom Code
// gibt es eine bessere Lösung, als switch um nach Datenbank zu agieren? > Interface soll mir die Antwort zu dieser Frage geben
// json = bson = xml > struct kann besser "marched" werden
// !!! WORK1: Interface einbauen, keine verschiedenen steckdosen bauen (SQLite Toaster, MongoDB Kühlschrank) (Objekte wie in Python) !!!
// UI ungleich Interface! Interface kann UI sein, aber als verständnis als objekt speichern
// WORK2: Auf MongoDB umstellen > wechsel soll fliegend passieren
// (WORK3: Pattern lernen, versuch Dinge miteinander kompatible zu machen)
// (WORK4: UnitTest "mocking", Werte selbst festlegen)
// Prinzipien nutzen, um guten Code zu bauen > für skalierbare / flexible Projekte Interface nutzen
// Fehler "fangen" verstehen:
// Fehler sollten so weit außen wie möglich aber so weit innen wie nötig behandelt werden.
// Der Use-Case entscheidet über die Fehlerbehandlung - nicht die Kernlogik.
// .log error values kennenlernen
