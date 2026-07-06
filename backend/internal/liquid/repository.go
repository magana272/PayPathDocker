package liquid

import (
	"context"
	"fmt"

	"paypath/internal/storage"
	"paypath/internal/storage/cache"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	All(userID int) ([]Liquid, error)
	Create(userID int, l Liquid) (Liquid, error)
	Update(userID, id int, l Liquid) (bool, error)
}

type mongoRepo struct {
	db *storage.DB
}

func NewRepository(db *storage.DB) Repository {
	r := &mongoRepo{db: db}
	r.db.Collection("liquid").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	return r
}

func cacheKey(userID int) string { return fmt.Sprintf("liq:%d", userID) }

func (r *mongoRepo) All(userID int) ([]Liquid, error) {
	return cache.CachedList(r.db.Cache(), r.db.SF(), cacheKey(userID), func() ([]Liquid, error) {
		cursor, err := r.db.Collection("liquid").Find(context.Background(), bson.M{"user_id": userID})
		if err != nil {
			return nil, err
		}
		var list []Liquid
		if err := cursor.All(context.Background(), &list); err != nil {
			return nil, err
		}
		return list, nil
	})
}

func (r *mongoRepo) Create(userID int, l Liquid) (Liquid, error) {
	l.ID = r.db.NextID("liquid")
	l.UserID = userID
	_, err := r.db.Collection("liquid").InsertOne(context.Background(), l)
	if err == nil {
		r.db.Cache().Delete(cacheKey(userID))
	}
	return l, err
}

func (r *mongoRepo) Update(userID, id int, l Liquid) (bool, error) {
	res, err := r.db.Collection("liquid").UpdateOne(
		context.Background(),
		bson.M{"id": id, "user_id": userID},
		bson.M{"$set": bson.M{
			"bank":    l.Bank,
			"balance": l.Balance,
		}},
	)
	if err != nil {
		return false, err
	}
	if res.ModifiedCount > 0 {
		r.db.Cache().Delete(cacheKey(userID))
	}
	return res.ModifiedCount > 0, nil
}
