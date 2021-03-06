package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createIndex(collection *mongo.Collection, keys interface{}, unique bool) error {
	opts := options.Index().SetUnique(unique)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    keys,
		Options: opts,
	})
	return err
}

// InitializeIndex populates all collections indexes
func InitializeIndex(db *mongo.Database) error {
	if err := initBucketIndex(db); err != nil {
		return err
	}
	if err := initUserIndex(db); err != nil {
		return err
	}
	if err := initSessionIndex(db); err != nil {
		return err
	}
	return nil
}
