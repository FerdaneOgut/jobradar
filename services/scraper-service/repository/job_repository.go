package repository

import (
	"context"
	"time"

	"github.com/FerdaneOgut/jobradar/scraper-service/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobRepository struct {
	collection *mongo.Collection
	redis      *redis.Client
}

func NewJobRepository(db *mongo.Database, redis *redis.Client) *JobRepository {
	return &JobRepository{
		collection: db.Collection("jobs"),
		redis:      redis,
	}
}

func (r *JobRepository) IsExist(ctx context.Context, externalID string) (bool, error) {
	result, err := r.redis.Get(ctx, externalID).Result()
	if err == nil && result != "" {
		return true, nil
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"external_id": externalID})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *JobRepository) Save(ctx context.Context, job models.Job) error {
	_, err := r.collection.InsertOne(ctx, job)
	if err != nil {
		return err
	}

	r.redis.Set(ctx, job.ExternalID, true, 24*time.Hour)
	return nil
}

func (r *JobRepository) GetAll(ctx context.Context) ([]models.Job, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(100)
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
