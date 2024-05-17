package database

import (
	"log"
)

// DBType is an enumeration of the supported database types
type DBType string

const (
	SQLiteDBType DBType = "sqlite"
	MongoDBType  DBType = "mongodb"
)

// CreateDatabase creates a new instance of the specified database type
func CreateDatabase(dbType DBType, connString string) (Database, error) {
	log.Println("Start Function CreateDatabase")
	switch dbType {
	case SQLiteDBType:
		return NewSQLiteDB(connString)
	//case MongoDBType:
	//	return NewMongoDB(connString) //Cannot use 'NewMongoDB(connString)' (type (*mongo.Database, error)) as the type DatabaseType does not implement 'Database' as some methods are missing:Connect(dbPath string) errorCreateUser(user *domain.User) errorCheckUser(user *domain.User) error…
	default:
		return nil, nil
	}
}
