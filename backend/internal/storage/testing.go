package storage

import (
	"context"
	"os"
	"time"

	"paypath/internal/storage/cache"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TestingT interface {
	Helper()
	Cleanup(func())
	Fatal(args ...any)
	Skip(args ...any)
}

func NewTest(t TestingT) *DB {
	t.Helper()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		t.Skip("MONGODB_URI not set; skipping integration test")
		return nil
	}
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		t.Fatal(err)
	}
	d := &DB{client: client, db: client.Database("paypath_test"), cache: cache.New()}
	t.Cleanup(func() {
		d.db.Drop(context.Background())
		client.Disconnect(context.Background())
	})
	return d
}
