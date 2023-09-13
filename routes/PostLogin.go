package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var request struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		IsRemember bool   `json:"isRemember"`
	}

	if err := c.BodyParser(&request); err != nil {
		fmt.Printf("Error parsing request body: %+v\n", request)
		fmt.Printf("Error parsing request body: %s\n", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	if request.Login == "" || request.Password == "" {
		fmt.Println("Login and password are required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Login and password are required"})
	}

	fmt.Printf("Received request data: %+v\n", request)

	apiURL := "https://poo.tomedu.ru/services/security/login"
	reqBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %s\n", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send request"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("External API returned an error: %s", resp.Status)
		fmt.Printf("Error: %s\n", errorMessage)
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": errorMessage})
	}

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding response: %s\n", err.Error())
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
