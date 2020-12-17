package lib

import (
	"strings"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// IsDocumentNotFoundError returns boolean
func IsDocumentNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

func IsMongoDupKey(err error) bool {
	wce, ok := err.(mongo.WriteConcernError)
	if !ok {
		return false
	}
	return wce.Code == 11000 || wce.Code == 11001 || wce.Code == 12582 || wce.Code == 16460 && strings.Contains(wce.Message, " E11000 ")
}

// NewUUID return UUID/v4
func NewUUID() string {
	v := uuid.NewV4()
	return v.String()
}
