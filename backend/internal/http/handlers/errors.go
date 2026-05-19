package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/repositories"
)

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repositories.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(err, repositories.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
