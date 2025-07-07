package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
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

func removeTask(tasks *[]taskStruct, id string) error {
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
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	validFlags := []string{"add", "remove", "complete"}

	if len(args) < 2 {
		return nil
	}

	command := args[1]

	if !slices.Contains(validFlags, command) {
		return fmt.Errorf("unknown command %q", command)
	}

	switch command {
	case "add":
		if len(args) < 3 {
			return fmt.Errorf("missing task content")
		}
		fileName := tasksFileName
		taskContent := args[2]
		// check if a custom file name is provided
		// todo add <fileName.json> <content>
		if len(args) > 3 && len(args[2]) > 5 && args[2][len(args[2])-5:] == ".json" {
			fileName = args[2]
			taskContent = args[3]
		}
		tasks, err := loadTasksFromFile(fileName)
		if err != nil {
			return err
		}
		addTask(&tasks, createTask(taskContent))
		return saveTasksToFile(&tasks, fileName)
	case "remove":
		if len(args) < 3 {
			return fmt.Errorf("missing task ID")
		}
		fileName := tasksFileName
		taskID := args[2]
		// check if a custom file name is provided
		// todo add <fileName.json> <content>
		if len(args) > 3 && len(args[2]) > 5 && args[2][len(args[2])-5:] == ".json" {
			fileName = args[2]
			taskID = args[3]
		}
		tasks, err := loadTasksFromFile(fileName)
		if err != nil {
			return err
		}
		removeTask(&tasks, taskID)
		return saveTasksToFile(&tasks, fileName)
	}

	return nil
}
