package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name"`
	BucketID  primitive.ObjectID `bson:"bucketId,omitempty" json:"bucketId"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy"`
	Size      int64              `bson:"size,omitempty" json:"size,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (f *File) CollectionName() string {
	return "files"
}
