package s3storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Config holds configuration parameters for connecting to and interacting with an S3-compatible storage service.
type Config struct {
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	SessionToken    string `mapstructure:"session_token"`

	UsePathStyle         bool   `mapstructure:"use_path_style"`
	ServerSideEncryption string `mapstructure:"server_side_encryption"`
	DefaultCacheControl  string `mapstructure:"default_cache_control"`

	UploadPartSize int64         `mapstructure:"upload_part_size"`
	Timeout        time.Duration `mapstructure:"timeout"`
}

// Client represents an S3 client with configuration for AWS SDK for Go, providing upload, download, and presigning capabilities.
type Client struct {
	config     Config
	awsconfig  aws.Config
	raw        *s3.Client
	uploader   *manager.Uploader
	downloader *manager.Downloader
	presign    *s3.PresignClient
}

// NewS3 initializes a new S3 client with the provided configuration and returns it, or an error if the configuration is invalid.
func NewS3(ctx context.Context, config Config) (*Client, error) {
	loadOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(config.Region),
	}

	if config.AccessKeyID != "" && config.SecretAccessKey != "" {
		loadOpts = append(loadOpts, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, config.SessionToken),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), loadOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Options client S3
	s3Opts := func(o *s3.Options) {
		o.UsePathStyle = config.UsePathStyle
		if config.Endpoint != "" {
			o.BaseEndpoint = aws.String(config.Endpoint)
		}
	}
	raw := s3.NewFromConfig(awsCfg, s3Opts)

	// Uploader/Downloader
	upOpts := func(u *manager.Uploader) {
		if config.UploadPartSize > 0 {
			u.PartSize = config.UploadPartSize
		}
	}

	uploader := manager.NewUploader(raw, func(o *manager.Uploader) {})
	upOpts(uploader)

	downloader := manager.NewDownloader(raw)
	presignClient := s3.NewPresignClient(raw)

	client := &Client{
		config:     config,
		awsconfig:  awsCfg,
		raw:        raw,
		uploader:   uploader,
		downloader: downloader,
		presign:    presignClient,
	}

	if err := client.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping S3 storage: %w", err)
	}

	return client, nil
}

// Put uploads an object to an S3 bucket with the specified key and optional metadata such as content type and cache control.
func (c *Client) Put(ctx context.Context, bucket, key string, body io.Reader, contentType, cacheControl string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: nil,
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}

	if cacheControl != "" {
		input.CacheControl = aws.String(cacheControl)
	} else if c.config.DefaultCacheControl != "" {
		input.CacheControl = aws.String(c.config.DefaultCacheControl)
	}

	if c.config.ServerSideEncryption != "" {
		input.ServerSideEncryption = types.ServerSideEncryption(c.config.ServerSideEncryption)
	}

	_, err := c.raw.PutObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

// PutLarge uploads a large object to the specified S3 bucket using multipart upload for improved performance.
// Supports optional content type and cache control metadata.
func (c *Client) PutLarge(ctx context.Context, bucket, key string, body io.Reader, contentType, cacheControl string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	upInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	if contentType != "" {
		upInput.ContentType = aws.String(contentType)
	}

	if cacheControl != "" {
		upInput.CacheControl = aws.String(cacheControl)
	} else if c.config.DefaultCacheControl != "" {
		upInput.CacheControl = aws.String(c.config.DefaultCacheControl)
	}

	if c.config.ServerSideEncryption != "" {
		upInput.ServerSideEncryption = types.ServerSideEncryption(c.config.ServerSideEncryption)
	}

	_, err := c.uploader.Upload(ctx, upInput)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

// Get retrieves an object from the specified S3 bucket by its key and returns a read closer for its content.
func (c *Client) Get(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	out, err := c.raw.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return out.Body, nil
}

// Head retrieves metadata of an object from the specified S3 bucket and key without downloading its content.
func (c *Client) Head(ctx context.Context, bucket, key string) (*s3.HeadObjectOutput, error) {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	head, err := c.raw.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get head object: %w", err)
	}

	return head, nil
}

// Delete removes an object identified by the given bucket and key from the S3 storage. Returns an error if the operation fails.
func (c *Client) Delete(ctx context.Context, bucket, key string) error {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	_, err := c.raw.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// PresignGet generates a presigned URL for downloading an object from S3 with a specified expiration duration.
func (c *Client) PresignGet(ctx context.Context, bucket, key string, expires time.Duration) (string, http.Header, error) {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	if expires <= 0 {
		expires = 15 * time.Minute
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	req, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(po *s3.PresignOptions) {
		po.Expires = expires
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return req.URL, req.SignedHeader, nil
}

// PresignPut generates a presigned URL for uploading an object to the specified S3 bucket with a given key and content type.
// It also generates the signed headers required for the PUT request. The URL expires after the specified duration.
func (c *Client) PresignPut(ctx context.Context, bucket, key, contentType string, expires time.Duration) (string, http.Header, error) {
	if bucket == "" {
		bucket = c.config.Bucket
	}

	if expires <= 0 {
		expires = 15 * time.Minute
	}

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	in := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if contentType != "" {
		in.ContentType = aws.String(contentType)
	}

	if c.config.ServerSideEncryption != "" {
		in.ServerSideEncryption = types.ServerSideEncryption(c.config.ServerSideEncryption)
	}

	req, err := c.presign.PresignPutObject(ctx, in, func(po *s3.PresignOptions) {
		po.Expires = expires
	})

	if err != nil {
		return "", nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return req.URL, req.SignedHeader, nil
}

// Ping checks if the S3 storage is reachable and returns an error if it is not.
func (c *Client) Ping(ctx context.Context) error {
	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	_, err := c.raw.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(c.config.Bucket),
	})

	return err
}

// Raw returns the underlying AWS S3 client used for raw operations.
func (c *Client) Raw() *s3.Client {
	return c.raw
}

// withTimeout returns a context with a timeout derived from the client's configuration or the original context if the timeout is not set.
func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.config.Timeout <= 0 {
		// Retourne un cancel no-op pour simplifier les appels (toujours defer cancel()).
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, c.config.Timeout)
}
