package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type TasksHandler struct {
	Service services.AppService
}

type createTaskRequest struct {
	ColumnID string `json:"columnId" binding:"required"`
	Title    string `json:"title" binding:"required,min=1"`
}

type updateTaskRequest struct {
	Title       *string             `json:"title"`
	Description *string             `json:"description"`
	DueDate     OptionalString      `json:"dueDate"`
	AssigneeID  OptionalString      `json:"assigneeId"`
	Labels      OptionalStringSlice `json:"labels"`
}

type moveTaskRequest struct {
	ToColumnID string `json:"toColumnId" binding:"required"`
	ToIndex    *int   `json:"toIndex"`
}

func (h TasksHandler) ListByProject(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	tasks, err := h.Service.ListTasks(c.Request.Context(), userID, projectID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h TasksHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	t, err := h.Service.CreateTask(c.Request.Context(), userID, projectID, req.ColumnID, req.Title)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": t})
}

func (h TasksHandler) Get(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	t, err := h.Service.GetTask(c.Request.Context(), userID, taskID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": t})
}

func (h TasksHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	var req updateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	patch := map[string]any{}
	if req.Title != nil {
		patch["title"] = *req.Title
	}
	if req.Description != nil {
		patch["description"] = *req.Description
	}
	if req.DueDate.Set {
		if req.DueDate.Value == nil {
			patch["due_date"] = nil
		} else {
			patch["due_date"] = *req.DueDate.Value
		}
	}
	if req.AssigneeID.Set {
		if req.AssigneeID.Value == nil {
			patch["assignee_id"] = nil
		} else {
			patch["assignee_id"] = *req.AssigneeID.Value
		}
	}

	var labels *[]string
	if req.Labels.Set {
		labels = &req.Labels.Value
	}

	if len(patch) == 0 && labels == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
		return
	}

	t, err := h.Service.UpdateTask(c.Request.Context(), userID, taskID, patch, labels)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": t})
}

func (h TasksHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	if err := h.Service.DeleteTask(c.Request.Context(), userID, taskID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h TasksHandler) Move(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	var req moveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.ToColumnID == "" || req.ToIndex == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Service.MoveTask(c.Request.Context(), userID, taskID, req.ToColumnID, *req.ToIndex); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
