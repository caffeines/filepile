package service

import "github.com/minio/minio-go/v7"

type MinioService interface {
	MakeBucket(minioClient *minio.Client, name string) (bool, error)
}
