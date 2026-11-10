package testutils

import (
	"os"
	"testing"

	"github.com/TylerWon/hike-log/backend/s3"
)

// Creates a S3Client to use during testing. Arguments are optional and can be provided to override default values.
func NewS3Client(
	t *testing.T,
	env *string,
	accessKeyID *string,
	secretAccessKey *string,
	region *string,
	endpointURL *string,
	externalEndpointURL *string,
	bucketName *string,
) s3.S3Client {
	t.Helper()

	envValue := os.Getenv("ENV")
	accessKeyIDValue := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKeyValue := os.Getenv("AWS_SECRET_ACCESS_KEY")
	regionValue := os.Getenv("AWS_REGION")
	endpointURLValue := os.Getenv("AWS_ENDPOINT_URL")
	externalEndpointURLValue := os.Getenv("AWS_EXTERNAL_ENDPOINT_URL")
	bucketNameValue := os.Getenv("AWS_S3_BUCKET_NAME")

	if env != nil {
		envValue = *env
	}
	if accessKeyID != nil {
		accessKeyIDValue = *accessKeyID
	}
	if secretAccessKey != nil {
		secretAccessKeyValue = *secretAccessKey
	}
	if region != nil {
		regionValue = *region
	}
	if endpointURL != nil {
		endpointURLValue = *endpointURL
	}
	if externalEndpointURL != nil {
		externalEndpointURLValue = *externalEndpointURL
	}
	if bucketName != nil {
		bucketNameValue = *bucketName
	}

	s3Client, err := s3.NewS3Client(
		envValue,
		accessKeyIDValue,
		secretAccessKeyValue,
		regionValue,
		endpointURLValue,
		externalEndpointURLValue,
		bucketNameValue,
	)
	if err != nil {
		t.Fatal("Failed to setup S3 client: ", err)
	}

	return s3Client
}
