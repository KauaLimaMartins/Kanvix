package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type ColumnsHandler struct {
	Service services.AppService
}

type createColumnRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

type updateColumnRequest struct {
	Name  *string `json:"name"`
	Order *int    `json:"order"`
}

func (h ColumnsHandler) ListByProject(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	cols, err := h.Service.ListColumns(c.Request.Context(), userID, projectID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"columns": cols})
}

func (h ColumnsHandler) Create(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	projectID := c.Param("projectId")
	var req createColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	col, err := h.Service.CreateColumn(c.Request.Context(), userID, projectID, req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"column": col})
}

func (h ColumnsHandler) Update(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	columnID := c.Param("columnId")
	var req updateColumnRequest
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
	col, err := h.Service.UpdateColumn(c.Request.Context(), userID, columnID, patch)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"column": col})
}

func (h ColumnsHandler) Delete(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	columnID := c.Param("columnId")
	if err := h.Service.DeleteColumn(c.Request.Context(), userID, columnID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
