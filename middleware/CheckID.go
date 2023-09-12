package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, "Something is bad #2")
			c.Abort()
			return
		}

		c.Next()
	}
}
