package services

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (s *EmailService) SendMatchNotification(to, name, jobTitle, company, jobURL string, score int, reason string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("🎯 New Job Match: %s at %s (Score: %d)", jobTitle, company, score))
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Hi %s!</h2>
		<p>We found a great job match for you!</p>
		<h3>%s at %s</h3>
		<p><strong>Match Score:</strong> %d/100</p>
		<p><strong>Why it matches:</strong> %s</p>
		<p><a href="%s">View Job</a></p>
	`, name, jobTitle, company, score, reason, jobURL))

	return s.send(m)
}

func (s *EmailService) SendFollowUpNotification(to, name, jobTitle, company, jobURL string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("⏰ Follow-up reminder: %s at %s", jobTitle, company))
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Hi %s!</h2>
		<p>It's been 7 days since you applied to <strong>%s at %s</strong>.</p>
		<p>Consider sending a follow-up email to check on your application status.</p>
		<p><a href="%s">View Job</a></p>
	`, name, jobTitle, company, jobURL))

	return s.send(m)
}

func (s *EmailService) send(m *gomail.Message) error {
	if s.host == "" {
		log.Println("SMTP not configured, skipping email")
		return nil
	}

	d := gomail.NewDialer(s.host, s.port, s.username, s.password)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email sent to %s", m.GetHeader("To")[0])
	return nil
}
