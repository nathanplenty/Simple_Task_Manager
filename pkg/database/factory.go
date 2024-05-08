package database

import (
	"errors"
)

// DatabaseType defines the supported database types
type DatabaseType string

const (
	SQLite  DatabaseType = "sqlite"
	MongoDB DatabaseType = "mongodb"
)

// NewDatabase creates a new database instance based on the specified type
func NewDatabase(dbType DatabaseType, dbPath string) (Database, error) {
	switch dbType {
	case SQLite:
		return NewSQLiteDB(dbPath), nil
	case MongoDB:
		return NewMongoDB(dbPath), nil
	default:
		return nil, errors.New("unsupported database type")
	}
}
