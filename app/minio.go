package app

import (
	"log"

	"github.com/caffeines/filepile/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func SetMinioClient() error {
	minioCfg := config.GetMinio()
	endpoint := "localhost:9000"
	accessKeyID := minioCfg.AccessKey
	secretAccessKey := minioCfg.SecretKey
	useSSL := false
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}
	mc, err := minio.New(endpoint, opts)
	minioClient = mc
	if err != nil {
		return err
	}
	log.Println("Successfully connected to Minio")
	return nil
}

func GetMinioClient() *minio.Client {
	if minioClient == nil {
		panic("Minio not initialized")
	}
	return minioClient
}
