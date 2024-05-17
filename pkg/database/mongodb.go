package database

import (
	"Simple_Task_Manager/pkg/domain"
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

//// NewMongoDB creates a new instance of MongoDB
//func NewMongoDB(connString string) (Database, error) {
//	log.Println("Start Function NewMongoDB")
//
//	clientOptions := options.Client().ApplyURI(connString)
//	client, err := mongo.Connect(context.Background(), clientOptions)
//	if err != nil {
//		return nil, err
//	}
//
//	err = client.Ping(context.Background(), nil)
//	if err != nil {
//		return nil, err
//	}
//
//	database := client.Database("task_manager")
//
//	return database, nil
//}

// InitializeDatabase initializes the MongoDB database with the necessary collections and indexes
func (m *MongoDB) InitializeDatabase() error {
	log.Println("Start Function InitializeDatabase")

	// Ensure indexes for collections
	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_name", Value: 1}},
		},
	}
	taskIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "task_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
	}

	_, err := m.database.Collection("users").Indexes().CreateMany(context.Background(), userIndexes)
	if err != nil {
		return err
	}

	_, err = m.database.Collection("tasks").Indexes().CreateMany(context.Background(), taskIndexes)
	if err != nil {
		return err
	}

	sessionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "session_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
		},
	}

	_, err = m.database.Collection("sessions").Indexes().CreateMany(context.Background(), sessionIndexes)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new user in the MongoDB database
func (m *MongoDB) CreateUser(user *domain.User) error {
	log.Println("Start Function CreateUser")

	_, err := m.database.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}

// CheckUser checks if the user exists in the MongoDB database
func (m *MongoDB) CheckUser(user *domain.User) error {
	log.Println("Start Function CheckUser")

	filter := bson.M{"user_name": user.UserName}
	err := m.database.Collection("users").FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		return errors.New("user does not exist")
	}

	return nil
}

// CreateTask creates a new task in the MongoDB database
func (m *MongoDB) CreateTask(task *domain.Task) error {
	log.Println("Start Function CreateTask")

	_, err := m.database.Collection("tasks").InsertOne(context.Background(), task)
	if err != nil {
		return err
	}

	return nil
}

// ReadTask reads a task from the MongoDB database
func (m *MongoDB) ReadTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function ReadTask")

	filter := bson.M{"task_id": task.TaskID, "user_id": user.UserID}
	err := m.database.Collection("tasks").FindOne(context.Background(), filter).Decode(task)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTask updates a task in the MongoDB database
func (m *MongoDB) UpdateTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function UpdateTask")

	filter := bson.M{"task_id": task.TaskID, "user_id": user.UserID}
	update := bson.M{
		"$set": bson.M{
			"task_name": task.TaskName,
			"due_date":  task.DueDate,
			"completed": task.Completed,
		},
	}
	_, err := m.database.Collection("tasks").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask deletes a task from the MongoDB database
func (m *MongoDB) DeleteTask(task *domain.Task, user *domain.User) error {
	log.Println("Start Function DeleteTask")

	filter := bson.M{"task_id": task.TaskID, "user_id": user.UserID}
	_, err := m.database.Collection("tasks").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

// CreateSession creates a new session in the MongoDB database
func (m *MongoDB) CreateSession(userID int) (string, error) {
	log.Println("Start Function CreateSession")

	sessionID := strconv.Itoa(userID) + "_session"
	session := &domain.Session{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Session expires in 24 hours
	}
	_, err := m.database.Collection("sessions").InsertOne(context.Background(), session)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// UpdateSession updates an existing session in the MongoDB database
func (m *MongoDB) UpdateSession(sessionID string) error {
	log.Println("Start Function UpdateSession")

	filter := bson.M{"session_id": sessionID}
	update := bson.M{
		"$set": bson.M{
			"expires_at": time.Now().Add(24 * time.Hour),
		},
	}
	_, err := m.database.Collection("sessions").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSession deletes a session from the MongoDB database
func (m *MongoDB) DeleteSession(sessionID string) error {
	log.Println("Start Function DeleteSession")

	filter := bson.M{"session_id": sessionID}
	_, err := m.database.Collection("sessions").DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}

// GetSession retrieves a session from the MongoDB database
func (m *MongoDB) GetSession(sessionID string) (*domain.Session, error) {
	log.Println("Start Function GetSession")

	filter := bson.M{"session_id": sessionID}
	session := &domain.Session{}
	err := m.database.Collection("sessions").FindOne(context.Background(), filter).Decode(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// CheckPassword checks if the password for the user is correct
func (m *MongoDB) CheckPassword(user *domain.User) error {
	log.Println("Start Function CheckPassword")

	filter := bson.M{"user_name": user.UserName}
	var storedUser domain.User
	err := m.database.Collection("users").FindOne(context.Background(), filter).Decode(&storedUser)
	if err != nil {
		return err
	}

	if user.Password != storedUser.Password {
		return errors.New("invalid password")
	}

	return nil
}

// GetUserIDByUsername retrieves the user ID by username from the MongoDB database
func (m *MongoDB) GetUserIDByUsername(userName string, userID *int) error {
	log.Println("Start Function GetUserIDByUsername")

	filter := bson.M{"user_name": userName}
	var user domain.User
	err := m.database.Collection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}

	*userID, err = strconv.Atoi(user.UserID)
	if err != nil {
		return err
	}

	return nil
}
