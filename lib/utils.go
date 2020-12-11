package lib

import "strings"

// IsDocumentNotFoundError returns boolean
func IsDocumentNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}
