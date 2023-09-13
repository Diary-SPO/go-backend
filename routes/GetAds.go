package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Diary-SPO/go-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func GetNotifications(c *fiber.Ctx) error {
	secret := c.Get("secret")

	apiURL := "https://poo.tomedu.ru/services/people/organization/news/last/10"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send request"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("External API returned an error: %s", resp.Status)
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": errorMessage})
	}

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode response"})
	}

	return c.JSON(data)
}

func AddNotificationsRoute(router fiber.Router) {
	router.Get("/ads", GetNotifications, middleware.CheckCookie())
}
