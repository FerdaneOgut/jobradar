package models

import "time"

type NotificationType string

const (
	NotificationTypeMatch    NotificationType = "match"
	NotificationTypeFollowUp NotificationType = "followup"
)

type MatchNotification struct {
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	JobTitle  string `json:"job_title"`
	Company   string `json:"company"`
	JobURL    string `json:"job_url"`
	Score     int    `json:"score"`
	Reason    string `json:"reason"`
}

type FollowUpNotification struct {
	UserEmail string    `json:"user_email"`
	UserName  string    `json:"user_name"`
	JobTitle  string    `json:"job_title"`
	Company   string    `json:"company"`
	JobURL    string    `json:"job_url"`
	AppliedAt time.Time `json:"applied_at"`
}
