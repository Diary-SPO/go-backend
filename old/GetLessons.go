//package routes
//
//import (
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"time"
//
//	"github.com/gin-gonic/gin"
//)
//
//func GetStudentLessons(c *gin.Context) {
//	id := c.Param("id")
//	startDate := c.Param("startDate")
//	endDate := c.Param("endDate")
//	secret := c.GetHeader("secret")
//
//	var formattedStartDate, formattedEndDate string
//
//	if startDate != "" && endDate != "" {
//		formattedStartDate = startDate
//		formattedEndDate = endDate
//
//		startTime, _ := time.Parse("2006-01-02", startDate)
//		endTime, _ := time.Parse("2006-01-02", endDate)
//		differenceInDays := endTime.Sub(startTime).Hours() / 24
//
//		if differenceInDays > 14 {
//			newEndDate := startTime.Add(14 * 24 * time.Hour)
//			formattedEndDate = newEndDate.Format("2006-01-02")
//		}
//	} else {
//		currentDate := time.Now()
//		formattedStartDate = currentDate.Format("2006-01-02")
//
//		endDate := currentDate.Add(14 * 24 * time.Hour)
//		formattedEndDate = endDate.Format("2006-01-02")
//	}
//
//	apiURL := fmt.Sprintf("https://poo.tomedu.ru/services/students/%s/lessons/%s/%s", id, formattedStartDate, formattedEndDate)
//	req, err := http.NewRequest("GET", apiURL, nil)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
//		return
//	}
//
//	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
//	req.Header.Set("Cookie", secret)
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
//		return
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		errorMessage := fmt.Sprintf("External API returned an error: %s", resp.Status)
//		c.JSON(resp.StatusCode, gin.H{"error": errorMessage})
//		return
//	}
//
//	var data interface{}
//	err = json.NewDecoder(resp.Body).Decode(&data)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
//		return
//	}
//
//	c.JSON(http.StatusOK, data)
//}
//
//func AddStudentLessonsRoute(router *gin.RouterGroup) {
//	lessonsRouter := router.Group("/lessons")
//	lessonsRouter.GET("/:id/:startDate/:endDate", GetStudentLessons)
//}
