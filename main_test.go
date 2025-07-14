package main

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateValidTask(t *testing.T) {
	taskContent := "Test"
	testTask := createTask(taskContent)
	if testTask.Content != taskContent {
		t.Errorf("Content was not set correctly")
	}

	if testTask.IsComplete != false {
		t.Errorf("complete mark was not set correctly")
	}

	if time.Since(testTask.CreateDate) >= time.Second {
		t.Errorf("initial create time was not set correctly")
	}

	if !testTask.CompleteDate.IsZero() {
		t.Errorf("initial complete time was not set correctly")
	}

	_, err := uuid.Parse(testTask.ID)
	if err != nil {
		t.Errorf("id is not a valid uuid: %v", err)
	}
}

func TestAddTask(t *testing.T) {
	// define a slice to include tasks first
	var tasks []taskStruct

	// create a content first
	testTaskContent := "This a task content"

	// create a task
	testTask := createTask(testTaskContent)
	addTask(&tasks, testTask)

	// verify if the tasks slice is not null now
	if len(tasks) == 0 {
		t.Errorf("Failed to add task")
	}

	// test if the content in the newlly add tasks is same as pre-defined
	if tasks[0].Content != testTaskContent {
		t.Errorf("Failed to create task: incorrect task content")
	}
}

func TestRemoveTask(t *testing.T) {
	// Setup: Create a slice of tasks
	tasks := []taskStruct{
		createTask("Task 1"), // This will be removed
		createTask("Task 2"),
		createTask("Task 3"),
	}
	initialCount := len(tasks)
	taskToRemoveID := tasks[0].ID

	// Execute the function
	err := removeTask(&tasks, taskToRemoveID)
	if err != nil {
		t.Fatalf("removeTask() returned an unexpected error: %v", err)
	}

	// Verify the results
	if len(tasks) != initialCount-1 {
		t.Errorf("Expected task count to be %d, but got %d", initialCount-1, len(tasks))
	}

	for _, task := range tasks {
		if task.ID == taskToRemoveID {
			t.Errorf("Task with ID %s was not removed", taskToRemoveID)
		}
	}

	// Test case for trying to remove a non-existent task
	err = removeTask(&tasks, "non-existent-id")
	if err == nil {
		t.Errorf("Expected an error when trying to remove a non-existent task, but got nil")
	}
}

func TestUpdateTask(t *testing.T) {
	tasks := []taskStruct{
		createTask("Original Task 1"),
		createTask("Original Task 2"),
	}
	t.Run("update existing task succeeds", func(t *testing.T) {
		newContent := "Updateded Task Content"
		err := updateTask(&tasks, tasks[0].ID, newContent)
		if err != nil {
			t.Errorf("updateTask should not return error for existing task: %v", err)
		}
		if tasks[0].Content != newContent {
			t.Errorf("Expected content %s, got %s", newContent, tasks[0].Content)
		}
	})
	t.Run("update non-existent task returns error", func(t *testing.T) {
		err := updateTasks(&tassks, "fake-uuid", "something")
		if err == nil {
			t.Errorf("updateTask should return error for non-existent task")
		}
	})
}

func TestCompleteTask(t *testing.T) {
	tasks := []taskStruct{
		createTask("Original Task 1"),
		createTask("Original Task 2"),
	}

	t.Run("complete existing task succeeds", func(t *testing.T) {
		err := completeTask(&tasks, tasks[0].ID)
		if err != nil {
			t.Errorf("completeTask should not return for existing task: %v", err)
		}

		if !tasks[0].IsComplete {
			t.Errorf("Task should be marked as complete")
		}
		if tasks[0].CompleteDate.IsZero() {
			t.Errorf("CompleteDate should be test")
		}
	})
	t.Run("complete non-existing task returns error", func(t *testing.T) {
		err := completeTask(&tasks, tasks[0].ID)
		if err == nil {
			t.Errorf("completeTask should return error for non-existent task")
		}
	})
}

