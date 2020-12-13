package data

import (
	"context"
	"time"

	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepoImpl struct {
}

var sessionRepo SessionRepository

func NewSessionRepo() SessionRepository {
	if sessionRepo == nil {
		sessionRepo = &SessionRepoImpl{}
	}
	return sessionRepo
}

func (s *SessionRepoImpl) CreateSession(db *mongo.Database, sess *models.Session) error {
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := sessionCollection.InsertOne(ctx, sess)
	return err
}
