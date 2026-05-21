package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type SubtasksHandler struct {
	Service services.AppService
}

type createSubtaskRequest struct {
	Title string `json:"title" binding:"required,min=1"`
}

type setSubtaskDoneRequest struct {
	Done *bool `json:"done"`
}

func (h SubtasksHandler) ListByTask(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	ss, err := h.Service.ListSubtasks(c.Request.Context(), userID, taskID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"subtasks": ss})
}

func (h SubtasksHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	var req createSubtaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	st, err := h.Service.CreateSubtask(c.Request.Context(), userID, taskID, req.Title)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"subtask": st})
}

func (h SubtasksHandler) SetDone(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	subtaskID := c.Param("subtaskId")
	var req setSubtaskDoneRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Done == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	st, err := h.Service.SetSubtaskDone(c.Request.Context(), userID, subtaskID, *req.Done)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"subtask": st})
}
