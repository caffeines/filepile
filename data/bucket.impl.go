package data

import (
	"context"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/models"
	"github.com/caffeines/filepile/service"
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
	isExists, bucketCreateErr := minioService.MakeBucket(minioClient, bucket.Name)

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
