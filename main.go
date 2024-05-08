package main

import (
	"Task_Manager/mongodb"
	"Task_Manager/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	option := 2

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

	http.HandleFunc("/task/all", func(w http.ResponseWriter, r *http.Request) {
		controller.GetTaskAll(w, r)
	})
	http.HandleFunc("/tasks/add", controller.AddTask)
	http.HandleFunc("/tasks/delete", controller.DeleteTaskByID)
	http.HandleFunc("/tasks/update", controller.UpdateTaskByID)
	http.HandleFunc("/tasks/byID", controller.GetTaskByID)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupMongoDB() {
	collection := mongodb.NewCollection()
	controller := mongodb.NewController(collection)
	router := http.NewServeMux()

	router.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		controller.GetAllTasks(w)
	})
	router.HandleFunc("/tasks/add", controller.AddTask())
	router.HandleFunc("/tasks/delete", controller.DeleteTask)
	router.HandleFunc("/tasks/update", controller.UpdateTask)

	log.Fatal(http.ListenAndServe(":8080", router))
}
