package s3storage

import (
	"context"
	"io"
	"os"
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

func TestIntegrationClient_Put(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		testFile      string
		bucket        string
		key           string
		contentType   string
		cacheControl  string
		expectedError error
	}{
		{
			name:          "Success Case",
			testFile:      "testdata.txt",
			bucket:        "default",
			key:           "test/testdata.txt",
			contentType:   "text/plain",
			cacheControl:  "max-age=86400",
			expectedError: nil,
		},
		{
			name:          "Success Case - Replace the file",
			testFile:      "testdata2.txt",
			bucket:        "default",
			key:           "test/testdata.txt",
			contentType:   "text/plain",
			cacheControl:  "max-age=86400",
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
			file, err := os.Open(testCase.testFile)
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				t.Fatal(err)
			}
			fileSize := fileInfo.Size()

			err = s3.Put(ctx, testCase.bucket, testCase.key, file, testCase.contentType, testCase.cacheControl)

			if testCase.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// check if the file is uploaded successfully
				head, err := s3.Head(ctx, testCase.bucket, testCase.key)
				assert.NoError(t, err)
				assert.Equal(t, testCase.contentType, *head.ContentType)
				assert.Equal(t, testCase.cacheControl, *head.CacheControl)
				assert.Equal(t, fileSize, *head.ContentLength)
			}
		})
	}
}

func TestIntegrationClient_Get(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		bucket        string
		key           string
		expectedError error
	}{
		{
			name:          "Success Case",
			bucket:        "default",
			key:           "test/testdata.txt",
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreateS3Container(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s3Client, err := NewS3(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}

	originalFile, err := os.Open("testdata.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer originalFile.Close()

	originalContent, err := io.ReadAll(originalFile)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			originalFile.Seek(0, 0)
			err := s3Client.Put(ctx, testCase.bucket, testCase.key, originalFile, "text/plain", "")
			if err != nil {
				t.Fatal(err)
			}

			reader, err := s3Client.Get(ctx, testCase.bucket, testCase.key)

			if testCase.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				defer reader.Close()

				content, err := io.ReadAll(reader)
				assert.NoError(t, err)

				assert.Equal(t, originalContent, content)
			}
		})
	}
}

func TestIntegrationClient_Delete(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		bucket        string
		key           string
		expectedError error
	}{
		{
			name:          "Success Case",
			bucket:        "default",
			key:           "test/testdata.txt",
			expectedError: nil,
		},
	}

	ctx := context.Background()
	container, err := CreateS3Container(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s3Client, err := NewS3(ctx, container.Config)
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open("testdata.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Upload du fichier avant de le supprimer
			err := s3Client.Put(ctx, testCase.bucket, testCase.key, file, "text/plain", "")
			if err != nil {
				t.Fatal(err)
			}

			// Test de la suppression
			err = s3Client.Delete(ctx, testCase.bucket, testCase.key)

			if testCase.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Vérifier que le fichier n'existe plus
				_, err := s3Client.Head(ctx, testCase.bucket, testCase.key)
				assert.Error(t, err)
			}
		})
	}
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
			t.Parallel()
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
			t.Parallel()
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
