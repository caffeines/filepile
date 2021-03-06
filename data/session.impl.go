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
	_, err := sessionCollection.InsertOne(context.Background(), sess)
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
	err := sessionCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(sess)
	return sess, err
}

func (s *SessionRepoImpl) Logout(db *mongo.Database, token string) error {
	sess := &models.Session{}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	filter := bson.D{{"refreshToken", token}}
	err := sessionCollection.FindOneAndDelete(context.Background(), filter).Decode(&sess)
	return err
}
