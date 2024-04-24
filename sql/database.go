package sql

import "database/sql"

type Database struct {
	db *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) Initialize() error {
	_, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY,
		task TEXT,
		due_date DATE,
		completed BOOLEAN
	)`)
	return err
}
