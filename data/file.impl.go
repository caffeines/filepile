package data

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type FileRepoImpl struct{}

var fileRepo FileRepository

func NewFileRepo() FileRepository {
	if fileRepo == nil {
		fileRepo = &FileRepoImpl{}
	}
	return fileRepo
}

func (f *FileRepoImpl) UploadToMinio(bucket, fileName, contentType string, reader io.Reader, size int64, client *minio.Client) error {
	_, err := client.PutObject(context.Background(), bucket, fileName, reader, size, minio.PutObjectOptions{
		ContentType:        contentType,
		ContentDisposition: fmt.Sprintf("attachment; filename=\"%s\"", fileName),
	})
	if err != nil {
		return err
	}
	return nil
}
