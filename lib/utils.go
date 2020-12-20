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
	wce, ok := err.(mongo.CommandError)
	if !ok {
		return false
	}
	isDup := strings.Contains(err.Error(), " E11000 ")
	return wce.Code == 11000 || wce.Code == 11001 || wce.Code == 12582 || wce.Code == 16460 || isDup
}

// NewUUID return UUID/v4
func NewUUID() string {
	v := uuid.NewV4()
	return v.String()
}
