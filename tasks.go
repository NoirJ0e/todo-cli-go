package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

var tasksFileName = "tasks.json"

// const tasksFileName = "tasks.json"

type taskStruct struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	CreateDate   time.Time `json:"createDate"`
	CompleteDate time.Time `json:"completeDate"`
	IsComplete   bool      `json:"isComplete"`
}

func createTask(content string) taskStruct {
	task := taskStruct{
		ID:         uuid.NewString(),
		Content:    content,
		CreateDate: time.Now(),
		IsComplete: false,
	}
	return task
}

func addTask(tasks *[]taskStruct, task taskStruct) {
	*tasks = append(*tasks, task)
}

func removeTask(tasks *[]taskStruct, id string) error {
	taskIndex := -1
	for taskID, task := range *tasks {
		if task.ID == id {
			taskIndex = taskID
			break
		}
	}
	if taskIndex == -1 {
		return fmt.Errorf("task id not found")
	}
	*tasks = append((*tasks)[:taskIndex], (*tasks)[taskIndex+1:]...)

	return nil
}

func saveTasksToFile(tasks *[]taskStruct, fileName string) error {
	jsonData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0o644)
	if err != nil {
		return err
	}

	return err
}

func loadTasksFromFile(fileName string) ([]taskStruct, error) {
	var tasks []taskStruct
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []taskStruct{}, nil
		}
		return nil, err
	}
	err = json.Unmarshal(fileContent, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func completeTask(tasks *[]taskStruct, taskID string) error {
	for i := range *tasks {
		if (*tasks)[i].ID == taskID {
			(*tasks)[i].IsComplete = true
			(*tasks)[i].CompleteDate = time.Now()
			return nil
		}
	}
	return fmt.Errorf("task with ID %s not found", taskID)
}

func updateTask(tasks *[]taskStruct, taskID string, newTaskContent string) error {
	for i := range *tasks {
		if (*tasks)[i].ID == taskID {
			(*tasks)[i].Content = newTaskContent
			return nil
		}
	}
	return fmt.Errorf("task with id %s not found", taskID)
}
