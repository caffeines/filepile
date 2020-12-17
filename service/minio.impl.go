package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
)

type MinioServiceImpl struct{}

var minioService MinioService

func NewMinioService() MinioService {
	if minioService == nil {
		minioService = &MinioServiceImpl{}
	}
	return minioService
}

// MakeBucket create new bucket if not exists if exist return error and boolean
func (m *MinioServiceImpl) MakeBucket(minioClient *minio.Client, bucketName string) (bool, error) {
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			return true, errors.New("bucketExist")
		}
		return exists, err
	}
	return false, nil
}

func (m *MinioServiceImpl) UploadToMinio(mc *minio.Client, bucket, fileName, contentType string, reader io.Reader, size int64) error {
	info, err := mc.PutObject(context.Background(), bucket, fileName, reader, size, minio.PutObjectOptions{
		ContentType:        contentType,
		ContentDisposition: fmt.Sprintf("attachment; filename=\"%s\"", fileName),
	})
	log.Println(info)
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioServiceImpl) GetObjectFromMinio(mc *minio.Client, bucket, fileName string) (*minio.Object, error) {
	obj, err := mc.GetObject(context.Background(), bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}
