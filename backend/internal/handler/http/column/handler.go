package column

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/column/create"
	"kanvix/backend/internal/usecase/column/delete"
	"kanvix/backend/internal/usecase/column/list"
	"kanvix/backend/internal/usecase/column/update"
)

type Handler struct {
	List   list.UseCase
	Create create.UseCase
	Update update.UseCase
	Delete delete.UseCase
}

type columnDTO struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
}

func toDTO(c2 entity.Column) columnDTO {
	return columnDTO{ID: c2.ID, ProjectID: c2.ProjectID, Name: c2.Name, Order: c2.Order}
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
	cols := make([]columnDTO, 0, len(out.Columns))
	for _, col := range out.Columns {
		cols = append(cols, toDTO(col))
	}
	c.JSON(http.StatusOK, gin.H{"columns": cols})
}

type createRequest struct {
	Name string `json:"name" binding:"required,min=1"`
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
	out, err := h.Create.Execute(c.Request.Context(), create.In{UserID: userID, ProjectID: projectID, Name: req.Name})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"column": toDTO(out.Column)})
}

type updateRequest struct {
	Name  *string `json:"name"`
	Order *int    `json:"order"`
}

func (h Handler) UpdateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	columnID := c.Param("columnId")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	patch := map[string]any{}
	if req.Name != nil {
		patch["name"] = *req.Name
	}
	if req.Order != nil {
		patch["order"] = *req.Order
	}
	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
		return
	}
	out, err := h.Update.Execute(c.Request.Context(), update.In{UserID: userID, ColumnID: columnID, Patch: patch})
	if err != nil {
		if err.Error() == "empty patch" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
			return
		}
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"column": toDTO(out.Column)})
}

func (h Handler) DeleteHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	columnID := c.Param("columnId")
	if err := h.Delete.Execute(c.Request.Context(), delete.In{UserID: userID, ColumnID: columnID}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

