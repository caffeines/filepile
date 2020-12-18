package data

import (
	"io"

	"github.com/minio/minio-go/v7"
)

type FileRepository interface {
	UploadToMinio(bucket, fileName, contentType string, reader io.Reader, size int64, client *minio.Client) error
}
