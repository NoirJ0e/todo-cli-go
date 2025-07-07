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

		// Add your test logic here
		_ = testTasks // Remove this line when you implement the test
	})

	t.Run("With `complete` flag", func(t *testing.T) {
		testFileName := "Test_File_Complete.json"
		testTasks := setupTestFile(t, testFileName)

		// Add your test logic here
		_ = testTasks // Remove this line when you implement the test
	})
}
