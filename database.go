package main

//goland:noinspection ALL,Annotator,Annotator
import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY,
			description TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database schema created successfully")
}
