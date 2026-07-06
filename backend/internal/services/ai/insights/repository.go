package insights

import (
	"context"
	"time"

	"paypath/internal/storage"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository interface {
	GetCached(userID int, topic string) (string, error)
	Save(userID int, topic, response string) error
}

type mongoRepo struct {
	db *storage.DB
}

func NewRepository(db *storage.DB) Repository {
	r := &mongoRepo{db: db}
	r.db.Collection("insights_cache").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "topic", Value: 1}},
	})
	return r
}

func (r *mongoRepo) GetCached(userID int, topic string) (string, error) {
	var ic InsightsCache
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := r.db.Collection("insights_cache").FindOne(context.Background(), bson.M{"user_id": userID, "topic": topic}, opts).Decode(&ic)
	if err != nil {
		return "", err
	}
	return ic.Response, nil
}

func (r *mongoRepo) Save(userID int, topic, response string) error {
	_, err := r.db.Collection("insights_cache").InsertOne(context.Background(), InsightsCache{
		UserID:    userID,
		Topic:     topic,
		Response:  response,
		CreatedAt: time.Now(),
	})
	return err
}
