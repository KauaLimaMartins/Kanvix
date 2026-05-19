package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/http/middleware"
)

func requireUserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", false
	}
	id, _ := v.(string)
	if id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", false
	}
	return id, true
}
