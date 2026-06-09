package consumers

import (
	"encoding/json"
	"log"

	"github.com/FerdaneOgut/jobradar/notification-service/models"
	"github.com/FerdaneOgut/jobradar/notification-service/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MatchConsumer struct {
	emailService *services.EmailService
}

func NewMatchConsumer(emailService *services.EmailService) *MatchConsumer {
	return &MatchConsumer{emailService: emailService}
}

func (c *MatchConsumer) Start(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"job.matches",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var notification models.MatchNotification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Error parsing match notification: %v", err)
				continue
			}

			log.Printf("Received match notification for %s - %s (score: %d)",
				notification.UserEmail, notification.JobTitle, notification.Score)

			if err := c.emailService.SendMatchNotification(
				notification.UserEmail,
				notification.UserName,
				notification.JobTitle,
				notification.Company,
				notification.JobURL,
				notification.Score,
				notification.Reason,
			); err != nil {
				log.Printf("Error sending match email: %v", err)
			}
		}
	}()

	log.Println("Match consumer started, waiting for messages...")
	return nil
}
