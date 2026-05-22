package statssearch

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/handler/http/httputil"
	searchws "kanvix/backend/internal/usecase/search/workspace"
	statsws "kanvix/backend/internal/usecase/stats/workspace"
)

type Handler struct {
	Stats  statsws.UseCase
	Search searchws.UseCase
}

func (h Handler) WorkspaceStatsHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
	if !ok {
		return
	}
	workspaceID := c.Param("workspaceId")
	out, err := h.Stats.Execute(c.Request.Context(), userID, workspaceID)
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h Handler) SearchHandler(c *gin.Context) {
	userID, ok := httputil.RequireUserID(c)
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
	out, err := h.Search.Execute(c.Request.Context(), userID, workspaceID, q, limit)
	if err != nil {
		httputil.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

