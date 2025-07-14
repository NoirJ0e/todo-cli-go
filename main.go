package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	validFlags := []string{"add", "remove", "complete", "update"}

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
	case "complete":
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
		completeTask(&tasks, taskID)
		return saveTasksToFile(&tasks, fileName)
	case "update":
		if len(args) < 3 {
			return fmt.Errorf("missing task ID")
		}
		fileName := tasksFileName
		taskID := args[2]
		var newTaskContent string
		// check if a custom file name is provided
		// todo add <fileName.json> <content>
		if len(args) > 3 && len(args[2]) > 5 && args[2][len(args[2])-5:] == ".json" {
			fileName = args[2]
			taskID = args[3]
			newTaskContent = args[4]
		}
		tasks, err := loadTasksFromFile(fileName)
		if err != nil {
			return err
		}
		updateTask(&tasks, taskID, newTaskContent)
		return saveTasksToFile(&tasks, fileName)
	}

	return nil
}
