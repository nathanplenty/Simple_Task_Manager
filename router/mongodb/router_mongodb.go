package routerMongoDB

import (
	"encoding/json"
	"log"
	"net/http"
)

type TaskHandler interface {
	HandleTasks(w http.ResponseWriter, r *http.Request)
}

type TaskManager interface {
	GetTaskByID(w http.ResponseWriter, taskID, userID string)
	UpdateTask(w http.ResponseWriter, taskID, userID string)
	DeleteTask(w http.ResponseWriter, taskID, userID string)
	GetTasks(w http.ResponseWriter)
	CreateTask(w http.ResponseWriter, userName, taskName, dueDate string)
}

type App struct {
	TaskManager TaskManager
}

func (app *App) HandleTasks(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	taskIDStr := vars.Get("task_id")
	userIDStr := vars.Get("user_id")

	if taskIDStr != "" {
		switch r.Method {
		case http.MethodGet:
			app.TaskManager.GetTaskByID(w, taskIDStr, userIDStr)
		case http.MethodPatch:
			userIDStr := vars.Get("user_id")
			if userIDStr == "" {
				log.Println("Missing user_id parameter")
				http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
				return
			}
			app.TaskManager.UpdateTask(w, taskIDStr, userIDStr)
		case http.MethodDelete:
			userIDStr := vars.Get("user_id")
			if userIDStr == "" {
				log.Println("Missing user_id parameter")
				http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
				return
			}
			app.TaskManager.DeleteTask(w, taskIDStr, userIDStr)
		default:
			log.Printf("Method %s not allowed", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			app.TaskManager.GetTasks(w)
		case http.MethodPost:
			var requestBody struct {
				UserName string `json:"user_name"`
				TaskName string `json:"task_name"`
				DueDate  string `json:"due_date"`
			}
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				log.Println("Error decoding request body:", err)
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			if requestBody.UserName == "" || requestBody.TaskName == "" || requestBody.DueDate == "" {
				log.Println("Missing required parameters")
				http.Error(w, "Missing required parameters", http.StatusBadRequest)
				return
			}
			app.TaskManager.CreateTask(w, requestBody.UserName, requestBody.TaskName, requestBody.DueDate)
		default:
			log.Printf("Method %s not allowed", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
