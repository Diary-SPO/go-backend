package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/scffs/go-backend/middleware"
	"github.com/scffs/go-backend/routes"
	"net"
	"net/http"
	"time"
)

var (
	limiter     = ratelimit.NewBucket(time.Minute, 80)
	banList     = make(map[string]time.Time)
	banDuration = 15 * time.Minute
)

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"secret", "Content-Type"}

	router.Use(cors.New(config))

	router.GET("/", handleRequest)

	performanceGroup := router.Group("/")
	performanceGroup.Use(middleware.CheckID())
	performanceGroup.Use(middleware.CheckCookie())
	routes.AddPerformanceRoutes(performanceGroup)

	organizationGroup := router.Group("/")
	organizationGroup.Use(middleware.CheckCookie())
	routes.AddOrganizationRoute(organizationGroup)

	adsGroup := router.Group("/")
	adsGroup.Use(middleware.CheckCookie())
	routes.AddNotificationsRoute(adsGroup)

	attestationGroup := router.Group("/")
	attestationGroup.Use(middleware.CheckID())
	attestationGroup.Use(middleware.CheckCookie())
	routes.AddGroupAttestationRoute(attestationGroup)

	lessonsGroup := router.Group("/")
	lessonsGroup.Use(middleware.CheckCookie())
	routes.AddStudentLessonsRoute(lessonsGroup)

	loginGroup := router.Group("/")
	routes.AddLoginRoute(loginGroup)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func handleRequest(c *gin.Context) {
	clientIP, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	if _, exists := banList[clientIP]; exists {
		c.String(http.StatusTooManyRequests, "You are temporarily banned")
		return
	}

	if limiter.TakeAvailable(1) > 0 {
		c.String(http.StatusOK, "Request successful")
	} else {
		banList[clientIP] = time.Now()
		c.String(http.StatusTooManyRequests, "Request limit exceeded. You are temporarily banned")
	}
}
