package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/caffeines/sharehub/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var instance *mongo.Client

// ConnectMongo connects with MongoDB
func ConnectMongo() error {
	db := config.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.URL))
	instance = client
	if err != nil {
		return err
	}
	if err := instance.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("\nSuccessfully connected to MongoDB")
	return nil
}

// DisconnectMongo disconnects with MongoDB
func DisconnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := instance.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

// GetMongoClient returns created mongo instance or error if not connected
func GetMongoClient() (*mongo.Client, error) {
	if instance == nil {
		return nil, errors.New("Database not connected")
	}
	return instance, nil
}
