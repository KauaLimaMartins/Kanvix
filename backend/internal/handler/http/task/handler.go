package task

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/task/create"
	"kanvix/backend/internal/usecase/task/delete"
	"kanvix/backend/internal/usecase/task/get"
	"kanvix/backend/internal/usecase/task/list"
	"kanvix/backend/internal/usecase/task/move"
	"kanvix/backend/internal/usecase/task/update"
)

type Handler struct {
	List   list.UseCase
	Get    get.UseCase
	Create create.UseCase
	Update update.UseCase
	Delete delete.UseCase
	Move   move.UseCase
}

type taskDTO struct {
	ID          string   `json:"id"`
	ProjectID   string   `json:"projectId"`
	ColumnID    string   `json:"columnId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	DueDate     *string  `json:"dueDate,omitempty"`
	AssigneeID  *string  `json:"assigneeId,omitempty"`
	Order       int      `json:"order"`
	CreatedAt   string   `json:"createdAt"`
}

func toDTO(t entity.Task, labels []string) taskDTO {
	if labels == nil {
		labels = []string{}
	}
	return taskDTO{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		ColumnID:    t.ColumnID,
		Title:       t.Title,
		Description: t.Description,
		Labels:      labels,
		DueDate:     t.DueDate,
		AssigneeID:  t.AssigneeID,
		Order:       t.Order,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func (h Handler) ListByProjectHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	out, err := h.List.Execute(c.Request.Context(), list.In{UserID: userID, ProjectID: projectID})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	tasks := make([]taskDTO, 0, len(out.Tasks))
	for _, t := range out.Tasks {
		tasks = append(tasks, toDTO(t, out.Labels[t.ID]))
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

type createRequest struct {
	ColumnID string `json:"columnId" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

func (h Handler) CreateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.Create.Execute(c.Request.Context(), create.In{
		UserID:    userID,
		ProjectID: projectID,
		ColumnID:  req.ColumnID,
		Title:     req.Title,
	})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": toDTO(out.Task, []string{})})
}

func (h Handler) GetHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	out, err := h.Get.Execute(c.Request.Context(), get.In{UserID: userID, TaskID: taskID})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": toDTO(out.Task, out.Labels)})
}

type updateRequest struct {
	Title       *string               `json:"title"`
	Description *string               `json:"description"`
	DueDate     httputil.OptionalString `json:"dueDate"`
	AssigneeID  httputil.OptionalString `json:"assigneeId"`
	Labels      httputil.OptionalStringSlice `json:"labels"`
}

func (h Handler) UpdateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	var req updateRequest
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

	out, err := h.Update.Execute(c.Request.Context(), update.In{
		UserID: userID,
		TaskID: taskID,
		Patch:  patch,
		Labels: labels,
	})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": toDTO(out.Task, out.Labels)})
}

func (h Handler) DeleteHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	if err := h.Delete.Execute(c.Request.Context(), delete.In{UserID: userID, TaskID: taskID}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type moveRequest struct {
	ToColumnID string `json:"toColumnId"`
	ToIndex    *int   `json:"toIndex"`
}

func (h Handler) MoveHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	var req moveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.ToColumnID == "" || req.ToIndex == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Move.Execute(c.Request.Context(), move.In{UserID: userID, TaskID: taskID, ToColumnID: req.ToColumnID, ToIndex: *req.ToIndex}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

