package routerHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskHandler interface {
	HandleTasks(w http.ResponseWriter, r *http.Request)
}

type TaskManager interface {
	GetTaskByID(w http.ResponseWriter, taskID int)
	UpdateTask(w http.ResponseWriter, taskID, userID int)
	DeleteTask(w http.ResponseWriter, taskID, userID int)
	GetTasks(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, userName, taskName, dueDate string)
}

type App struct {
	TaskManager TaskManager
}

func (app *App) HandleTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	taskIDStr := query.Get("task_id")
	if taskIDStr != "" {
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			log.Println("Invalid task_id parameter:", err)
			http.Error(w, "Invalid task_id parameter", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			app.TaskManager.GetTaskByID(w, taskID)
		case http.MethodPatch:
			userIDStr := query.Get("user_id")
			if userIDStr == "" {
				log.Println("Missing user_id parameter")
				http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
				return
			}
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Println("Invalid user_id parameter:", err)
				http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
				return
			}
			app.TaskManager.UpdateTask(w, taskID, userID)
		case http.MethodDelete:
			userIDStr := query.Get("user_id")
			if userIDStr == "" {
				log.Println("Missing user_id parameter")
				http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
				return
			}
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Println("Invalid user_id parameter:", err)
				http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
				return
			}
			app.TaskManager.DeleteTask(w, taskID, userID)
		default:
			log.Printf("Method %s not allowed", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			app.TaskManager.GetTasks(w, r)
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
