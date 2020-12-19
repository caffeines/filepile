package service

import (
	"io"

	"github.com/minio/minio-go/v7"
)

type MinioService interface {
	MakeBucket(minioClient *minio.Client, name string) (bool, error)
	UploadToMinio(bucket, fileName, contentType string, reader io.Reader, size int64, client *minio.Client) error
}
