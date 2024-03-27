package taskhandler

import (
	taskmanager "Simple_Task_Manager/task_manager"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

type App struct {
	DB          *sql.DB
	TaskManager *taskmanager.App
}

func (app *App) HandleTasks(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	log.Printf("Decoded request body: %+v\n", requestBody)

	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if len(requestBody) == 1 && requestBody["task_id"] != nil {
			taskID := int(requestBody["task_id"].(float64))
			app.TaskManager.GetTaskByID(w, taskID)
		} else {
			app.TaskManager.GetTasks(w)
		}
	case http.MethodPost:
		userName := requestBody["user_name"].(string)
		taskName := requestBody["task_name"].(string)
		dueDate := requestBody["due_date"].(string)
		app.TaskManager.CreateTask(w, userName, taskName, dueDate)
	case http.MethodPatch:
		taskID := int(requestBody["task_id"].(float64))
		userID := int(requestBody["user_id"].(float64))
		app.TaskManager.UpdateTask(w, taskID, userID)
	case http.MethodDelete:
		taskID := int(requestBody["task_id"].(float64))
		userID := int(requestBody["user_id"].(float64))
		app.TaskManager.DeleteTask(w, taskID, userID)
	default:
		log.Printf("Method %s not allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
