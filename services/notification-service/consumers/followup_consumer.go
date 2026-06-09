package consumers

import (
	"encoding/json"
	"log"

	"github.com/FerdaneOgut/jobradar/notification-service/models"
	"github.com/FerdaneOgut/jobradar/notification-service/services"
	amqp "github.com/rabbitmq/amqp091-go"
)

type FollowUpConsumer struct {
	emailService *services.EmailService
}

func NewFollowUpConsumer(emailService *services.EmailService) *FollowUpConsumer {
	return &FollowUpConsumer{emailService: emailService}
}

func (c *FollowUpConsumer) Start(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"job.followups",
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
			var notification models.FollowUpNotification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Error parsing followup notification: %v", err)
				continue
			}

			log.Printf("Received followup notification for %s - %s",
				notification.UserEmail, notification.JobTitle)

			if err := c.emailService.SendFollowUpNotification(
				notification.UserEmail,
				notification.UserName,
				notification.JobTitle,
				notification.Company,
				notification.JobURL,
			); err != nil {
				log.Printf("Error sending followup email: %v", err)
			}
		}
	}()

	log.Println("FollowUp consumer started, waiting for messages...")
	return nil
}
