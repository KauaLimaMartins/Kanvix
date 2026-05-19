package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/http/middleware"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/services"
)

type BootstrapHandler struct {
	Service services.AppService
}

func (h BootstrapHandler) Get(c *gin.Context) {
	uval, ok := c.Get(middleware.CtxUserKey)
	user, _ := uval.(models.User)
	if !ok || user.ID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	out, err := h.Service.Bootstrap(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, out)
}
