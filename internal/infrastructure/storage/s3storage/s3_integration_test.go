package s3storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_NewS3(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	container, err := CreateS3Container(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Success Case", func(t *testing.T) {
		t.Parallel()
		_, err := NewS3(ctx, container.Config)
		assert.NoError(t, err)
	})

}

func TestIntegrationClient_PresignGet(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		bucket        string
		key           string
		expires       time.Duration
		expectedError error
	}{
		{
			name:          "Success Case",
			bucket:        "default",
			key:           "test.txt",
			expires:       time.Hour,
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreateS3Container(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s3, err := NewS3(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			url, header, err := s3.PresignGet(ctx, testCase.bucket, testCase.key, testCase.expires)

			if testCase.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
				assert.NotEmpty(t, header)
			}
		})
	}
}

func TestIntegrationClient_PresignPut(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		bucket        string
		key           string
		contentType   string
		expires       time.Duration
		expectedError error
	}{
		{
			name:          "Success Case",
			bucket:        "default",
			key:           "test.txt",
			contentType:   "text/plain",
			expires:       time.Hour,
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreateS3Container(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s3, err := NewS3(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			url, header, err := s3.PresignPut(ctx, testCase.bucket, testCase.key, testCase.contentType, testCase.expires)

			if testCase.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
				assert.NotEmpty(t, header)
			}
		})
	}
}
