package storage

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func InitMinIO() {

	endpoint := "localhost:9000"
	accessKey := "admin"
	secretKey := "password"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal(err)
	}

	Client = client
}

func UploadFile(bucket, objectName string, filePath string) (string, error) {

	_, err := Client.FPutObject(context.Background(),
		bucket,
		objectName,
		filePath,
		minio.PutObjectOptions{})

	if err != nil {
		return "", err
	}

	url := "http://localhost:9000/" + bucket + "/" + objectName

	return url, nil
}
