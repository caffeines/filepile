package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model holds the user's data
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"name,omitempty"`
	Username         string             `bson:"username,omitempty"`
	Email            string             `bson:"email,omitempty"`
	Password         string             `bson:"password,omitempty"`
	ProfilePic       string             `bson:"profilePic,omitempty"`
	VerificationCode string             `bson:"verificationCode,omitempty"`
	CodeGenTime      time.Time          `bson:"codeGenTime,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty"`
	UpdateAt         time.Time          `bson:"updatedAt,omitempty"`
}

// CollectionName returns name of the models
func (u User) CollectionName() string {
	return "users"
}
