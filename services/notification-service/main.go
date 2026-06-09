package main

import (
	"log"
	"os"
	"time"

	"github.com/FerdaneOgut/jobradar/notification-service/consumers"
	"github.com/FerdaneOgut/jobradar/notification-service/services"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func connectRabbitMQ() *amqp.Connection {
	rabbitURL := os.Getenv("RABBITMQ_CONNECTION")
	var conn *amqp.Connection
	var err error

	for i := 0; i < 20; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			log.Println("Connected to RabbitMQ")
			return conn
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5s... (%d/20)", i+1)
		time.Sleep(5 * time.Second)
	}

	log.Fatalf("Failed to connect to RabbitMQ after 20 attempts: %v", err)
	return nil
}

func main() {
	// Connect to RabbitMQ
	conn := connectRabbitMQ()
	defer conn.Close()

	// Email service
	emailService := services.NewEmailService()

	// Start consumers
	matchConsumer := consumers.NewMatchConsumer(emailService)
	if err := matchConsumer.Start(conn); err != nil {
		log.Fatalf("Failed to start match consumer: %v", err)
	}

	followUpConsumer := consumers.NewFollowUpConsumer(emailService)
	if err := followUpConsumer.Start(conn); err != nil {
		log.Fatalf("Failed to start followup consumer: %v", err)
	}

	// HTTP server
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "notification-service",
			"status":  "healthy",
		})
	})

	log.Println("Notification Service running on port 6083")
	r.Run(":6083")
}
