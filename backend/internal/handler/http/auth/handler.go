package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/config"
	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/usecase/auth/login"
	"kanvix/backend/internal/usecase/auth/logout"
	"kanvix/backend/internal/usecase/auth/me"
	"kanvix/backend/internal/usecase/auth/register"
	"kanvix/backend/internal/usecase/auth/setup"
)

type Handler struct {
	Cfg config.Config

	Setup    setup.UseCase
	Register register.UseCase
	Login    login.UseCase
	Me       me.UseCase
	Logout   logout.UseCase
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password"`
}

func (h Handler) LoginHandler(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	out, err := h.Login.Execute(c.Request.Context(), login.In{Email: req.Email, Password: req.Password})
	if err != nil {
		if err == domain.ErrForbidden || err == domain.ErrNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.Cfg.CookieName,
		Value:    out.Token,
		Path:     "/",
		Domain:   h.Cfg.CookieDomain,
		HttpOnly: true,
		Secure:   h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.Cfg.SessionTTL.Seconds()),
		Expires:  time.Now().Add(h.Cfg.SessionTTL),
	})

	c.JSON(http.StatusOK, gin.H{"user": out.User})
}

func (h Handler) SetupHandler(c *gin.Context) {
	out, err := h.Setup.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, out)
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

func (h Handler) FirstSignupHandler(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	out, err := h.Register.Execute(c.Request.Context(), register.In{Email: req.Email, Password: req.Password, Name: req.Name})
	if err != nil {
		if err == domain.ErrForbidden {
			c.JSON(http.StatusForbidden, gin.H{"error": "setup already completed"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.Cfg.CookieName,
		Value:    out.Token,
		Path:     "/",
		Domain:   h.Cfg.CookieDomain,
		HttpOnly: true,
		Secure:   h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.Cfg.SessionTTL.Seconds()),
		Expires:  time.Now().Add(h.Cfg.SessionTTL),
	})

	c.JSON(http.StatusCreated, gin.H{"user": out.User})
}

func (h Handler) MeHandler(c *gin.Context) {
	token, _ := c.Cookie(h.Cfg.CookieName)
	out, err := h.Me.Execute(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h Handler) LogoutHandler(c *gin.Context) {
	token, _ := c.Cookie(h.Cfg.CookieName)
	_ = h.Logout.Execute(c.Request.Context(), token)
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
