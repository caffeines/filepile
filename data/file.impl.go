package data

import (
	"context"
	"io"

	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type FileRepoImpl struct{}

var fileRepo FileRepository

func NewFileRepo() FileRepository {
	if fileRepo == nil {
		fileRepo = &FileRepoImpl{}
	}
	return fileRepo
}

func (f *FileRepoImpl) UploadFile(file *models.File, reader io.Reader, db *mongo.Database) error {
	fileCollection := db.Collection(file.CollectionName())

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	client := db.Client()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	insertionError := mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(txnOpts); err != nil {
			return err
		}
		_, err := fileCollection.InsertOne(sessionContext, file)
		if err != nil {
			return err
		}
		//TODO: add file in bucket also
		return nil
	})
	if insertionError != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			return abortErr
		}
		return insertionError
	}
	return nil
}
