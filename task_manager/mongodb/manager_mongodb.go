package taskManagerMongoDB

import (
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type TaskManager interface {
	GetTaskByID(w http.ResponseWriter, taskID, userID string)
	UpdateTask(w http.ResponseWriter, taskID, userID string)
	DeleteTask(w http.ResponseWriter, taskID, userID string)
	GetTasks(w http.ResponseWriter)
	CreateTask(w http.ResponseWriter, userName, taskName, dueDate string)
}

type App struct {
	DB    *mongo.Database
	Users *mongo.Collection
	Tasks *mongo.Collection
}

type Task struct {
	TaskID    primitive.ObjectID `json:"task_id" bson:"_id"`
	TaskName  string             `json:"task_name" bson:"task_name"`
	DueDate   string             `json:"due_date" bson:"due_date"`
	Completed bool               `json:"completed" bson:"completed"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type User struct {
	UserID   primitive.ObjectID `json:"user_id" bson:"_id"`
	UserName string             `json:"user_name" bson:"user_name"`
}

func (app *App) GetTaskByID(w http.ResponseWriter, taskID, userID string) {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		log.Println("Invalid task ID:", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskCollection := app.Tasks
	filter := bson.M{"_id": objectID, "user_id": userID}

	var task Task
	err = taskCollection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("Task not found")
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving task: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Printf("Error encoding task to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	log.Println("Task retrieved successfully")
}

func (app *App) UpdateTask(w http.ResponseWriter, taskID, userID string) {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		log.Println("Invalid task ID:", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskCollection := app.Tasks
	filter := bson.M{"_id": objectID, "user_id": userID}

	update := bson.M{
		"$set": bson.M{
			"completed": true,
		},
	}

	_, err = taskCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Task updated successfully")
	_, _ = w.Write([]byte("Task updated successfully"))
}

func (app *App) DeleteTask(w http.ResponseWriter, taskID, userID string) {
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		log.Println("Invalid task ID:", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	taskCollection := app.Tasks
	filter := bson.M{"_id": objectID, "user_id": userID}

	update := bson.M{
		"$set": bson.M{
			"task_name": "X",
			"due_date":  "0001-01-01T00:00:00Z",
			"completed": false,
		},
	}

	_, err = taskCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error anonymizing task: %v", err)
		http.Error(w, "Error anonymizing task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Task content anonymized successfully")
	_, _ = w.Write([]byte("Task anonymized successfully"))
}

func (app *App) GetTasks(w http.ResponseWriter) {
	taskCollection := app.Tasks
	cursor, err := taskCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error querying tasks from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var tasks []Task
	for cursor.Next(context.Background()) {
		var task Task
		err := cursor.Decode(&task)
		if err != nil {
			log.Printf("Error decoding task: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	err = cursor.Err()
	if err != nil {
		log.Printf("Error iterating over task cursor: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Error encoding tasks to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(responseBody)
	if err != nil {
		log.Printf("Error writing response body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Tasks gathered successfully")
}

func (app *App) CreateTask(w http.ResponseWriter, userName, taskName, dueDate string) {
	user := User{
		UserID:   primitive.NewObjectID(),
		UserName: userName,
	}

	_, err := app.Users.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error creating new user: %v", err)
		http.Error(w, "Error creating new user", http.StatusInternalServerError)
		return
	}

	task := Task{
		TaskID:    primitive.NewObjectID(),
		TaskName:  taskName,
		DueDate:   dueDate,
		Completed: false,
		UserID:    user.UserID,
	}

	_, err = app.Tasks.InsertOne(context.Background(), task)
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Task created successfully")
	_, _ = w.Write([]byte("Task created successfully"))
}
