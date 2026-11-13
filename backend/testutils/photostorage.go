package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/s3client"
)

// Creates a PhotoStorage to use during testing. Arguments are optional and can be provided to override default values.
func NewPhotoStorage(
	t *testing.T,
	s3Client s3client.S3Client,
	env string,
	awsConfig config.AWSConfig,
) photostorage.PhotoStorage {
	t.Helper()

	testConfig := NewTestConfig()

	if env != "" {
		testConfig.Env = env
	}
	if awsConfig.S3.BucketName != "" {
		testConfig.AWS.S3.BucketName = awsConfig.S3.BucketName
	}
	if awsConfig.Region != "" {
		testConfig.AWS.Region = awsConfig.Region
	}
	if awsConfig.ExternalEndpointURL != "" {
		testConfig.AWS.ExternalEndpointURL = awsConfig.ExternalEndpointURL
	}

	photoStorage, err := photostorage.New(s3Client, testConfig.Env, testConfig.AWS)
	if err != nil {
		t.Fatal("Failed to setup photo storage: ", err)
	}

	return photoStorage
}
