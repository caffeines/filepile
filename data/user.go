package data

import (
	"github.com/caffeines/sharehub/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository ...
type UserRepository interface {
	Register(db *mongo.Database, u *models.User) (*mongo.InsertOneResult, error)
}