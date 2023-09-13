package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

func Login(c *fiber.Ctx) error {
	var request struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		IsRemember bool   `json:"isRemember"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Invalid login or password")
	}

	apiURL := "https://poo.tomedu.ru/services/security/login"
	reqBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

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

	setCookieHeader := resp.Header.Values("Set-Cookie")
	cookieString := strings.Join(setCookieHeader, "; ")

	successResponse := fiber.Map{
		"data":   data,
		"cookie": cookieString,
	}
	return c.JSON(successResponse)
}

func AddLoginRoute(router fiber.Router) {
	router.Post("/login", Login)
}
