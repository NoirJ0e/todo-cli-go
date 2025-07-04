package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateValidTask(t *testing.T) {
	taskContent := "Test"
	testTask := createTask(taskContent)
	if testTask.content != taskContent {
		t.Errorf("Content was not set correctly")
	}

	if testTask.isComplete != false {
		t.Errorf("complete mark was not set correctly")
	}

	if time.Since(testTask.createDate) >= time.Second {
		t.Errorf("initial create time was not set correctly")
	}

	if !testTask.completeDate.IsZero() {
		t.Errorf("initial complete time was not set correctly")
	}

	_, err := uuid.Parse(testTask.id)
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
	if tasks[0].content != testTaskContent {
		t.Errorf("Failed to create task: incorrect task content")
	}
}

func TestRecordExist(t *testing.T) {
}
