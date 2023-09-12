package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var request struct {
		Login      string `json:"login"`
		Password   string `json:"password"`
		IsRemember bool   `json:"isRemember"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid login or password")
		return
	}

	apiURL := "https://poo.tomedu.ru/services/security/login"
	reqBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Получение текста ошибки от внешнего сервера
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

	setCookieHeader := resp.Header.Values("Set-Cookie")
	cookieString := strings.Join(setCookieHeader, "; ")

	c.JSON(http.StatusOK, gin.H{"data": data, "cookie": cookieString})
}

func AddLoginRoute(router *gin.RouterGroup) {
	router.POST("/login", Login)
}
