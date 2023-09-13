package main

import (
	"github.com/Diary-SPO/go-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"time"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Expiration: 10 * time.Second,
		Max:        10,
	}))

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