func TestSaveDataToFile(t *testing.T) {
	// verify if the os has the test file already, if so delete it so this test
	// can run multiple times
	testFileName := "Test_File.json"

	// create test tasks slice for saving
	tasks := []taskStruct{
		createTask("Test1"),
		createTask("Test2"),
	}
	err := saveTasksToFile(&tasks, testFileName)
	if err != nil {
		t.Errorf("saveTasksToFile() error: %v", err)
	}

	t.Cleanup(func() { os.Remove(testFileName) })

	if _, err := os.Stat(testFileName); os.IsNotExist(err) {
		t.Errorf("saveTasksToFile() did not create file: %s", testFileName)
	}
}

func TestLoadTasksFromFile(t *testing.T) {
	testFileName := "Test_File.json"

	// create test tasks slice for saving
	originalTasks := []taskStruct{
		createTask("Test1"),
		createTask("Test2"),
	}
	err := saveTasksToFile(&originalTasks, testFileName)
	if err != nil {
		t.Errorf("saveTasksToFile() error: %v", err)
	}

	t.Cleanup(func() { os.Remove(testFileName) })

	loadedTasks, err := loadTasksFromFile(testFileName)
	if err != nil {
		t.Errorf("loadTasksFromFile() did not load file %s: %v", testFileName, err)
	}

	if len(loadedTasks) != len(originalTasks) {
		t.Errorf("loadTasksFromFile() did not load file %s correctly: expecting load %d ; have %d", testFileName, len(originalTasks), len(loadedTasks))
	}
	for i, task := range loadedTasks {
		if task.Content != originalTasks[i].Content {
			t.Errorf("loadTasksFromFile() loaded task compromised. Expecting %s, have %s", originalTasks[i].Content, task.Content)
		}
	}
}

