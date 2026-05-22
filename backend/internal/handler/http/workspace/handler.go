package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/workspace/create"
	"kanvix/backend/internal/usecase/workspace/delete"
	"kanvix/backend/internal/usecase/workspace/list"
	"kanvix/backend/internal/usecase/workspace/update"
)

type Handler struct {
	List   list.UseCase
	Create create.UseCase
	Update update.UseCase
	Delete delete.UseCase
}

type workspaceDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
	Role  string `json:"role,omitempty"`
}

func toDTO(w entity.Workspace) workspaceDTO {
	return workspaceDTO{ID: w.ID, Name: w.Name, Icon: w.Icon, Color: w.Color}
}

func (h Handler) ListHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	out, err := h.List.Execute(c.Request.Context(), userID)
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	ws := make([]workspaceDTO, 0, len(out.Workspaces))
	for _, w := range out.Workspaces {
		ws = append(ws, toDTO(w))
	}
	c.JSON(http.StatusOK, gin.H{"workspaces": ws})
}

type createRequest struct {
	Name  string  `json:"name" binding:"required,min=1"`
	Icon  *string `json:"icon"`
	Color *string `json:"color"`
}

func (h Handler) CreateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.Create.Execute(c.Request.Context(), create.In{
		UserID: userID,
		Name:   req.Name,
		Icon:   req.Icon,
		Color:  req.Color,
	})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"workspace": toDTO(out.Workspace)})
}

type updateRequest struct {
	Name  *string `json:"name"`
	Icon  *string `json:"icon"`
	Color *string `json:"color"`
}

func (h Handler) UpdateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	id := c.Param("workspaceId")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.Update.Execute(c.Request.Context(), update.In{
		UserID:      userID,
		WorkspaceID: id,
		Name:        req.Name,
		Icon:        req.Icon,
		Color:       req.Color,
	})
	if err != nil {
		if err.Error() == "empty patch" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
			return
		}
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workspace": toDTO(out.Workspace)})
}

func (h Handler) DeleteHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	id := c.Param("workspaceId")
	if err := h.Delete.Execute(c.Request.Context(), delete.In{UserID: userID, WorkspaceID: id}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

