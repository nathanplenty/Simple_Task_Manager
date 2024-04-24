package databaseMongoDB

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBManager interface {
	OpenDatabase() (*mongo.Database, error)
	InitializeDatabase() error
}

type MongoDB struct {
	ConnectionString string
	DatabaseName     string
}

func NewMongoDB(connectionString, databaseName string) *MongoDB {
	return &MongoDB{
		ConnectionString: connectionString,
		DatabaseName:     databaseName,
	}
}

func (db *MongoDB) OpenDatabase() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.ConnectionString))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
		return nil, err
	}

	database := client.Database(db.DatabaseName)
	return database, nil
}

func (db *MongoDB) InitializeDatabase() error {
	log.Println("MongoDB initialized successfully")
	return nil
}
