package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type WorkspacesHandler struct {
	Service services.AppService
}

type createWorkspaceRequest struct {
	Name  string  `json:"name" binding:"required,min=1"`
	Icon  *string `json:"icon"`
	Color *string `json:"color"`
}

type updateWorkspaceRequest struct {
	Name  *string `json:"name"`
	Icon  *string `json:"icon"`
	Color *string `json:"color"`
}

func (h WorkspacesHandler) List(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	ws, err := h.Service.ListWorkspaces(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workspaces": ws})
}

func (h WorkspacesHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	var req createWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	w, err := h.Service.CreateWorkspace(c.Request.Context(), userID, req.Name, req.Icon, req.Color)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"workspace": w})
}

func (h WorkspacesHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	id := c.Param("workspaceId")
	var req updateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	patch := map[string]any{}
	if req.Name != nil {
		patch["name"] = *req.Name
	}
	if req.Icon != nil {
		patch["icon"] = *req.Icon
	}
	if req.Color != nil {
		patch["color"] = *req.Color
	}
	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty patch"})
		return
	}
	w, err := h.Service.UpdateWorkspace(c.Request.Context(), userID, id, patch)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"workspace": w})
}

func (h WorkspacesHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	id := c.Param("workspaceId")
	if err := h.Service.DeleteWorkspace(c.Request.Context(), userID, id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
