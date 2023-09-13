package middleware

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func CheckID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if id == "" {
			return c.Status(http.StatusBadRequest).JSON("Something is bad #2")
		}

		return c.Next()
	}
}
