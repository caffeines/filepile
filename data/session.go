package data

import (
	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository interface {
	CreateSession(db *mongo.Database, sess *models.Session) error
	UpdateSession(db *mongo.Database, token, accessToken string) (*models.Session, error)
	Logout(db *mongo.Database, token string) error
}
