package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User model holds the user's data
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name,omitempty" json:"name"`
	Username         string             `bson:"username,omitempty" json:"username"`
	Email            string             `bson:"email,omitempty" json:"email"`
	Password         string             `bson:"password,omitempty" json:"-"`
	ProfilePic       string             `bson:"profilePic,omitempty" json:"profilePic,omitempty"`
	VerificationCode string             `bson:"verificationCode,omitempty" json:"-"`
	CodeGenTime      time.Time          `bson:"codeGenTime,omitempty" json:"-"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt         time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (u User) CollectionName() string {
	return "users"
}

func initUserIndex(db *mongo.Database) error {
	user := User{}
	userCol := db.Collection(user.CollectionName())
	if err := createIndex(userCol, bson.M{"username": 1}, true); err != nil {
		return err
	}
	if err := createIndex(userCol, bson.M{"email": 1}, true); err != nil {
		return err
	}
	return nil
}
