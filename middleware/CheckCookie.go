package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func CheckCookie() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secret := c.Get("secret")

		if secret == "" {
			return c.Status(http.StatusBadRequest).JSON("Something is bad #1")
		}

		return c.Next()
	}
}
