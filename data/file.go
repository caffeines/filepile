package data

import (
	"io"

	"github.com/caffeines/filepile/models"
	"github.com/minio/minio-go/v7"
)

type FileRepository interface {
	UploadFile(file *models.File, reader io.Reader, client *minio.Client) error
}
