package auth

import (
	"context"
	"fmt"
	"time"

	"paypath/internal/storage"
	"paypath/internal/storage/cache"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository interface {
	CreateUser(email, password, name string) (int64, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	DeleteUser(id int) (bool, error)
	RevokeToken(token string) error
	IsTokenRevoked(token string) (bool, error)
}

type mongoRepo struct {
	db *storage.DB
}

func NewRepository(db *storage.DB) Repository {
	r := &mongoRepo{db: db}
	ctx := context.Background()
	r.db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	r.db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	r.db.Collection("revoked_tokens").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "token", Value: 1}},
	})
	return r
}

func (r *mongoRepo) CreateUser(email, password, name string) (int64, error) {
	id := r.db.NextID("users")
	user := User{ID: id, Email: email, Password: password, Name: name}
	_, err := r.db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (r *mongoRepo) GetUserByEmail(email string) (*User, error) {
	var u User
	err := r.db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.db.Cache().Set(fmt.Sprintf("usr:%d", u.ID), &u, cache.DataCacheTTL)
	return &u, nil
}

func (r *mongoRepo) GetUserByID(id int) (*User, error) {
	key := fmt.Sprintf("usr:%d", id)
	if v, ok := r.db.Cache().Get(key); ok {
		return v.(*User), nil
	}
	var u User
	err := r.db.Collection("users").FindOne(context.Background(), bson.M{"id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.db.Cache().Set(key, &u, cache.DataCacheTTL)
	return &u, nil
}

func (r *mongoRepo) DeleteUser(id int) (bool, error) {
	res, err := r.db.Collection("users").DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return false, err
	}
	if res.DeletedCount > 0 {
		r.db.Cache().Delete(fmt.Sprintf("usr:%d", id))
	}
	return res.DeletedCount > 0, nil
}

func (r *mongoRepo) RevokeToken(token string) error {
	_, err := r.db.Collection("revoked_tokens").InsertOne(context.Background(), RevokedToken{
		Token:     token,
		RevokedAt: time.Now(),
	})
	if err == nil {
		r.db.Cache().Delete("tok:" + token)
	}
	return err
}

func (r *mongoRepo) IsTokenRevoked(token string) (bool, error) {
	key := "tok:" + token
	if v, ok := r.db.Cache().Get(key); ok {
		return v.(bool), nil
	}
	v, err, _ := r.db.SF().Do(key, func() (any, error) {
		if v, ok := r.db.Cache().Get(key); ok {
			return v, nil
		}
		err := r.db.Collection("revoked_tokens").FindOne(context.Background(), bson.M{"token": token}).Err()
		if err == mongo.ErrNoDocuments {
			r.db.Cache().Set(key, false, cache.TokenCacheTTL)
			return false, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	})
	if err != nil {
		return false, err
	}
	return v.(bool), nil
}
