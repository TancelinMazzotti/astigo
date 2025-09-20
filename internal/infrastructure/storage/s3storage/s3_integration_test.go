package s3storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration_NewS3(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		expectedError error
	}{
		{
			name: "Success Case",
		},
	}

	ctx := context.Background()
	container, err := CreateS3Container(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewS3(ctx, container.Config)
			if testCase.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
