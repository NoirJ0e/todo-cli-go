package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// create a struct to define a "task"
// that should contain task content, task creation date, task id and if task is
// a done at least
type taskStruct struct {
	id           string // will use uuid
	content      string
	createDate   time.Time
	completeDate time.Time
	isComplete   bool
}

func createTask(content string) taskStruct {
	task := taskStruct{
		id:         uuid.NewString(),
		content:    content,
		createDate: time.Now(),
		isComplete: false,
	}
	return task
}

func addTask(tasks *[]taskStruct, task taskStruct) {
	*tasks = append(*tasks, task)
}

func main() {
	fmt.Println("TODO App")
}
