package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader("secret")

		if secret == "" {
			c.JSON(http.StatusBadRequest, "Something is bad #1")
			c.Abort()
			return
		}

		c.Next()
	}
}
