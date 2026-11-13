package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/s3client"
)

// Creates a S3Client to use during testing. Arguments are optional and can be provided to override default values.
func NewS3Client(t *testing.T, env string, awsConfig config.AWSConfig) s3client.S3Client {
	t.Helper()

	testConfig := NewTestConfig()

	if env != "" {
		testConfig.Env = env
	}
	if awsConfig.AccessKeyID != "" {
		testConfig.AWS.AccessKeyID = awsConfig.AccessKeyID
	}
	if awsConfig.SecretAccessKey != "" {
		testConfig.AWS.SecretAccessKey = awsConfig.SecretAccessKey
	}
	if awsConfig.Region != "" {
		testConfig.AWS.Region = awsConfig.Region
	}
	if awsConfig.EndpointURL != "" {
		testConfig.AWS.EndpointURL = awsConfig.EndpointURL
	}
	if awsConfig.ExternalEndpointURL != "" {
		testConfig.AWS.ExternalEndpointURL = awsConfig.ExternalEndpointURL
	}
	if awsConfig.S3.BucketName != "" {
		testConfig.AWS.S3.BucketName = awsConfig.S3.BucketName
	}

	s3Client, err := s3client.New(testConfig.Env, testConfig.AWS)
	if err != nil {
		t.Fatal("Failed to setup S3 client: ", err)
	}

	return s3Client
}
