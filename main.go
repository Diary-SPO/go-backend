package main

import (
	"github.com/Diary-SPO/go-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"time"
)

var blockedIPs = make(map[string]time.Time)

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     os.Stdout,
	}))

	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{
		Max:        70,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			ip := c.IP()
			blockedIPs[ip] = time.Now().Add(5 * time.Minute)
			return c.Status(fiber.StatusTooManyRequests).SendString("Rate limit exceeded.")
		},
	}))

	// Middleware для проверки блокировки клиентов.
	app.Use(func(c *fiber.Ctx) error {
		ip := c.IP()
		if _, ok := blockedIPs[ip]; ok {
			if time.Now().After(blockedIPs[ip]) {
				delete(blockedIPs, ip)
			} else {
				return c.Status(fiber.StatusTooManyRequests).SendString("You are temporarily blocked.")
			}
		}
		return c.Next()
	})

	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	routes.AddLoginRoute(app)
	routes.AddAttestationRoute(app)
	routes.AddPerformanceRoute(app)
	routes.AddOrganizationRoute(app)
	routes.AddNotificationsRoute(app)
	routes.AddStudentLessonsRoute(app)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
