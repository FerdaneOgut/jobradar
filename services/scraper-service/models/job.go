package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ExternalID  string             `bson:"external_id" json:"external_id"`
	Title       string             `bson:"title" json:"title"`
	Company     string             `bson:"company" json:"company"`
	Location    string             `bson:"location" json:"location"`
	Description string             `bson:"description" json:"description"`
	URL         string             `bson:"url" json:"url"`
	Source      string             `bson:"source" json:"source"`
	Tags        []string           `bson:"tags" json:"tags"`
	IsRemote    bool               `bson:"is_remote" json:"is_remote"`
	Salary      string             `bson:"salary" json:"salary"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
