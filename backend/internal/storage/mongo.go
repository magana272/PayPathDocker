package storage

import (
	"context"
	"sync"
	"time"

	"paypath/internal/storage/cache"
	"paypath/pkg/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/singleflight"
)

const minPoolSize = 8

type DB struct {
	client *mongo.Client
	db     *mongo.Database
	cache  *cache.Cache
	sf     singleflight.Group
}

func Connect(uri string) *DB {
	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(20)
	client, err := mongo.Connect(opts)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to connect to MongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed to ping MongoDB")
	}

	warmPool(client, minPoolSize)

	logger.Log.Info().Msg("connected to MongoDB")
	return &DB{client: client, db: client.Database("paypath"), cache: cache.New()}
}

func warmPool(client *mongo.Client, n int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := client.Ping(ctx, nil); err != nil {
				logger.Log.Warn().Err(err).Msg("connection pool warm-up ping failed")
			}
		}()
	}
	wg.Wait()
}

func (d *DB) Close() error {
	return d.client.Disconnect(context.Background())
}

func (d *DB) Collection(name string) *mongo.Collection {
	return d.db.Collection(name)
}

func (d *DB) Cache() *cache.Cache {
	return d.cache
}

func (d *DB) SF() *singleflight.Group {
	return &d.sf
}
