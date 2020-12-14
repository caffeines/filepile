package config

import (
	"github.com/spf13/viper"
)

type Minio struct {
	AccessKey string
	SecretKey string
}

var minio Minio

func LoadMinio() {
	mu.Lock()
	defer mu.Unlock()
	minio = Minio{
		AccessKey: viper.GetString("minio.accessKey"),
		SecretKey: viper.GetString("minio.secretKey"),
	}
}

func GetMinio() Minio {
	return minio
}
