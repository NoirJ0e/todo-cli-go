package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/tasks", getTasksHandler)
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
}
