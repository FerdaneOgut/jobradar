package main

import (
	"context"
	"log"
	"os"

	"github.com/FerdaneOgut/jobradar/ai-matching-service/providers"
	"github.com/FerdaneOgut/jobradar/ai-matching-service/repository"
	"github.com/FerdaneOgut/jobradar/ai-matching-service/services"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	// MongoDB
	mongoURI := os.Getenv("MONGODB_CONNECTION_STRING")
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database("jobradar")

	// LLM Provider
	llmProvider, err := providers.NewLLMProvider()
	if err != nil {
		log.Fatalf("Failed to create LLM provider: %v", err)
	}
	log.Printf("Using LLM provider: %s", llmProvider.Name())

	// Repository and Service
	matchRepo := repository.NewMatchRepository(db)
	matchingService := services.NewMatchingService(matchRepo, llmProvider)

	// Run matching immediately on startup
	log.Println("Running initial matching...")
	go matchingService.MatchAllUsers()

	// Schedule matching every 6 hours
	c := cron.New()
	c.AddFunc("0 */6 * * *", func() {
		log.Println("Running scheduled matching...")
		matchingService.MatchAllUsers()
	})
	c.Start()

	// HTTP server
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "ai-matching-service",
			"status":  "healthy",
		})
	})

	r.GET("/matches/:userId", func(c *gin.Context) {
		userID := c.Param("userId")
		matches, err := matchRepo.GetMatchesByUserID(context.Background(), userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, matches)
	})

	r.POST("/match", func(c *gin.Context) {
		go matchingService.MatchAllUsers()
		c.JSON(200, gin.H{"message": "Matching started"})
	})

	log.Println("AI Matching Service running on port 6082")
	r.Run(":6082")
}
