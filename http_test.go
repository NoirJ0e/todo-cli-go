package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetTasksHTTP(t *testing.T) {
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

func TestCreateTaskHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_create_tasks.json"
	err := saveTasksToFile(&[]taskStruct{}, testFileName)
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}
	t.Cleanup(func() { os.Remove(testFileName) })

	originalFileName := tasksFileName
	tasksFileName = testFileName
	t.Cleanup(func() { tasksFileName = originalFileName })

	router := setupRouter()

	requestBody := `{"content": "New task from HTTP"}`
	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(requestBody))
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

func TestUpdateTaskHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_update_tasks.json"
	testTasks := []taskStruct{
		createTask("Task 1"),
		createTask("Task 2"),
	}
	err := saveTasksToFile(&testTasks, testFileName)
	if err != nil {
		t.Fatalf("Failed to setup test file: %v", err)
	}
	t.Cleanup(func() { os.Remove(testFileName) })

	originalFileName := tasksFileName
	tasksFileName = testFileName
	t.Cleanup(func() { tasksFileName = originalFileName })

	router := setupRouter()

	taskToUpdate := testTasks[0]
	requestBody := `{"content": "Updated Task Content"}`
	url := fmt.Sprintf("/tasks/%s", taskToUpdate.ID)
	req, _ := http.NewRequest("PUT", url, strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var updatedTask taskStruct
	err = json.Unmarshal(w.Body.Bytes(), &updatedTask)
	if err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if updatedTask.ID != taskToUpdate.ID {
		t.Errorf("Expected ID %s, got %s", taskToUpdate.ID, updatedTask.ID)
	}

	tasks, err := loadTasksFromFile(testFileName)
	if err != nil {
		t.Errorf("Failed to load tasks from file: %v", err)
	}

	var foundTask taskStruct
	for _, task := range tasks {
		if task.ID == taskToUpdate.ID {
			foundTask = task
			break
		}
	}

	if foundTask.Content != "Updated Task Content" {
		t.Errorf("Task was not updated in file. Expected 'Updated Task Content', but got %s", foundTask.Content)
	}
}
