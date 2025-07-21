package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateTaskRequest struct {
	Content string `json:"content" binding:"required"`
}

type CompleteTaskRequest struct {
	IsComplete bool `json:"isComplete" binding:"required"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/tasks", getTasksHandler)
	router.POST("/tasks", createTaskHandler)
	router.PUT("/tasks/:id", updateTaskHandler)
	router.PATCH("/tasks/:id", completeTaskHandler)
	router.DELETE("/tasks/:id", deleteTaskHandler)
	return router
}

func getTasksHandler(c *gin.Context) {
	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get query parameter "isComplete" from the URL
	// c.Query() returns the value as a string, or empty string if not found
	isCompleteParam := c.Query("isComplete")

	// If no filter is specified, return all tasks (original behavior)
	if isCompleteParam == "" {
		c.JSON(http.StatusOK, tasks)
		return
	}

	// Convert string parameter to boolean
	// "true" -> true, "false" -> false, anything else -> error
	var filterComplete bool
	switch isCompleteParam {
	case "true":
		filterComplete = true
	case "false":
		filterComplete = false
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "isComplete must be 'true' or 'false'"})
		return

	}

	// Filter the tasks based on the isComplete status
	var filteredTasks []taskStruct
	for _, task := range tasks {
		if task.IsComplete == filterComplete {
			filteredTasks = append(filteredTasks, task)
		}
	}

	// Return the filtered results
	c.JSON(http.StatusOK, filteredTasks)
}

func createTaskHandler(c *gin.Context) {
	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// parse the JSON input
	var request CreateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask := createTask(request.Content)
	addTask(&tasks, newTask)

	if err := saveTasksToFile(&tasks, tasksFileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newTask)
}

func updateTaskHandler(c *gin.Context) {
	taskID := c.Param("id")
	var request UpdateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	err = updateTask(&tasks, taskID, request.Content)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := saveTasksToFile(&tasks, tasksFileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedTask taskStruct
	for _, task := range tasks {
		if task.ID == taskID {
			updatedTask = task
			break
		}
	}
	c.JSON(http.StatusOK, updatedTask)
}

func completeTaskHandler(c *gin.Context) {
	taskID := c.Param("id")
	var request CompleteTaskRequest
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = completeTask(&tasks, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := saveTasksToFile(&tasks, tasksFileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedTask taskStruct
	for _, task := range tasks {
		if task.ID == taskID {
			updatedTask = task
			break
		}
	}
	c.JSON(http.StatusOK, updatedTask)
}

func deleteTaskHandler(c *gin.Context) {
	taskID := c.Param("id")

	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = removeTask(&tasks, taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := saveTasksToFile(&tasks, tasksFileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func filterTaskHandler(c *gin.Context) {
}
