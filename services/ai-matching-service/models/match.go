package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Job struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ExternalID  string        `bson:"external_id" json:"external_id"`
	Title       string        `bson:"title" json:"title"`
	Company     string        `bson:"company" json:"company"`
	Location    string        `bson:"location" json:"location"`
	Description string        `bson:"description" json:"description"`
	URL         string        `bson:"url" json:"url"`
	Source      string        `bson:"source" json:"source"`
	Tags        []string      `bson:"tags" json:"tags"`
	IsRemote    bool          `bson:"is_remote" json:"is_remote"`
	Salary      string        `bson:"salary" json:"salary"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
}

type User struct {
	ID     bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Email  string        `bson:"email" json:"email"`
	CvPath string        `bson:"cvPath" json:"cvPath"`
}

type MatchResult struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string        `bson:"user_id" json:"user_id"`
	JobID     string        `bson:"job_id" json:"job_id"`
	Job       Job           `bson:"job" json:"job"`
	Score     int           `bson:"score" json:"score"`
	Reason    string        `bson:"reason" json:"reason"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
