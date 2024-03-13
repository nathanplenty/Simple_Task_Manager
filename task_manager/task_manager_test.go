package taskmanager

import (
	"testing"
)

func TestInitDatabase(t *testing.T) {

}

func TestHandleTasks(t *testing.T) {
	// Testen der handleTasks-Funktion
}

func TestGetTasks(t *testing.T) {
	// Testen der getTasks-Funktion
}

func TestCreateTasks(t *testing.T) {
	// Testen der createTasks-Funktion
}

func TestUpdateTasks(t *testing.T) {
	// Testen der updateTasks-Funktion
}

func TestDeleteTasks(t *testing.T) {
	// Testen der deleteTasks-Funktion
}

func TestGetTasksByID(t *testing.T) {
	// Testen der getTasksByID-Funktion
}

func TestSimple(t *testing.T) {
	result := simple()
	if result != true {
		t.Errorf("Boolean = %v; want true", result)
	}
}
