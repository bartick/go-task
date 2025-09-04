package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/bartick/go-task/app/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func HandlerGetTasks(c *gin.Context) {
	db, ok := c.MustGet("db").(*sqlx.DB)
	if !ok {
		log.Error("Failed to get database connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	tasks, err := model.GetAllTasks(db)
	if err != nil {
		log.Error("Failed to get tasks: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved tasks",
		"data":    tasks,
	})
}

func HandlerGetTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	db, ok := c.MustGet("db").(*sqlx.DB)
	if !ok {
		log.Error("Failed to get database connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	task, err := model.GetByID(db, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		log.Error("Failed to get task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func HandlerGetSubTasks(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	db, ok := c.MustGet("db").(*sqlx.DB)
	if !ok {
		log.Error("Failed to get database connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	task, err := model.GetTaskWithSubtasks(db, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		log.Error("Failed to get task subtasks", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task subtasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func HandlerCreateTasks(c *gin.Context) {
	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db, ok := c.MustGet("db").(*sqlx.DB)
	if !ok {
		log.Error("Failed to get database connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	task, err := model.CreateTask(db, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Error("Failed to create task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func HandlerUpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if everything is nil
	if !req.Title.Valid() && !req.Description.Valid() && !req.Status.Valid() && !req.Priority.Valid() &&
		!req.DueDate.Valid() && !req.CompletedAt.Valid() && !req.ParentTaskID.Valid() && !req.CategoryName.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	db, ok := c.MustGet("db").(*sqlx.DB)
	if !ok {
		log.Error("Failed to get database connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	effected, err := model.UpdateTask(db, taskID, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Unable to update Task"})
			return
		}
		log.Error("Failed to update task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	if effected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func HandlerDeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	db, ok := c.MustGet("db").(*sqlx.DB)
	if !ok {
		log.Error("Failed to get database connection")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	effected, err := model.DeleteTask(db, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		log.Error("Failed to delete task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if effected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
