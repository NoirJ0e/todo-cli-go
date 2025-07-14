package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Expected JSON content, got %s", contentType)
	}
}

func TestCreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_create_tasks.json"
	err := saveTasksToFile(&[]taskStruct{}, testFileName)
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}
	t.Cleanup(func() { os.Remove(testFileName) })

	router := setupRouter()

	requestBody := `{"content": "New task from HTTP"}`
	req, _ := http.NewRequest("POST", "/tasks/create", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdTask taskStruct
	err = json.Unmarshal(w.Body.Bytes(), &createdTask)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if createdTask.Content != "New task from HTTP" {
		t.Errorf("Expected content 'New task from HTTP', got %s", createdTask.Content)
	}

	if createdTask.ID == "" {
		t.Error("Expected task to have an ID")
	}

	if createdTask.IsComplete {
		t.Error("Expected new task to be incomplete")
	}

	tasks, err := loadTasksFromFile(testFileName)
	if err != nil {
		t.Errorf("Failed to load tasks from file: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task in file, got %d", len(tasks))
	}
}
