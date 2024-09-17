package s3base

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 interface {
	Upload(ctx context.Context, file []byte, key string) error
	Download(ctx context.Context, key string) ([]byte, error)
}

type S3Client struct {
	Client *s3.Client
}

func NewS3Client(awsEndpoint string) (*S3Client, error) {
	awsRegion := "us-east-1"
	optsFunc := func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: awsEndpoint,
		}, nil
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(optsFunc)

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET_KEY", "TOKEN")),
	)
	if err != nil {
		log.Fatalf("failed to load config: error: %s", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{
		Client: s3Client,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, bucketName, key string, file []byte) error {
	uploader := manager.NewUploader(s.Client)
	if _, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	}); err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func (s *S3Client) DownloadFile(ctx context.Context, bucketName, objectKey, fileName string) error {
	result, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Printf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Printf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}

	_, err = file.Write(body)
	return err
}

func (s *S3Client) CreateBucket(bucketName string) error {
	if _, err := s.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}); err != nil {
		return err
	}

	return nil
}
