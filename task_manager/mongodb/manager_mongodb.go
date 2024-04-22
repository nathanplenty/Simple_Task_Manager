package taskManagerMongodb

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskManager interface {
	GetTaskByID(w http.ResponseWriter, taskID string)
	UpdateTask(w http.ResponseWriter, taskID, userID string)
	DeleteTask(w http.ResponseWriter, taskID, userID string)
	GetTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, userName, taskName, dueDate string)
}

type App struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

type Task struct {
	TaskID    int       `json:"task_id" bson:"task_id"`
	TaskName  string    `json:"task_name" bson:"task_name"`
	DueDate   time.Time `json:"due_date" bson:"due_date"`
	Completed bool      `json:"completed" bson:"completed"`
}

type User struct {
	UserID   int    `json:"user_id" bson:"user_id"`
	UserName string `json:"user_name" bson:"user_name"`
}

func (app *App) GetTaskByID(w http.ResponseWriter, taskID int) {
	var task Task
	filter := bson.M{"task_id": taskID}
	err := app.Collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
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

func (app *App) UpdateTask(w http.ResponseWriter, taskID int) {
	filter := bson.M{"task_id": taskID}
	update := bson.M{"$set": bson.M{"completed": true}}
	result, err := app.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Task updated successfully")
	_, _ = w.Write([]byte("Task updated successfully"))
}

func (app *App) DeleteTask(w http.ResponseWriter, taskID int) {
	filter := bson.M{"task_id": taskID}
	update := bson.M{"$set": bson.M{"task_name": "X", "due_date": time.Time{}, "completed": false}}
	result, err := app.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error anonymizing task: %v", err)
		http.Error(w, "Error anonymizing task", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		log.Println("Task not found")
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Task content anonymized successfully")
	_, _ = w.Write([]byte("Task anonymized successfully"))
}

func (app *App) GetTasks(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("task_id")
	if taskIDStr != "" {
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			log.Println("Invalid task_id parameter:", err)
			http.Error(w, "Invalid task_id parameter", http.StatusBadRequest)
			return
		}
		app.GetTaskByID(w, taskID)
		return
	}

	filter := bson.M{}
	cursor, err := app.Collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error querying tasks from database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if cerr := cursor.Close(context.Background()); cerr != nil {
			log.Printf("Error closing cursor: %v", cerr)
		}
	}()

	var tasks []Task
	for cursor.Next(context.Background()) {
		var task Task
		err = cursor.Decode(&task)
		if err != nil {
			log.Printf("Error decoding task document: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err = cursor.Err(); err != nil {
		log.Printf("Error iterating over task documents: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Error encoding tasks to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		log.Printf("Error writing response body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Tasks gathered successfully")
}

func (app *App) CreateTask(w http.ResponseWriter, userName, taskName, dueDate string) {
	dueDateTime, err := time.Parse(time.RFC3339, dueDate)
	if err != nil {
		log.Printf("Error parsing due date: %v", err)
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	_, err = app.Collection.InsertOne(context.Background(), bson.M{
		"user_name": userName,
		"task_name": taskName,
		"due_date":  dueDateTime,
		"completed": false,
	})
	if err != nil {
		log.Printf("Error inserting task: %v", err)
		http.Error(w, "Error inserting task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Task created successfully")
	_, _ = w.Write([]byte("Task created successfully"))
}
