package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type ProjectsHandler struct {
	Service services.AppService
}

type createProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=1"`
	Description *string `json:"description"`
}

type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (h ProjectsHandler) ListByWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	ps, err := h.Service.ListProjects(c.Request.Context(), userID, workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": ps})
}

func (h ProjectsHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	p, cols, err := h.Service.CreateProject(c.Request.Context(), userID, workspaceID, req.Name, req.Description)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"project": p, "columns": cols})
}

func (h ProjectsHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	var req updateProjectRequest
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
	p, err := h.Service.UpdateProject(c.Request.Context(), userID, projectID, patch)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"project": p})
}

func (h ProjectsHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	if err := h.Service.DeleteProject(c.Request.Context(), userID, projectID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
