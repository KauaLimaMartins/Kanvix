package subtask

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/subtask/create"
	"kanvix/backend/internal/usecase/subtask/delete"
	"kanvix/backend/internal/usecase/subtask/list"
	"kanvix/backend/internal/usecase/subtask/update"
)

type Handler struct {
	List   list.UseCase
	Create create.UseCase
	Update update.UseCase
	Delete delete.UseCase
}

type subtaskDTO struct {
	ID        string `json:"id"`
	TaskID    string `json:"taskId"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}

func toDTO(s entity.Subtask) subtaskDTO {
	return subtaskDTO{
		ID:        s.ID,
		TaskID:    s.TaskID,
		Title:     s.Title,
		Done:      s.Done,
		CreatedAt: s.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func (h Handler) ListByTaskHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	out, err := h.List.Execute(c.Request.Context(), list.In{UserID: userID, TaskID: taskID})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	ss := make([]subtaskDTO, 0, len(out.Subtasks))
	for _, s := range out.Subtasks {
		ss = append(ss, toDTO(s))
	}
	c.JSON(http.StatusOK, gin.H{"subtasks": ss})
}

type createRequest struct {
	Title string `json:"title" binding:"required,min=1"`
}

func (h Handler) CreateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	taskID := c.Param("taskId")
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.Create.Execute(c.Request.Context(), create.In{UserID: userID, TaskID: taskID, Title: req.Title})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"subtask": toDTO(out.Subtask)})
}

type updateRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

func (h Handler) PatchHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	subtaskID := c.Param("subtaskId")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Title == nil && req.Done == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
		return
	}
	if req.Title != nil {
		t := strings.TrimSpace(*req.Title)
		if t == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		req.Title = &t
	}

	out, err := h.Update.Execute(c.Request.Context(), update.In{UserID: userID, SubtaskID: subtaskID, Title: req.Title, Done: req.Done})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"subtask": toDTO(out.Subtask)})
}

func (h Handler) DeleteHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	subtaskID := c.Param("subtaskId")
	if err := h.Delete.Execute(c.Request.Context(), delete.In{UserID: userID, SubtaskID: subtaskID}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

