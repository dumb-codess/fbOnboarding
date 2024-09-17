package s3base

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestS(t *testing.T) {
	awsEndpoint := "http://localhost:4566"
	bucketName := "test"
	s3client, err := NewS3Client(awsEndpoint)
	if err != nil {
		t.Fatalf("failed to create s3 client: %v", err)
	}

	testFile, err := os.Create("testfile")
	if err != nil {
		t.Fatalf("failed to create file")
	}

	if _, err := testFile.Write([]byte("this is test content")); err != nil {
		t.Fatal(err)
	}

	t.Run("creating bucket", func(t *testing.T) {
		if err := s3client.CreateBucket(bucketName); err != nil {
			t.Errorf("failed to create bucket %v", err)
			return
		}

		result, err := s3client.Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
		if err != nil {
			t.Errorf("failed to get List of buckets: %v", err)
			return
		}

		if GetPointerValue(result.Buckets[0].Name) != "test" {
			t.Errorf("wrong bucket name recieved got:%v want: %v ", GetPointerValue(result.Buckets[0].Name), "test")
		}

	})

	t.Run("uploading file to s3 bucket", func(t *testing.T) {
		fileBts, err := io.ReadAll(testFile)
		if err != nil {
			t.Errorf("failed to read file %v", err)
			return
		}
		if err := s3client.UploadFile(context.Background(), bucketName, "uploadedFiles/new", fileBts); err != nil {
			t.Errorf("failed to upload file: %v", err)
			return
		}

		// List Objects
		output, err := s3client.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			t.Errorf("failed to get object")
			return
		}

		if GetPointerValue(output.Contents[0].Key) != "uploadedFiles" {
			t.Errorf("Recieved wrong object ")
			return
		}

	})

	t.Run("download file from s3 bucket", func(t *testing.T) {
		if err := s3client.DownloadFile(context.Background(), bucketName, "uploadedFiles/new", "downloadedFile"); err != nil {
			t.Errorf("failed to download file from s3: %v", err)
			return
		}

	})

}

func GetPointerValue[T any](ptr *T) T {
	if ptr == nil {
		var x interface{}
		ZeroValue, _ := x.(T)
		return ZeroValue
	}
	return *ptr
}
