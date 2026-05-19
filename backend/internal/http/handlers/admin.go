package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/services"
)

type AdminHandler struct {
	Service services.AppService
}

func (h AdminHandler) ResetDemo(c *gin.Context) {
	_, ok := requireUserID(c)
	if !ok {
		return
	}
	if err := h.Service.ResetDemo(c.Request.Context()); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
