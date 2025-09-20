package s3storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

type S3Container struct {
	*minio.MinioContainer
	Config Config
}

func CreateS3Container(ctx context.Context) (*S3Container, error) {
	var err error
	config := Config{
		Endpoint:        "http://localhost",
		Region:          "us-east-1",
		Bucket:          "default",
		AccessKeyID:     "minio-root-user",
		SecretAccessKey: "minio-root-password",
		UsePathStyle:    true,
		Timeout:         10 * time.Second,
	}

	s3Container, err := minio.Run(ctx, "quay.io/minio/minio",
		minio.WithUsername("minio-root-user"),
		minio.WithPassword("minio-root-password"),
	)
	if err != nil {
		return nil, err
	}

	containerPort, err := s3Container.MappedPort(ctx, "9000/tcp")
	if err != nil {
		return nil, err
	}
	config.Endpoint += ":" + containerPort.Port()

	client, err := NewS3(ctx, config)
	if err != nil {
		return nil, err
	}

	_, err = client.raw.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &config.Bucket,
	})
	if err != nil {
		return nil, err
	}

	return &S3Container{
		MinioContainer: s3Container,
		Config:         config,
	}, nil
}
