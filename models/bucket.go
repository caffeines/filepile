package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Bucket holds bucket data
type Bucket struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name,omitempty" json:"name"`
	CreatedBy primitive.ObjectID   `bson:"createdBy,omitempty" json:"createdBy"`
	Files     []primitive.ObjectID `bson:"files,omitempty" json:"files"`
	CreatedAt time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt  time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (b *Bucket) CollectionName() string {
	return "buckets"
}

func initBucketIndex(db *mongo.Database) error {
	bucket := Bucket{}
	bucketCol := db.Collection(bucket.CollectionName())
	if err := createIndex(bucketCol, bson.M{"name": 1, "CreatedBy": 1}, true); err != nil {
		return err
	}
	return nil
}
