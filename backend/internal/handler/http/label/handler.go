package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/label/create"
	"kanvix/backend/internal/usecase/label/delete"
	"kanvix/backend/internal/usecase/label/list"
	"kanvix/backend/internal/usecase/label/update"
)

type Handler struct {
	List   list.UseCase
	Create create.UseCase
	Update update.UseCase
	Delete delete.UseCase
}

type labelDTO struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
}

func toDTO(l entity.Label) labelDTO {
	return labelDTO{ID: l.ID, WorkspaceID: l.WorkspaceID, Name: l.Name, Color: l.Color}
}

func (h Handler) ListByWorkspaceHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	out, err := h.List.Execute(c.Request.Context(), list.In{UserID: userID, WorkspaceID: workspaceID})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	labels := make([]labelDTO, 0, len(out.Labels))
	for _, l := range out.Labels {
		labels = append(labels, toDTO(l))
	}
	c.JSON(http.StatusOK, gin.H{"labels": labels})
}

type createRequest struct {
	Name  string `json:"name" binding:"required,min=1"`
	Color string `json:"color" binding:"required"`
}

func (h Handler) CreateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.Create.Execute(c.Request.Context(), create.In{UserID: userID, WorkspaceID: workspaceID, Name: req.Name, Color: req.Color})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"label": toDTO(out.Label)})
}

type updateRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

func (h Handler) UpdateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	labelID := c.Param("labelId")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	patch := map[string]any{}
	if req.Name != nil {
		patch["name"] = *req.Name
	}
	if req.Color != nil {
		patch["color"] = *req.Color
	}
	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
		return
	}
	out, err := h.Update.Execute(c.Request.Context(), update.In{UserID: userID, LabelID: labelID, Patch: patch})
	if err != nil {
		if err.Error() == "empty patch" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
			return
		}
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"label": toDTO(out.Label)})
}

func (h Handler) DeleteHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	labelID := c.Param("labelId")
	if err := h.Delete.Execute(c.Request.Context(), delete.In{UserID: userID, LabelID: labelID}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

