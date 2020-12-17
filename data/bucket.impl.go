package data

import (
	"context"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/models"
	"github.com/caffeines/filepile/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type StorageRepoImpl struct{}

var storageRepo StorageRepository

func NewStorageRepo() StorageRepository {
	if storageRepo == nil {
		storageRepo = &StorageRepoImpl{}
	}
	return storageRepo
}

func (s *StorageRepoImpl) CreateNewBucket(db *mongo.Database, bucket *models.Bucket) (bool, error) {
	bucketCollection := db.Collection(bucket.CollectionName())

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	client := db.Client()
	session, err := client.StartSession()
	if err != nil {
		return false, err
	}
	defer session.EndSession(context.Background())

	insertionError := mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(txnOpts); err != nil {
			return err
		}
		_, err := bucketCollection.InsertOne(sessionContext, bucket)
		if err != nil {
			return err
		}
		return nil
	})

	minioClient := app.GetMinioClient()
	minioService := service.NewMinioService()
	isExists, bucketCreateErr := minioService.MakeBucket(minioClient, bucket.ID.String())

	if insertionError != nil || bucketCreateErr != nil || isExists {
		var abortErr error
		if lib.IsMongoDupKey(insertionError) || isExists {
			abortErr = session.AbortTransaction(context.Background())
			return true, abortErr
		}
		abortErr = session.AbortTransaction(context.Background())
		if abortErr != nil {
			return false, abortErr
		} else if bucketCreateErr != nil {
			return false, bucketCreateErr
		}
		return false, insertionError

	}
	if err := session.CommitTransaction(context.Background()); err != nil {
		return false, err
	}
	return false, nil
}

// FindBucketsByCreatorID returns folders by creator userId
func (s *SessionRepoImpl) FindBucketsByCreatorID(db *mongo.Database, creatorID primitive.ObjectID, lastID string) ([]models.Bucket, error) {
	bucket := models.Bucket{}
	bucketCollection := db.Collection(bucket.CollectionName())
	opts := options.Find().SetSort(bson.M{"_id": 1}).SetLimit(10)
	query := bson.M{"createdBy": creatorID}
	if lastID == "" {
		id, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		query = bson.M{"createdBy": creatorID, "_id": bson.M{"$gt": id}}
	}
	cursor, err := bucketCollection.Find(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}
	var buckets []models.Bucket
	if err = cursor.All(context.Background(), &buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}
