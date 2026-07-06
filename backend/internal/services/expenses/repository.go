package expenses

import (
	"context"
	"fmt"

	"paypath/internal/storage"
	"paypath/internal/storage/cache"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository interface {
	All(userID int) ([]Expense, error)
	Create(userID int, e Expense) (Expense, error)
	Update(userID, id int, e Expense) (*Expense, error)
	Delete(userID, id int) (bool, error)
}

type mongoRepo struct {
	db *storage.DB
}

func NewRepository(db *storage.DB) Repository {
	r := &mongoRepo{db: db}
	r.db.Collection("expenses").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	return r
}

func cacheKey(userID int) string { return fmt.Sprintf("exp:%d", userID) }

func (r *mongoRepo) All(userID int) ([]Expense, error) {
	return cache.CachedList(r.db.Cache(), r.db.SF(), cacheKey(userID), func() ([]Expense, error) {
		cursor, err := r.db.Collection("expenses").Find(context.Background(), bson.M{"user_id": userID})
		if err != nil {
			return nil, err
		}
		var list []Expense
		if err := cursor.All(context.Background(), &list); err != nil {
			return nil, err
		}
		return list, nil
	})
}

func (r *mongoRepo) Create(userID int, e Expense) (Expense, error) {
	e.ID = r.db.NextID("expenses")
	e.UserID = userID
	_, err := r.db.Collection("expenses").InsertOne(context.Background(), e)
	if err == nil {
		r.db.Cache().Delete(cacheKey(userID))
	}
	return e, err
}

func (r *mongoRepo) Update(userID, id int, e Expense) (*Expense, error) {
	set := bson.M{}
	if e.Expense != "" {
		set["expense"] = e.Expense
	}
	if e.Cost != 0 {
		set["cost"] = e.Cost
	}
	if e.Date != nil {
		set["date"] = e.Date
	}
	if e.DueDate != nil {
		set["due_date"] = e.DueDate
	}
	if e.Frequency != "" {
		set["frequency"] = e.Frequency
	}
	if e.Exceptions != nil {
		set["exceptions"] = e.Exceptions
	}
	if len(set) == 0 {
		return r.get(userID, id)
	}
	filter := bson.M{"id": id, "user_id": userID}
	res, err := r.db.Collection("expenses").UpdateOne(context.Background(), filter, bson.M{"$set": set})
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, nil
	}
	r.db.Cache().Delete(cacheKey(userID))
	return r.get(userID, id)
}

func (r *mongoRepo) get(userID, id int) (*Expense, error) {
	var e Expense
	err := r.db.Collection("expenses").FindOne(context.Background(), bson.M{"id": id, "user_id": userID}).Decode(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *mongoRepo) Delete(userID, id int) (bool, error) {
	res, err := r.db.Collection("expenses").DeleteOne(context.Background(), bson.M{"id": id, "user_id": userID})
	if err != nil {
		return false, err
	}
	if res.DeletedCount > 0 {
		r.db.Cache().Delete(cacheKey(userID))
	}
	return res.DeletedCount > 0, nil
}
