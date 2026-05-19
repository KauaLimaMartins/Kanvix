package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type LabelsHandler struct {
	Service services.AppService
}

type createLabelRequest struct {
	Name  string `json:"name" binding:"required,min=1"`
	Color string `json:"color" binding:"required,min=1"`
}

type updateLabelRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

func (h LabelsHandler) ListByWorkspace(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	ls, err := h.Service.ListLabels(c.Request.Context(), userID, workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"labels": ls})
}

func (h LabelsHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	var req createLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	l, err := h.Service.CreateLabel(c.Request.Context(), userID, workspaceID, req.Name, req.Color)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"label": l})
}

func (h LabelsHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	labelID := c.Param("labelId")
	var req updateLabelRequest
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
	l, err := h.Service.UpdateLabel(c.Request.Context(), userID, labelID, patch)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"label": l})
}

func (h LabelsHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	labelID := c.Param("labelId")
	if err := h.Service.DeleteLabel(c.Request.Context(), userID, labelID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
