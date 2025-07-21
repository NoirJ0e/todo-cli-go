package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

func TestUpdateNonExistentTaskHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_update_empty.json"
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

	requestBody := `{"content": "Updated Content"}`
	req, _ := http.NewRequest("PUT", "/tasks/fake-uuid-123", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, but got %d", w.Code)
	}
}

func TestCompleteTaskHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_update_empty.json"
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
	requestBody := `{"isComplete": true}`
	url := fmt.Sprintf("/tasks/%s", taskToUpdate.ID)
	req, _ := http.NewRequest("PATCH", url, strings.NewReader(requestBody))
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

	if foundTask.IsComplete != true {
		t.Errorf("Task was not updated in file. Expected mark as complete, but got %t", foundTask.IsComplete)
	}
	if foundTask.CompleteDate.IsZero() {
		t.Error("Expected CompleteDate to be set, but it is zero")
	}
	if time.Since(foundTask.CompleteDate) > time.Second {
		t.Error("CompleteDate should be recent")
	}
}

func TestDeleteTaskHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_delete.json"
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

	taskToDelete := testTasks[0]
	url := fmt.Sprintf("/tasks/%s", taskToDelete.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	tasks, err := loadTasksFromFile(testFileName)
	if err != nil {
		t.Errorf("Failed to load tasks from file: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected to have 1 taks left, but have %d", len(tasks))
	}

	for _, task := range tasks {
		if task.ID == taskToDelete.ID {
			t.Errorf("Expected task with ID %s be deleted, but it still exist", task.ID)
			break
		}
	}
}

func TestFilterTaskBasedOnStatusHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testFileName := "test_http_filter_status.json"
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

	// Mark task 1 as complete and save back to file
	if err := completeTask(&testTasks, testTasks[0].ID); err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}
	saveTasksToFile(&testTasks, tasksFileName)

	router := setupRouter()

	t.Run("filter completed tasks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks?isComplete=true", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, but got %d", w.Code)
		}
		var results []taskStruct
		err := json.Unmarshal(w.Body.Bytes(), &results)
		if err != nil {
			t.Errorf("Error handling return result: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("Expected results has length of 1, but got %d", len(results))
		}
		if results[0].IsComplete != true {
			t.Error("Expected results to be complete, but it's not")
		}
	})

	t.Run("filter incomplete tasks", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks?isComplete=false", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, but got %d", w.Code)
		}
		var results []taskStruct
		err := json.Unmarshal(w.Body.Bytes(), &results)
		if err != nil {
			t.Errorf("Error handling return result: %v", err)
		}
		if len(results) != 1 {
			t.Errorf("Expected results has length of 1, but got %d", len(results))
		}
		if results[0].IsComplete != false {
			t.Error("Expected results to be incomplete, but it's not")
		}
	})
}
