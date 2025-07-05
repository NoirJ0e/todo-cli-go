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

	//
	t.Cleanup(func() { os.Remove(testFileName) })

	if _, err := os.Stat(testFileName); os.IsNotExist(err) {
		t.Errorf("saveTasksToFile() did not create file: %s", testFileName)
	}
}
