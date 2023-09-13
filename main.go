package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Создаем новый экземпляр Fiber
	app := fiber.New()

	// Разрешаем CORS (Cross-Origin Resource Sharing)
	app.Use(cors.New())

	// Создаем роут "hello"
	app.Post("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	// Запускаем сервер на порту 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
