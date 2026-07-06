package main

import (
	"context"
	"fmt"
	"time"

	"paypath/pkg/setting"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	cfg := setting.Load()
	if cfg.MongoURI == "" {
		fmt.Println("MONGODB_URI not set")
		return
	}

	opts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(opts)
	if err != nil {
		fmt.Println("connect:", err)
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Database("paypath").Drop(ctx); err != nil {
		fmt.Println("drop:", err)
		return
	}
	fmt.Println("dropped paypath database — restart server to re-seed")
}
