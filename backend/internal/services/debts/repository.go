package debts

import (
	"context"
	"fmt"

	"paypath/internal/storage"
	"paypath/internal/storage/cache"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	All(userID int) ([]Debt, error)
	Create(userID int, d Debt) (Debt, error)
	Update(userID, id int, d Debt) (bool, error)
	Delete(userID, id int) (bool, error)
}

type mongoRepo struct {
	db *storage.DB
}

func NewRepository(db *storage.DB) Repository {
	r := &mongoRepo{db: db}
	r.db.Collection("debts").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	return r
}

func cacheKey(userID int) string { return fmt.Sprintf("dbt:%d", userID) }

func (r *mongoRepo) All(userID int) ([]Debt, error) {
	return cache.CachedList(r.db.Cache(), r.db.SF(), cacheKey(userID), func() ([]Debt, error) {
		cursor, err := r.db.Collection("debts").Find(context.Background(), bson.M{"user_id": userID})
		if err != nil {
			return nil, err
		}
		var list []Debt
		if err := cursor.All(context.Background(), &list); err != nil {
			return nil, err
		}
		return list, nil
	})
}

func (r *mongoRepo) Create(userID int, d Debt) (Debt, error) {
	d.ID = r.db.NextID("debts")
	d.UserID = userID
	_, err := r.db.Collection("debts").InsertOne(context.Background(), d)
	if err == nil {
		r.db.Cache().Delete(cacheKey(userID))
	}
	return d, err
}

func (r *mongoRepo) Update(userID, id int, d Debt) (bool, error) {
	res, err := r.db.Collection("debts").UpdateOne(
		context.Background(),
		bson.M{"id": id, "user_id": userID},
		bson.M{"$set": bson.M{
			"bank":    d.Bank,
			"type":    d.Type,
			"name":    d.Name,
			"apy":     d.APY,
			"balance": d.Balance,
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

func (r *mongoRepo) Delete(userID, id int) (bool, error) {
	res, err := r.db.Collection("debts").DeleteOne(context.Background(), bson.M{"id": id, "user_id": userID})
	if err != nil {
		return false, err
	}
	if res.DeletedCount > 0 {
		r.db.Cache().Delete(cacheKey(userID))
	}
	return res.DeletedCount > 0, nil
}
