package main

import (
	"context"
	"log"
	"os"

	"github.com/FerdaneOgut/jobradar/scraper-service/repository"
	"github.com/FerdaneOgut/jobradar/scraper-service/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB
	mongoURI := os.Getenv("MONGODB_CONNECTION_STRING")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := mongoClient.Database("jobradar")

	// Redis
	redisAddr := os.Getenv("REDIS_CONNECTION")
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Repository and Service
	jobRepo := repository.NewJobRepository(db, redisClient)
	scraperService := services.NewScraperService(jobRepo)

	// Run scraper immediately on startup
	log.Println("Running initial scrape...")
	scraperService.ScrapeAll()

	// Schedule scraper every 6 hours
	c := cron.New()
	c.AddFunc("0 */6 * * *", func() {
		log.Println("Running scheduled scrape...")
		scraperService.ScrapeAll()
	})
	c.Start()

	// HTTP server
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "scraper-service",
			"status":  "healthy",
		})
	})

	r.GET("/jobs", func(c *gin.Context) {
		jobs, err := jobRepo.GetAll(context.Background())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, jobs)
	})

	r.GET("/scrape", func(c *gin.Context) {
		go scraperService.ScrapeAll()
		c.JSON(200, gin.H{"message": "Scraping started"})
	})

	log.Println("Scraper Service running on port 6081")
	r.Run(":6081")
}
