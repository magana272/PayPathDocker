package storage

import (
	"context"

	"paypath/pkg/logger"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (d *DB) NextID(collection string) int {
	var result struct {
		Seq int `bson:"seq"`
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := d.db.Collection("counters").FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": collection},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&result)
	if err != nil {
		logger.Log.Fatal().Err(err).Str("collection", collection).Msg("failed to generate ID")
	}
	return result.Seq
}
