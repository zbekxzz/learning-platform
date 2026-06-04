package storage

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	// Ensure standard buckets are created and have public read access
	EnsureBucketPublic("certificates")
	EnsureBucketPublic("materials")
}

func EnsureBucketPublic(bucketName string) {
	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, bucketName)
	if err == nil && !exists {
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Println("Warning: failed to make bucket:", bucketName, err)
			return
		}
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)

	err = Client.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		log.Println("Warning: failed to set public bucket policy for:", bucketName, err)
	}
}

func UploadFile(bucket, objectName string, filePath string) (string, error) {
	opts := minio.PutObjectOptions{}

	// Set appropriate Content-Type and Content-Disposition to force download for certificates
	if strings.HasSuffix(strings.ToLower(objectName), ".svg") {
		opts.ContentType = "image/svg+xml"
		opts.ContentDisposition = "attachment; filename=\"" + objectName + "\""
	} else if strings.HasSuffix(strings.ToLower(objectName), ".pdf") {
		opts.ContentType = "application/pdf"
	} else if strings.HasSuffix(strings.ToLower(objectName), ".png") {
		opts.ContentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(objectName), ".jpg") || strings.HasSuffix(strings.ToLower(objectName), ".jpeg") {
		opts.ContentType = "image/jpeg"
	}

	_, err := Client.FPutObject(context.Background(),
		bucket,
		objectName,
		filePath,
		opts)

	if err != nil {
		return "", err
	}

	url := "http://localhost:9000/" + bucket + "/" + objectName

	return url, nil
}

