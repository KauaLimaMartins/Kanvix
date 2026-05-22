package project

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/handler/http/httputil"
	"kanvix/backend/internal/usecase/project/create"
	"kanvix/backend/internal/usecase/project/delete"
	"kanvix/backend/internal/usecase/project/list"
	"kanvix/backend/internal/usecase/project/update"
)

type Handler struct {
	List   list.UseCase
	Create create.UseCase
	Update update.UseCase
	Delete delete.UseCase
}

type projectDTO struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type columnDTO struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
}

func toProjectDTO(p entity.Project) projectDTO {
	return projectDTO{ID: p.ID, WorkspaceID: p.WorkspaceID, Name: p.Name, Description: p.Description}
}

func toColumnDTO(c entity.Column) columnDTO {
	return columnDTO{ID: c.ID, ProjectID: c.ProjectID, Name: c.Name, Order: c.Order}
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
	projects := make([]projectDTO, 0, len(out.Projects))
	for _, p := range out.Projects {
		projects = append(projects, toProjectDTO(p))
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

type createRequest struct {
	Name string `json:"name" binding:"required,min=1"`
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
	out, err := h.Create.Execute(c.Request.Context(), create.In{UserID: userID, WorkspaceID: workspaceID, Name: req.Name})
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	cols := make([]columnDTO, 0, len(out.Columns))
	for _, col := range out.Columns {
		cols = append(cols, toColumnDTO(col))
	}
	c.JSON(http.StatusCreated, gin.H{"project": toProjectDTO(out.Project), "columns": cols})
}

type updateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h Handler) UpdateHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	patch := map[string]any{}
	if req.Name != nil {
		patch["name"] = *req.Name
	}
	if req.Description != nil {
		patch["description"] = *req.Description
	}
	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
		return
	}
	out, err := h.Update.Execute(c.Request.Context(), update.In{UserID: userID, ProjectID: projectID, Patch: patch})
	if err != nil {
		if err.Error() == "empty patch" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
			return
		}
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": toProjectDTO(out.Project)})
}

func (h Handler) DeleteHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	if err := h.Delete.Execute(c.Request.Context(), delete.In{UserID: userID, ProjectID: projectID}); err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

