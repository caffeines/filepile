package data

import (
	"context"
	"errors"
	"time"

	"github.com/caffeines/filepile/app"
	"github.com/minio/minio-go/v7"
)

type MinioImplRepo struct{}

var minioRepo MinioRepository

func NewMinioRepo() MinioRepository {
	if minioRepo == nil {
		minioRepo = &MinioImplRepo{}
	}
	return minioRepo
}

// MakeBucket create new bucket if not exists if exist return error and boolean
func (m *MinioImplRepo) MakeBucket(bucketName string) (bool, error) {
	minioClient := app.GetMinioClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return true, errors.New("bucketExist")
		}
		return exists, err
	}
	return false, nil
}
