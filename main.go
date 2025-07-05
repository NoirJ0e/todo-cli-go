package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

// create a struct to define a "task"
// that should contain task content, task creation date, task id and if task is
// a done at least
type taskStruct struct {
	ID           string // will use uuid
	Content      string
	CreateDate   time.Time
	CompleteDate time.Time
	IsComplete   bool
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

func main() {
	fmt.Println("TODO App")
}
