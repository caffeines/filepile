package data

import (
	"context"
	"time"

	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *SessionRepoImpl) UpdateSession(db *mongo.Database, token, accessToken string) (*models.Session, error) {
	sess := &models.Session{}
	filter := bson.D{{"refreshToken", token}, {"expiresOn", bson.D{
		{"$gt", time.Now().Unix()},
	}}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.M{
		"refreshToken": lib.NewRefresToken(),
		"accesstoken":  accessToken,
		"createdAt":    time.Now().UTC(),
		"expiresOn":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}}}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := sessionCollection.FindOneAndUpdate(ctx, filter, update, &opt).Decode(sess)
	defer cancel()
	return sess, err
}
