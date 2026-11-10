package testutils

import (
	"os"
	"testing"

	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/s3"
)

// Creates a PhotoStorage to use during testing. Arguments are optional and can be provided to override default values.
func NewPhotoStorage(
	t *testing.T,
	s3Client s3.S3Client,
	env *string,
	bucketName *string,
	region *string,
	externalEndpointURL *string,
) photostorage.PhotoStorage {
	t.Helper()

	envValue := os.Getenv("ENV")
	bucketNameValue := os.Getenv("AWS_S3_BUCKET_NAME")
	regionValue := os.Getenv("AWS_REGION")
	externalEndpointURLValue := os.Getenv("AWS_EXTERNAL_ENDPOINT_URL")

	if env != nil {
		envValue = *env
	}
	if bucketName != nil {
		bucketNameValue = *bucketName
	}
	if region != nil {
		regionValue = *region
	}
	if externalEndpointURL != nil {
		externalEndpointURLValue = *externalEndpointURL
	}

	photoStorage, err := photostorage.New(
		s3Client,
		envValue,
		bucketNameValue,
		regionValue,
		externalEndpointURLValue,
	)
	if err != nil {
		t.Fatal("Failed to setup photo storage: ", err)
	}

	return photoStorage
}
