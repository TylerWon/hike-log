package photostorage_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

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
	_, err := photostorage.New(nil, "", "bucket", "us-west-2", "http://localhost:4566")
	suite.Require().EqualError(err, "env is empty")
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenBucketNameIsEmpty() {
	_, err := photostorage.New(nil, "local", "", "us-west-2", "http://localhost:4566")
	suite.Require().EqualError(err, "bucketName is empty")
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenRegionIsEmpty() {
	_, err := photostorage.New(nil, "local", "bucket", "", "http://localhost:4566")
	suite.Require().EqualError(err, "region is empty")
}

func (suite *photoStorageTestSuite) TestNew_ReturnsErrorWhenExternalEndpointURLIsEmptyInLocalEnvironment() {
	_, err := photostorage.New(nil, "local", "bucket", "us-west-2", "")
	suite.Require().EqualError(err, "externalEndpointURL is empty")
}

func (suite *photoStorageTestSuite) TestNew_AllowsEmptyExternalEndpointURLInNonLocalEnvironment() {
	storage, err := photostorage.New(nil, "production", "bucket", "us-west-2", "")
	suite.Require().NoError(err)
	suite.Require().NotNil(storage)
}

func (suite *photoStorageTestSuite) TestCreatePhotoObjectKey_ReturnsKey() {
	hikeID := uint(1)
	storage := testutils.NewPhotoStorage(suite.T(), nil, nil, nil, nil, nil)
	objectKey := storage.CreatePhotoObjectKey(hikeID)
	suite.Require().Regexp(fmt.Sprintf(`^hikes/%d/photos/[0-9a-f-]{36}$`, hikeID), objectKey)
}

func (suite *photoStorageTestSuite) TestCreatePhotoObjectKey_ReturnsUniqueKeys() {
	storage := testutils.NewPhotoStorage(suite.T(), nil, nil, nil, nil, nil)
	key1 := storage.CreatePhotoObjectKey(1)
	key2 := storage.CreatePhotoObjectKey(1)
	suite.Require().NotEqual(key1, key2)
}

func (suite *photoStorageTestSuite) TestCreatePresignedUploadURL_ReturnsErrorWhenS3Errors() {
	mockS3Client := &s3.MockS3Client{
		CreatePresignedPutObjectRequestError: errors.New("Something went wrong"),
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	_, _, err := storage.CreatePresignedUploadURL(context.TODO(), 1, "image/jpeg", 100)

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestCreatePresignedUploadURL_ReturnsURLAndObjectKey() {
	mockS3Client := &s3.MockS3Client{
		CreatePresignedPutObjectRequestResult: &v4.PresignedHTTPRequest{URL: "https://presigned.example/upload"},
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	presignedURL, objectKey, err := storage.CreatePresignedUploadURL(context.TODO(), 7, "image/jpeg", 100)

	suite.Require().NoError(err)
	suite.Require().Equal("https://presigned.example/upload", presignedURL)
	suite.Require().NotEmpty(objectKey)
}

func (suite *photoStorageTestSuite) TestDeletePhoto_ReturnsErrorWhenS3Errors() {
	mockS3Client := &s3.MockS3Client{DeleteObjectError: errors.New("Something went wrong")}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	err := storage.DeletePhoto(context.TODO(), storage.GetPhotoObjectURL("abc"))

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestDeletePhoto_DeletesPhoto() {
	mockS3Client := &s3.MockS3Client{DeleteObjectResult: &s3Sdk.DeleteObjectOutput{}}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	err := storage.DeletePhoto(context.TODO(), storage.GetPhotoObjectURL("abc"))

	suite.Require().NoError(err)
}

func (suite *photoStorageTestSuite) TestDeletePhotosForHike_ReturnsErrorWhenS3ListErrors() {
	mockS3Client := &s3.MockS3Client{ListObjectsError: errors.New("Something went wrong")}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	err := storage.DeletePhotosForHike(context.TODO(), 1)

	suite.Require().Error(err)
}

func (suite *photoStorageTestSuite) TestDeletePhotosForHike_NoOpWhenNoPhotos() {
	mockS3Client := &s3.MockS3Client{
		ListObjectsResult: &s3Sdk.ListObjectsV2Output{Contents: []s3types.Object{}},
	}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

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
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

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
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	err := storage.DeletePhotosForHike(context.TODO(), 1)

	suite.Require().NoError(err)
}

func (suite *photoStorageTestSuite) TestDoesPhotoExist_ReturnsS3Result() {
	mockS3Client := &s3.MockS3Client{DoesObjectExistResult: true}
	storage := testutils.NewPhotoStorage(suite.T(), mockS3Client, nil, nil, nil, nil)

	exists, err := storage.DoesPhotoExist(context.TODO(), storage.CreatePhotoObjectKey(1))

	suite.Require().NoError(err)
	suite.Require().True(exists)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectKey_ReturnsObjectKeyFromPathStyleURLInLocalEnvironment() {
	env, bucketName, externalEndpointURL := "local", "bucket", "http://localhost:4566"
	storage := testutils.NewPhotoStorage(suite.T(), nil, &env, &bucketName, nil, &externalEndpointURL)
	objectKey := storage.CreatePhotoObjectKey(1)
	objectURL := fmt.Sprintf("%s/%s/%s", externalEndpointURL, bucketName, objectKey)

	actual := storage.GetPhotoObjectKey(objectURL)

	suite.Require().Equal(objectKey, actual)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectKey_ReturnsObjectKeyFromVirtualHostStyleURLInNonLocalEnvironment() {
	env, bucketName, region := "production", "bucket", "us-west-2"
	storage := testutils.NewPhotoStorage(suite.T(), nil, &env, &bucketName, &region, nil)
	objectKey := storage.CreatePhotoObjectKey(1)
	objectURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)

	actual := storage.GetPhotoObjectKey(objectURL)

	suite.Require().Equal(objectKey, actual)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectURL_ReturnsPathStyleURLInLocalEnvironment() {
	env, bucketName, externalEndpointURL := "local", "bucket", "http://localhost:4566"
	storage := testutils.NewPhotoStorage(suite.T(), nil, &env, &bucketName, nil, &externalEndpointURL)
	objectKey := storage.CreatePhotoObjectKey(1)

	url := storage.GetPhotoObjectURL(objectKey)

	expected := fmt.Sprintf("%s/%s/%s", externalEndpointURL, bucketName, objectKey)
	suite.Require().Equal(expected, url)
}

func (suite *photoStorageTestSuite) TestGetPhotoObjectURL_ReturnsVirtualHostStyleURLInNonLocalEnvironment() {
	env, bucketName, region := "production", "bucket", "us-west-2"
	storage := testutils.NewPhotoStorage(suite.T(), nil, &env, &bucketName, &region, nil)
	objectKey := storage.CreatePhotoObjectKey(1)

	url := storage.GetPhotoObjectURL(objectKey)

	expected := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)
	suite.Require().Equal(expected, url)
}
