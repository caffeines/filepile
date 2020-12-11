package data

import (
	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository ...
type UserRepository interface {
	Register(db *mongo.Database, u *models.User) (*mongo.InsertOneResult, error)
	FindUserByEmail(db *mongo.Database, email string) (*models.User, error)
}
