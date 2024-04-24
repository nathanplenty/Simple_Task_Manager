package main

import (
	"Task_Manager/mongodb"
	"Task_Manager/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	option := 1

	switch option {
	case 1:
		setupSQL()
	case 2:
		setupMongoDB()
	default:
		setupSQL()
	}
}

func setupSQL() {
	db, err := sql.NewDatabase("sql/tasks.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	err = db.Initialize()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}

	controller := sql.NewController(db)

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		controller.GetAllTasks(w)
	})
	http.HandleFunc("/tasks/add", controller.AddTask)
	http.HandleFunc("/tasks/delete", controller.DeleteTask)
	http.HandleFunc("/tasks/update", controller.UpdateTask)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupMongoDB() {
	collection := mongodb.NewCollection()
	controller := mongodb.NewController(collection)

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		controller.GetAllTasks(w)
	})
	http.HandleFunc("/tasks/add", controller.AddTask)
	http.HandleFunc("/tasks/delete", controller.DeleteTask)
	http.HandleFunc("/tasks/update", controller.UpdateTask)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
