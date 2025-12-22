package s3storage

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/minio"
)

type S3Container struct {
	*minio.MinioContainer
	Config   Config
	Endpoint url.URL
	Host     string
	Port     int
}

func CreateS3Container(ctx context.Context) (*S3Container, error) {
	var err error
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

	return &S3Container{
		MinioContainer: s3Container,
		Config: Config{
			Bucket:       "default",
			UsePathStyle: true,
			Timeout:      10 * time.Second,
		},
		Endpoint: url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort("localhost", containerPort.Port()),
		},
		Host: "localhost",
		Port: containerPort.Int(),
	}, nil
}
