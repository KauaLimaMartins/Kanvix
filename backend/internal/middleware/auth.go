package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/usecase/auth/me"
)

const CtxUserIDKey = "kanvix_user_id"
const CtxUserKey = "kanvix_user"

type Auth struct {
	CookieName string
	Me         me.UseCase
}

func (a Auth) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(a.CookieName)
		out, err := a.Me.Execute(c.Request.Context(), token)
		if err != nil {
			status := http.StatusUnauthorized
			if err == domain.ErrForbidden {
				status = http.StatusUnauthorized
			}
			c.AbortWithStatusJSON(status, gin.H{"error": "unauthorized"})
			return
		}
		c.Set(CtxUserIDKey, out.User.ID)
		c.Set(CtxUserKey, out.User)
		c.Next()
	}
}
