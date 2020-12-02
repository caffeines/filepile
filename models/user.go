package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
