package database

import (
	"errors"
)

// SelectDatabaseType defines the supported database types
type SelectDatabaseType string

const (
	SQLite SelectDatabaseType = "sqlite"
	Mongo  SelectDatabaseType = "mongodb"
)

// NewDatabase creates a new database instance based on the specified type
func NewDatabase(dbType SelectDatabaseType, dbPath string) (Database, error) {
	switch dbType {
	case SQLite:
		return NewSQLiteDB(dbPath)
	case Mongo:
		return NewMongoDB(), nil
	default:
		return nil, errors.New("unsupported database type: " + string(dbType))
	}
}

//var drivers map[string]func() Database
//drivers["mongo"] = 7
//dbFunc, ok := drivers["mongo"]
//return dbFunc()
