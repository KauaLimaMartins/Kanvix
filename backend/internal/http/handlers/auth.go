package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/config"
	"kanvix/backend/internal/services"
)

type AuthHandler struct {
	Cfg     config.Config
	Service services.AuthService
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
}

func (h AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	u, token, err := h.Service.Login(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.Cfg.CookieName,
		Value:    token,
		Path:     "/",
		Domain:   h.Cfg.CookieDomain,
		HttpOnly: true,
		Secure:   h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.Cfg.SessionTTL.Seconds()),
		Expires:  time.Now().Add(h.Cfg.SessionTTL),
	})

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":          u.ID,
			"email":       u.Email,
			"name":        u.Name,
			"avatarColor": u.AvatarColor,
		},
	})
}

func (h AuthHandler) Me(c *gin.Context) {
	token, _ := c.Cookie(h.Cfg.CookieName)
	u, err := h.Service.Me(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":          u.ID,
			"email":       u.Email,
			"name":        u.Name,
			"avatarColor": u.AvatarColor,
		},
	})
}

func (h AuthHandler) Logout(c *gin.Context) {
	token, _ := c.Cookie(h.Cfg.CookieName)
	_ = h.Service.Logout(c.Request.Context(), token)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.Cfg.CookieName,
		Value:    "",
		Path:     "/",
		Domain:   h.Cfg.CookieDomain,
		HttpOnly: true,
		Secure:   h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
