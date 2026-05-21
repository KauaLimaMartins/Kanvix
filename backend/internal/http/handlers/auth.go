package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/config"
	"kanvix/backend/internal/repositories"
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

	u, token, err := h.Service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == repositories.ErrForbidden || err == repositories.ErrNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
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
			"role":        u.Role,
		},
	})
}

func (h AuthHandler) Setup(c *gin.Context) {
	need, err := h.Service.NeedsFirstSignup(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"needsFirstSignup": need})
}

type firstSignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func (h AuthHandler) FirstSignup(c *gin.Context) {
	var req firstSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	u, token, err := h.Service.FirstSignup(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if err == repositories.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "setup already completed"})
			return
		}
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

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":          u.ID,
			"email":       u.Email,
			"name":        u.Name,
			"avatarColor": u.AvatarColor,
			"role":        u.Role,
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
			"role":        u.Role,
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
