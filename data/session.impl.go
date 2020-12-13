package data

import (
	"context"
	"time"

	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *SessionRepoImpl) UpdateSession(db *mongo.Database, token string) (*models.Session, error) {
	sess := &models.Session{}
	filter := bson.D{{"refreshToken", token}}
	update := bson.D{{"$set", bson.M{
		"refreshToken": lib.NewRefresToken(),
		"createdAt":    time.Now().UTC(),
		"expiresOn":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}}}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := sessionCollection.FindOneAndUpdate(ctx, filter, update).Decode(sess)
	defer cancel()
	return sess, err
}
