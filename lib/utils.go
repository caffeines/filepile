package lib

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

// IsDocumentNotFoundError returns boolean
func IsDocumentNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

// NewUUID return UUID/v4
func NewUUID() string {
	v := uuid.NewV4()
	return v.String()
}