func TestRun(t *testing.T) {
	// Helper function to set up test data
	setupTestFile := func(t *testing.T, fileName string) []taskStruct {
		t.Helper() // This marks the function as a test helper
		testTasks := []taskStruct{
			createTask("Test1"),
			createTask("Test2"),
		}
		err := saveTasksToFile(&testTasks, fileName)
		if err != nil {
			t.Fatalf("Failed to setup test file: %v", err)
		}
		t.Cleanup(func() { os.Remove(fileName) })
		return testTasks
	}

	t.Run("With no args", func(t *testing.T) {
		args := []string{"todo"}
		err := run(args)
		if err != nil {
			t.Errorf("run() with no args error: %v", err)
		}
	})

	t.Run("With invalid args", func(t *testing.T) {
		args := []string{"todo", "foo"}
		err := run(args)
		if err == nil {
			t.Errorf("run() with invalid args should return an error, but got nil")
		}
	})

	// t.Run("With `add` flag", func(t *testing.T) {
	// 	testFileName := "Test_File.json"
	// 	testTasks := setupTestFile(t, testFileName)
	//
	// 	// Now you can use testTasks and testFileName in your test
	// 	// Add your test logic here
	// 	args := []string{"todo", "add", "newly added task"}
	// 	err := run(args)
	// 	if err != nil {
	// 		t.Errorf("run() with invalid args error: %v", err)
	// 	}
	// 	if len(testTasks) != 3 {
	// 		t.Errorf("run() with `add` flag error: should have 3 tasks in total, but now has %d", len(testTasks))
	// 	} else if testTasks[2].Content != "newly added task" {
	// 		t.Errorf("run() with `add` flag error: should create new task with \"newly added task\", but now has %s", testTasks[2].Content)
	// 	}
	// })

	t.Run("add commands adds a task to the default file", func(t *testing.T) {
		os.Remove(tasksFileName)
		t.Cleanup(func() { os.Remove(tasksFileName) })

		taskContent := "buy milk from default file"
		args := []string{"todo", "add", taskContent}
		if err := run(args); err != nil {
			t.Fatalf("run() returned an error: %v", err)
		}
		tasks, err := loadTasksFromFile(tasksFileName)
		if err != nil {
			t.Fatalf("Could not load tasks from default file: %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("Expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Content != taskContent {
			t.Fatalf("Expected task content %s, got %s", taskContent, tasks[0].Content)
		}
	})
	t.Run("add commands adds a task to a specific file", func(t *testing.T) {
		testFileName := "custom_add_test.json"
		os.Remove(testFileName)
		t.Cleanup(func() { os.Remove(testFileName) })

		taskContent := "buy milk from specific file"

		args := []string{"todo", "add", testFileName, taskContent}
		if err := run(args); err != nil {
			t.Fatalf("run() returned an error: %v", err)
		}
		tasks, err := loadTasksFromFile(testFileName)
		if err != nil {
			t.Fatalf("Could not load tasks from %s file: %v", testFileName, err)
		}
		if len(tasks) != 1 {
			t.Fatalf("Expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Content != taskContent {
			t.Fatalf("Expected task content %s, got %s", taskContent, tasks[0].Content)
		}
	})

	t.Run("With `remove` flag", func(t *testing.T) {
		testFileName := "Test_File_Remove.json"
		testTasks := setupTestFile(t, testFileName)
		taskToRemove := testTasks[0]

		args := []string{"todo", "remove", testFileName, taskToRemove.ID}
		if err := run(args); err != nil {
			t.Fatalf("run() with remove command returned an error: %v", err)
		}
		remainingTasks, err := loadTasksFromFile(testFileName)
		if err != nil {
			t.Fatalf("run() with remove command returned an error: %v", err)
		}
		if len(remainingTasks) != 1 {
			t.Fatalf("run() with remove command returned an error: expected 1 task left, got %d", len(remainingTasks))
		}
		if remainingTasks[0].ID == taskToRemove.ID {
			t.Fatalf("run() with remove command returned an error: task with ID %s should be removed, but it still exist", taskToRemove.ID)
		}
		if remainingTasks[0].ID != testTasks[1].ID {
			t.Errorf("run() with remove command returned an error: wrong task been removed from the file")
		}
	})

	t.Run("complete command marks a task as complete", func(t *testing.T) {
		testFileName := "tasks_for_complete_test.json"
		initialTasks := setupTestFile(t, testFileName)
		taskToComplete := initialTasks[0]

		args := []string{"todo", "complete", testFileName, taskToComplete.ID}
		if err := run(args); err != nil {
			t.Fatalf("run() with complete command returned an error: %v", err)
		}
		updatedTasks, err := loadTasksFromFile(testFileName)
		if err != nil {
			t.Fatalf("failed to load tasks after complete: %v", err)
		}
		var completedTask taskStruct
		for _, task := range updatedTasks {
			if task.ID == taskToComplete.ID {
				completedTask = task
				break
			}
		}
		if completedTask.ID == "" {
			t.Fatalf("could not find the task that supposed to be completed")
		}
		if !completedTask.IsComplete {
			t.Fatalf("taskw ith ID %s should have been marked as complete, but it wasnt", completedTask.ID)
		}
		if completedTask.CompleteDate.IsZero() {
			t.Errorf("the CompleteDate for task %s was not set.", completedTask.ID)
		}
	})
	t.Run("update command should update the task with new content", func(t *testing.T) {
		testFileName := "tasks_for_list_test.json"
		initialTasks := setupTestFile(t, testFileName)
		initialTaskID := initialTasks[0].ID
		newTaskContent := "This is a new task content"

		args := []string{"todo", "update", testFileName, initialTaskID, newTaskContent}
		if err := run(args); err != nil {
			t.Fatalf("run() with list command returned an error: %v", err)
		}
		updatedTasks, err := loadTasksFromFile(testFileName)
		if err != nil {
			t.Errorf("loadTasksFromFile() error: %v", err)
		}
		if updatedTasks[0].Content != newTaskContent {
			t.Errorf("task with ID %s should contain content as %s, but have %s ", updatedTasks[0].ID, newTaskContent, updatedTasks[0].Content)
		}
	})
}
