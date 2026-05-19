package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/repositories"
	"kanvix/backend/internal/services"
)

const CtxUserIDKey = "kanvix_user_id"
const CtxUserKey = "kanvix_user"

type Auth struct {
	CookieName string
	Service    services.AuthService
}

func (a Auth) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(a.CookieName)
		u, err := a.Service.Me(c.Request.Context(), token)
		if err != nil {
			status := http.StatusUnauthorized
			if err == repositories.ErrForbidden {
				status = http.StatusUnauthorized
			}
			c.AbortWithStatusJSON(status, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(CtxUserIDKey, u.ID)
		c.Set(CtxUserKey, u)
		c.Next()
	}
}
