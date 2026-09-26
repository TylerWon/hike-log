package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/s3"
)

// Creates a S3Client to use during testing.
func NewS3Client(t *testing.T) s3.S3Client {
	s3Client, err := s3.NewS3Client()
	if err != nil {
		t.Fatal("Failed to setup S3 client: ", err)
	}

	return s3Client
}
