package data

import (
	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageRepository interface {
	CreateNewBucket(db *mongo.Database, bucket *models.Bucket) (bool, error)
}
