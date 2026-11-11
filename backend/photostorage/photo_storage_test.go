package photostorage_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3Sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/suite"
)

type photoStorageTestSuite struct {
	suite.Suite
}

func TestPhotoStorageTestSuite(t *testing.T) {
	suite.Run(t, new(photoStorageTestSuite))
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenEnvIsEmpty() {
	testConfig := testutils.TestConfig()
	_, err := photostorage.New(nil, "", testConfig.AWS)
	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenBucketNameIsEmpty() {
	testConfig := testutils.TestConfig()
	testConfig.AWS.S3.BucketName = ""
	_, err := photostorage.New(nil, testConfig.Env, testConfig.AWS)
	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenRegionIsEmpty() {
	testConfig := testutils.TestConfig()
	testConfig.AWS.Region = ""
	_, err := photostorage.New(nil, testConfig.Env, testConfig.AWS)
	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenExternalEndpointURLIsEmptyInLocalEnvironment() {
	testConfig := testutils.TestConfig()
	testConfig.Env = "local"
	testConfig.AWS.ExternalEndpointURL = ""
	_, err := photostorage.New(nil, testConfig.Env, testConfig.AWS)
	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestNew_AllowsEmptyExternalEndpointURLInNonLocalEnvironment() {
	testConfig := testutils.TestConfig()
	testConfig.Env = "production"
	testConfig.AWS.ExternalEndpointURL = ""
	storage, err := photostorage.New(nil, testConfig.Env, testConfig.AWS)
	suite.Require().NoError(err)
	suite.Require().NotNil(storage)
}

func (suite *photoStorageTestSuite) TestCreatePhotoObjectKey_ReturnsKey() {
	storage := testutils.NewPhotoStorage(suite.T(), nil, "", config.AWSConfig{})
	hikeID := uint(1)
	objectKey := storage.CreatePhotoObjectKey(hikeID)
	suite.Require().Regexp(fmt.Sprintf(`^hikes/%d/photos/[0-9a-f-]{36}$`, hikeID), objectKey)
}

func (suite *photoStorageTestSuite) TestCreatePhotoObjectKey_ReturnsUniqueKeys() {
	storage := testutils.NewPhotoStorage(suite.T(), nil, "", config.AWSConfig{})
	key1 := storage.CreatePhotoObjectKey(1)
	key2 := storage.CreatePhotoObjectKey(1)
	suite.Require().NotEqual(key1, key2)
}

func (suite *photoStorageTestSuite) TestCreatePresignedUploadURL_ReturnsErrorWhenS3Errors() {
	mockS3Client := &s3.MockS3Client{
		CreatePresignedPutObjectRequestError: errors.New("Something went wrong"),
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	_, _, err := storage.CreatePresignedUploadURL(context.TODO(), 1, "image/jpeg", 100)

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestCreatePresignedUploadURL_ReturnsURLAndObjectKey() {
	mockS3Client := &s3.MockS3Client{
		CreatePresignedPutObjectRequestResult: &v4.PresignedHTTPRequest{URL: "https://presigned.example/upload"},
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	presignedURL, objectKey, err := storage.CreatePresignedUploadURL(context.TODO(), 7, "image/jpeg", 100)

	suite.Require().NoError(err)
	suite.Require().Equal("https://presigned.example/upload", presignedURL)
	suite.Require().NotEmpty(objectKey)
}

func (suite *photoStorageTestSuite) TestDeletePhoto_ReturnsErrorWhenS3Errors() {
	mockS3Client := &s3.MockS3Client{DeleteObjectError: errors.New("Something went wrong")}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	err := storage.DeletePhoto(context.TODO(), storage.GetPhotoObjectURL("abc"))

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestDeletePhoto_DeletesPhoto() {
	mockS3Client := &s3.MockS3Client{DeleteObjectResult: &s3Sdk.DeleteObjectOutput{}}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	err := storage.DeletePhoto(context.TODO(), storage.GetPhotoObjectURL("abc"))

	suite.Require().NoError(err)
}

func (suite *photoStorageTestSuite) TestDeletePhotosForHike_ReturnsErrorWhenS3ListErrors() {
	mockS3Client := &s3.MockS3Client{ListObjectsError: errors.New("Something went wrong")}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	err := storage.DeletePhotosForHike(context.TODO(), 1)

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestDeletePhotosForHike_NoOpWhenNoPhotos() {
	mockS3Client := &s3.MockS3Client{
		ListObjectsResult: &s3Sdk.ListObjectsV2Output{Contents: []s3types.Object{}},
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	err := storage.DeletePhotosForHike(context.TODO(), 1)

	suite.Require().NoError(err)
}

func (suite *photoStorageTestSuite) TestDeletePhotosForHike_ReturnsErrorWhenS3DeleteErrors() {
	mockS3Client := &s3.MockS3Client{
		ListObjectsResult: &s3Sdk.ListObjectsV2Output{
			Contents: []s3types.Object{{Key: aws.String("hikes/1/photos/abc")}},
		},
		DeleteObjectsError: errors.New("Something went wrong"),
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	err := storage.DeletePhotosForHike(context.TODO(), 1)

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestDeletePhotosForHike_DeletesPhotos() {
	mockS3Client := &s3.MockS3Client{
		ListObjectsResult: &s3Sdk.ListObjectsV2Output{
			Contents: []s3types.Object{
				{Key: aws.String("hikes/1/photos/abc")},
				{Key: aws.String("hikes/1/photos/def")},
			},
		},
		DeleteObjectsResult: &s3Sdk.DeleteObjectsOutput{},
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	err := storage.DeletePhotosForHike(context.TODO(), 1)

	suite.Require().NoError(err)
}

func (suite *photoStorageTestSuite) TestDoesPhotoExist_ReturnsS3Result() {
	mockS3Client := &s3.MockS3Client{DoesObjectExistResult: true}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, "", config.AWSConfig{})

	exists, err := storage.DoesPhotoExist(context.TODO(), storage.CreatePhotoObjectKey(1))

	suite.Require().NoError(err)
	suite.Require().True(exists)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectKey_ReturnsObjectKeyFromPathStyleURLWhenInLocalEnvironment() {
	testConfig := testutils.TestConfig()
	testConfig.Env = "local"
	testConfig.AWS.S3.BucketName = "bucket"
	testConfig.AWS.ExternalEndpointURL = "http://localhost:4566"

	storage := testutils.NewPhotoStorage(suite.T(), nil, testConfig.Env, testConfig.AWS)

	objectKey := storage.CreatePhotoObjectKey(1)
	objectURL := fmt.Sprintf("%s/%s/%s", testConfig.AWS.ExternalEndpointURL, testConfig.AWS.S3.BucketName, objectKey)
	actual := storage.GetPhotoObjectKey(objectURL)

	suite.Require().Equal(objectKey, actual)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectKey_ReturnsObjectKeyFromVirtualHostStyleURLWhenInNonLocalEnvironment() {
	testConfig := testutils.TestConfig()
	testConfig.Env = "production"
	testConfig.AWS.S3.BucketName = "bucket"
	testConfig.AWS.Region = "us-east-1"

	storage := testutils.NewPhotoStorage(suite.T(), nil, testConfig.Env, testConfig.AWS)

	objectKey := storage.CreatePhotoObjectKey(1)
	objectURL := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		testConfig.AWS.S3.BucketName,
		testConfig.AWS.Region,
		objectKey,
	)
	actual := storage.GetPhotoObjectKey(objectURL)

	suite.Require().Equal(objectKey, actual)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectURL_ReturnsPathStyleURLInLocalEnvironment() {
	testConfig := testutils.TestConfig()
	testConfig.Env = "local"
	testConfig.AWS.S3.BucketName = "bucket"
	testConfig.AWS.ExternalEndpointURL = "http://localhost:4566"

	storage := testutils.NewPhotoStorage(suite.T(), nil, testConfig.Env, testConfig.AWS)

	objectKey := storage.CreatePhotoObjectKey(1)
	url := storage.GetPhotoObjectURL(objectKey)

	expected := fmt.Sprintf("%s/%s/%s", testConfig.AWS.ExternalEndpointURL, testConfig.AWS.S3.BucketName, objectKey)
	suite.Require().Equal(expected, url)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectURL_ReturnsVirtualHostStyleURLInNonLocalEnvironment() {
	testConfig := testutils.TestConfig()
	testConfig.Env = "production"
	testConfig.AWS.S3.BucketName = "bucket"
	testConfig.AWS.Region = "us-east-1"

	storage := testutils.NewPhotoStorage(suite.T(), nil, testConfig.Env, testConfig.AWS)

	objectKey := storage.CreatePhotoObjectKey(1)
	url := storage.GetPhotoObjectURL(objectKey)

	expected := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		testConfig.AWS.S3.BucketName,
		testConfig.AWS.Region,
		objectKey,
	)
	suite.Require().Equal(expected, url)
}
