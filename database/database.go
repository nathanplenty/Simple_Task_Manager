package taskdatabase

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//import (
//	"log"
//	mgo "gopkg.in/mgo.v2"
//)

func CreateTables() error {
	db, err := sql.Open("sqlite3", "./database/tasks.db")
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			user_name TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Error creating 'users' table: %v", err)
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			task_id INTEGER PRIMARY KEY,
			user_id INTEGER,
			task_name TEXT NOT NULL,
			due_date DATE,
			completed BOOLEAN,
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
	`)
	if err != nil {
		log.Fatalf("Error creating 'tasks' table: %v", err)
		return err
	}

	log.Println("Database created successfully")
	return nil
}

//func CreateTables() error {
//	session, err := mgo.Dial("localhost")
//// check connection
//	db := session.DB("tasks")
//
//	err = db.C("users").EnsureIndex(mgo.Index{
//		Key:    []string{"user_id"},
//		Unique: true,
//	})
//// check users container
//	err = db.C("tasks").EnsureIndex(mgo.Index{
//		Key:    []string{"task_id"},
//		Unique: true,
//	})
//// check tasks container
//// finish end of function
//}
