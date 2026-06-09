package repository

import (
	"context"
	"log"
	"time"

	"github.com/FerdaneOgut/jobradar/ai-matching-service/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MatchRepository struct {
	matches *mongo.Collection
	jobs    *mongo.Collection
	users   *mongo.Collection
}

func NewMatchRepository(db *mongo.Database) *MatchRepository {
	return &MatchRepository{
		matches: db.Collection("matches"),
		jobs:    db.Collection("jobs"),
		users:   db.Collection("users"),
	}
}

func (r *MatchRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	cursor, err := r.users.Find(ctx, bson.M{"cvPath": bson.M{"$exists": true, "$ne": ""}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MatchRepository) GetRecentJobs(ctx context.Context) ([]models.Job, error) {
	opts := options.Find().SetLimit(100)

	// Try recent jobs first (last 6 hours)
	since := time.Now().Add(-6 * time.Hour)
	cursor, err := r.jobs.Find(ctx, bson.M{"created_at": bson.M{"$gte": since}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	// If no recent jobs, fall back to all jobs
	if len(jobs) == 0 {
		log.Println("No recent jobs found, fetching all jobs...")
		cursor2, err := r.jobs.Find(ctx, bson.M{}, opts)
		if err != nil {
			return nil, err
		}
		defer cursor2.Close(ctx)

		if err := cursor2.All(ctx, &jobs); err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func (r *MatchRepository) GetMatchesByUserID(ctx context.Context, userID string) ([]models.MatchResult, error) {
	opts := options.Find().
		SetSort(bson.M{"score": -1}).
		SetLimit(20)

	cursor, err := r.matches.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var matches []models.MatchResult
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *MatchRepository) SaveMatch(ctx context.Context, match models.MatchResult) error {
	_, err := r.matches.InsertOne(ctx, match)
	return err
}
