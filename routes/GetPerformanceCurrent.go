package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPerformanceCurrent(c *gin.Context) {
	id := c.Param("id")
	secret := c.GetHeader("secret")

	apiURL := fmt.Sprintf("https://poo.tomedu.ru/services/reports/current/performance/%s", id)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Cookie", secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("External API returned an error: %s", resp.Status)
		c.JSON(resp.StatusCode, gin.H{"error": errorMessage})
		return
	}

	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func AddPerformanceRoutes(router *gin.RouterGroup) {
	performanceRouter := router.Group("/performance.current")
	performanceRouter.GET("/:id", GetPerformanceCurrent)
}
