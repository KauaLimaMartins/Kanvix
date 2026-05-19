package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type StatsSearchHandler struct {
	Service services.AppService
}

func (h StatsSearchHandler) WorkspaceStats(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	out, err := h.Service.WorkspaceStats(c.Request.Context(), userID, workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h StatsSearchHandler) Search(c *gin.Context) {
	userID, ok := requireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	q := c.Query("q")
	limit := 0
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	out, err := h.Service.Search(c.Request.Context(), userID, workspaceID, q, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
