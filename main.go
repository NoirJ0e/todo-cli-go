package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

const tasksFileName = "tasks.json"

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

func main() {
	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- TODOs ---")
	if len(tasks) == 0 {
		fmt.Println("No tasks found. Add one!")
	} else {
		for i, task := range tasks {
			status := " "
			if task.IsComplete {
				status = "✔"
			}
			fmt.Printf("%d. [%s] %s\n", i+1, status, task.Content)
		}
	}
}
