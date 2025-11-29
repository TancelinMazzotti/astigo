package s3storage

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	globalContainer *S3Container
)

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()

	globalContainer, err = CreateS3Container(ctx)
	if err != nil {
		panic(err)
	}

	if err := os.Setenv("AWS_ENDPOINT_URL", globalContainer.Endpoint.String()); err != nil {
		panic(err)
	}
	if err := os.Setenv("AWS_ACCESS_KEY_ID", "minio-root-user"); err != nil {
		panic(err)
	}
	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", "minio-root-password"); err != nil {
		panic(err)
	}
	if err := os.Setenv("AWS_REGION", "us-east-1"); err != nil {
		panic(err)
	}

	client, err := NewS3(ctx, globalContainer.Config)
	if err != nil {
		panic(err)
	}

	_, err = client.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &globalContainer.Config.Bucket,
	})
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
