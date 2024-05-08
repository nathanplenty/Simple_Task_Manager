package database

import (
	"Simple_Task_Manager/pkg/domain"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB represents the MongoDB database
type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// NewMongoDB creates a new instance of MongoDB
func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

// Connect connects to the MongoDB database
func (m *MongoDB) Connect(connectionURI, dbName, collectionName string) error {
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}
	m.client = client
	m.database = client.Database(dbName)
	m.collection = m.database.Collection(collectionName)
	return nil
}

// CreateTask creates a new task in the MongoDB database
func (m *MongoDB) CreateTask(task *domain.Task, user *domain.User) error {
	userID, err := m.CheckUser(user)
	if err != nil {
		return err
	}
	user.UserID = userID
	task.Completed = "false"
	_, err = m.collection.InsertOne(context.Background(), task)
	if err != nil {
		return err
	}
	return nil
}

// ReadTask reads a task from the MongoDB database
func (m *MongoDB) ReadTask(taskID string, user *domain.User) (*domain.Task, error) {
	userID, err := m.CheckUser(user)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_id": userID, "_id": taskID}
	var task domain.Task
	err = m.collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask updates a task in the MongoDB database
func (m *MongoDB) UpdateTask(task *domain.Task, user *domain.User) error {
	userID, err := m.CheckUser(user)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": userID, "_id": task.TaskID}
	update := bson.M{"$set": bson.M{"task": task.TaskName, "due_date": task.DueDate, "completed": task.Completed}}

	_, err = m.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask deletes a task from the MongoDB database
func (m *MongoDB) DeleteTask(taskID string, user *domain.User) error {
	userID, err := m.CheckUser(user)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": userID, "_id": taskID}
	update := bson.M{"$unset": bson.M{"task": "", "due_date": "", "completed": ""}}

	_, err = m.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new user in the MongoDB database
func (m *MongoDB) CreateUser(user *domain.User) (string, error) {
	_, err := m.collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	return user.UserID, nil
}

// CheckUser checks if a user exists in the MongoDB database
func (m *MongoDB) CheckUser(user *domain.User) (string, error) {
	filter := bson.M{"user_name": user.UserName, "password": user.Password}
	var result domain.User
	err := m.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", nil
		}
		return "", err
	}
	return result.UserID, nil
}
