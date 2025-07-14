package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Content string `json:"content" binding: "required"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/tasks", getTasksHandler)
	router.POST("/tasks", createTaskHandler)
	return router
}

func getTasksHandler(c *gin.Context) {
	tasks, err := loadTasksFromFile(tasksFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, tasks)
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
