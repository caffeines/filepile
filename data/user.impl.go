package data

import (
	"context"

	"github.com/caffeines/filepile/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepoImpl ...
type UserRepoImpl struct {
}

var userRepo UserRepository

func NewUserRepo() UserRepository {
	if userRepo == nil {
		userRepo = &UserRepoImpl{}
	}
	return userRepo
}

// Register creates new user
func (usr *UserRepoImpl) Register(db *mongo.Database, user *models.User) (*mongo.InsertOneResult, error) {
	userCollection := db.Collection(user.CollectionName())
	createdUser, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

//FindUserByEmail returns the matched document with email
func (usr *UserRepoImpl) FindUserByEmail(db *mongo.Database, email string) (*models.User, error) {
	user := models.User{}
	userCollection := db.Collection(user.CollectionName())
	if err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

//FindUserByID returns the matched document with id
func (usr *UserRepoImpl) FindUserByID(db *mongo.Database, id primitive.ObjectID) (*models.User, error) {
	user := models.User{}
	userCollection := db.Collection(user.CollectionName())
	if err := userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
