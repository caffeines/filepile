package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session model holds the session's data
type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	RefreshToken string             `bson:"refreshToken,omitempty" json:"refreshToken"`
	AccessToken  string             `bson:"accesstoken,omitempty" json:"accesstoken"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	ExpiresOn    int64              `bson:"expiresAt,omitempty" json:"expiresAt"`
}

// CollectionName returns name of the models
func (s Session) CollectionName() string {
	return "sessions"
}
