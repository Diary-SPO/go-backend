package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Diary-SPO/go-backend/middleware"
	"github.com/gofiber/fiber/v2"
)

func GetStudentLessons(c *fiber.Ctx) error {
	id := c.Params("id")
	startDate := c.Params("startDate")
	endDate := c.Params("endDate")
	secret := c.Get("secret")

	var formattedStartDate, formattedEndDate string

	if startDate != "" && endDate != "" {
		formattedStartDate = startDate
		formattedEndDate = endDate

		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		differenceInDays := endTime.Sub(startTime).Hours() / 24

		if differenceInDays > 14 {
			newEndDate := startTime.Add(14 * 24 * time.Hour)
			formattedEndDate = newEndDate.Format("2006-01-02")
		}
	} else {
		currentDate := time.Now()
		formattedStartDate = currentDate.Format("2006-01-02")

		endDate := currentDate.Add(14 * 24 * time.Hour)
		formattedEndDate = endDate.Format("2006-01-02")
	}

	apiURL := fmt.Sprintf("https://poo.tomedu.ru/services/students/%s/lessons/%s/%s", id, formattedStartDate, formattedEndDate)
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

func AddStudentLessonsRoute(router fiber.Router) {
	lessonsRouter := router.Group("/lessons")
	lessonsRouter.Get("/:id/:startDate/:endDate", middleware.CheckID(), middleware.CheckCookie(), GetStudentLessons)
}
