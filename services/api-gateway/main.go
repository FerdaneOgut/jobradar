package main

import (
	"log"
	"os"

	"github.com/FerdaneOgut/jobradar/api-gateway/handlers"
	"github.com/FerdaneOgut/jobradar/api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Service URLs
	userService := os.Getenv("USER_SERVICE_URL")
	scraperService := os.Getenv("SCRAPER_SERVICE_URL")
	aiMatchingService := os.Getenv("AI_MATCHING_SERVICE_URL")
	trackerService := os.Getenv("TRACKER_SERVICE_URL")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "api-gateway",
			"status":  "healthy",
		})
	})

	// Public routes — no auth required
	r.POST("/api/auth/register", handlers.ProxyRequest(userService))
	r.POST("/api/auth/login", handlers.ProxyRequest(userService))

	// Protected routes — auth required
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	// User Service
	protected.POST("/api/cv/upload", handlers.ProxyRequest(userService))

	// Scraper Service
	protected.GET("/api/jobs", handlers.ProxyRequest(scraperService))
	protected.GET("/api/scrape", handlers.ProxyRequest(scraperService))

	// AI Matching Service
	protected.GET("/api/matches/:userId", handlers.ProxyRequest(aiMatchingService))
	protected.POST("/api/match", handlers.ProxyRequest(aiMatchingService))

	// Tracker Service
	protected.GET("/api/applications", handlers.ProxyRequest(trackerService))
	protected.POST("/api/applications", handlers.ProxyRequest(trackerService))
	protected.PUT("/api/applications/:id", handlers.ProxyRequest(trackerService))
	protected.DELETE("/api/applications/:id", handlers.ProxyRequest(trackerService))

	log.Println("API Gateway running on port 6080")
	r.Run(":6080")
}
