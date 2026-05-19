package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type UsersHandler struct {
	Service services.AppService
}

func (h UsersHandler) List(c *gin.Context) {
	_, ok := requireUserID(c)
	if !ok {
		return
	}
	users, err := h.Service.ListUsers(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
