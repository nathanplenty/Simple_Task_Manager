package databaseMongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DBManager interface {
	OpenDatabase() (*mongo.Client, error)
	InitializeDatabase(client *mongo.Client) error
}

type MongoDB struct {
	ConnectionString string
}

func NewMongoDB(connectionString string) *MongoDB {
	return &MongoDB{ConnectionString: connectionString}
}

func (db *MongoDB) OpenDatabase() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db.ConnectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MongoDB!")
	return client, nil
}

func (db *MongoDB) InitializeDatabase(client *mongo.Client) error {
	database := client.Database("mongo_db")

	collections := []struct {
		Name   string
		Schema interface{}
	}{
		{
			Name: "users",
			Schema: bson.D{
				{"user_id", 0},
				{"user_name", ""},
			},
		},
		{
			Name: "tasks",
			Schema: bson.D{
				{"task_id", 0},
				{"user_id", 0},
				{"task_name", ""},
				{"due_date", time.Now()},
				{"completed", false},
			},
		},
	}

	for _, col := range collections {
		_, err := database.Collection(col.Name).InsertOne(context.Background(), col.Schema)
		if err != nil {
			return err
		}
	}

	log.Println("Database initialized successfully!")
	return nil
}
