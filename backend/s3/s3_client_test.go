package s3_test

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

type s3ClientTestSuite struct {
	suite.Suite
}

func TestS3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(s3ClientTestSuite))
}

func (suite *s3ClientTestSuite) TestCreatePresignedPutObjectRequest_ReturnsErrorWhenS3Errors() {
	mockPresignClient := s3.MockPresignClient{
		PresignPutObjectResult: nil,
		PresignPutObjectError:  errors.New("Something went wrong"),
	}

	s3Client := s3.NewTestS3Client(&s3.MockClient{}, &mockPresignClient, "hike-log")

	_, err := s3Client.CreatePresignedPutObjectRequest(context.TODO(), "abc", "application/json", 5000)
	suite.Require().Error(err)
}

func (suite *s3ClientTestSuite) TestCreatePresignedPutObjectRequest_ReturnsPresignedRequestWithExternalEndpointWhenLocalEnvironment() {
	awsConfig := config.AWSConfig{ExternalEndpointURL: "http://localhost:4566"}
	s3Client := testutils.NewS3Client(suite.T(), "local", awsConfig)

	objectKey := "abc"
	request, err := s3Client.CreatePresignedPutObjectRequest(context.TODO(), objectKey, "application/json", 5000)
	suite.Require().NoError(err)

	suite.Require().Equal("PUT", request.Method)

	suite.Require().Equal("application/json", request.SignedHeader.Get("Content-Type"))
	suite.Require().Equal("5000", request.SignedHeader.Get("Content-Length"))

	// The presigned URL has dynamic values so only verify certain parts of it
	parsedURL, err := url.Parse(request.URL)
	suite.Require().NoError(err)
	suite.Require().Equal("localhost:4566", parsedURL.Host)
	suite.Require().Contains(parsedURL.Path, objectKey)

	parsedQueryParams, err := url.ParseQuery(parsedURL.RawQuery)
	suite.Require().NoError(err)
	suite.Require().Contains(parsedQueryParams, "X-Amz-SignedHeaders")
	suite.Require().Contains(parsedQueryParams, "X-Amz-Expires")

	signedHeadersParam := parsedQueryParams.Get("X-Amz-SignedHeaders")
	suite.Require().Contains(signedHeadersParam, "content-type")
	suite.Require().Contains(signedHeadersParam, "content-length")

	expiresParam := parsedQueryParams.Get("X-Amz-Expires")
	suite.Require().Contains(expiresParam, "900")
}

func (suite *s3ClientTestSuite) TestCreatePresignedPutObjectRequest_ReturnsPresignedRequestWithDefaultAWSEndpointWhenNonLocalEnvironment() {
	awsConfig := config.AWSConfig{Region: "us-east-1", S3: config.S3Config{BucketName: "bucket"}}
	s3Client := testutils.NewS3Client(suite.T(), "production", awsConfig)

	objectKey := "abc"
	request, err := s3Client.CreatePresignedPutObjectRequest(context.TODO(), objectKey, "application/json", 5000)
	suite.Require().NoError(err)

	suite.Require().Equal("PUT", request.Method)

	suite.Require().Equal("application/json", request.SignedHeader.Get("Content-Type"))
	suite.Require().Equal("5000", request.SignedHeader.Get("Content-Length"))

	// The presigned URL has dynamic values so only verify certain parts of it
	parsedURL, err := url.Parse(request.URL)
	suite.Require().NoError(err)
	expectedHost := fmt.Sprintf("%s.s3.%s.amazonaws.com", awsConfig.S3.BucketName, awsConfig.Region)
	suite.Require().Equal(expectedHost, parsedURL.Host)
	suite.Require().Contains(parsedURL.Path, objectKey)

	parsedQueryParams, err := url.ParseQuery(parsedURL.RawQuery)
	suite.Require().NoError(err)
	suite.Require().Contains(parsedQueryParams, "X-Amz-SignedHeaders")
	suite.Require().Contains(parsedQueryParams, "X-Amz-Expires")

	signedHeadersParam := parsedQueryParams.Get("X-Amz-SignedHeaders")
	suite.Require().Contains(signedHeadersParam, "content-type")
	suite.Require().Contains(signedHeadersParam, "content-length")

	expiresParam := parsedQueryParams.Get("X-Amz-Expires")
	suite.Require().Contains(expiresParam, "900")
}

func (suite *s3ClientTestSuite) TestDeleteObject_ReturnsErrorWhenS3Errors() {
	mockClient := s3.MockClient{
		DeleteObjectResult: nil,
		DeleteObjectError:  errors.New("Something went wrong"),
	}

	s3Client := s3.NewTestS3Client(&mockClient, &s3.MockPresignClient{}, "hike-log")

	_, err := s3Client.DeleteObject(context.TODO(), "abc")
	suite.Require().Error(err)
}

func (suite *s3ClientTestSuite) TestDeleteObject_DeletesObject() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	objectKey := fmt.Sprintf("delete-object-%d", time.Now().UnixNano())
	_, err := s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.Require().NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().False(exists)
}

func (suite *s3ClientTestSuite) TestDeleteObject_NoErrorWhenObjectDoesNotExist() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	_, err := s3Client.DeleteObject(context.TODO(), "test")
	suite.Require().NoError(err)
}

func (suite *s3ClientTestSuite) TestDeleteObjects_ReturnsErrorWhenS3Errors() {
	mockClient := s3.MockClient{
		DeleteObjectsResult: nil,
		DeleteObjectsError:  errors.New("Something went wrong"),
	}

	s3Client := s3.NewTestS3Client(&mockClient, &s3.MockPresignClient{}, "hike-log")

	_, err := s3Client.DeleteObjects(context.TODO(), []string{"abc", "def"})
	suite.Require().Error(err)
}

func (suite *s3ClientTestSuite) TestDeleteObjects_DeletesObjects() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	objectKey1 := fmt.Sprintf("delete-objects-1-%d", time.Now().UnixNano())
	_, err := s3Client.PutObject(context.TODO(), objectKey1, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	objectKey2 := fmt.Sprintf("delete-objects-2-%d", time.Now().UnixNano())
	_, err = s3Client.PutObject(context.TODO(), objectKey2, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	_, err = s3Client.DeleteObjects(context.TODO(), []string{objectKey1, objectKey2})
	suite.Require().NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey1)
	suite.Require().NoError(err)
	suite.Require().False(exists)

	exists, err = s3Client.DoesObjectExist(context.TODO(), objectKey2)
	suite.Require().NoError(err)
	suite.Require().False(exists)
}

func (suite *s3ClientTestSuite) TestDeleteObjects_NoErrorWhenObjectDoesNotExist() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	objectKey := fmt.Sprintf("delete-objects-%d", time.Now().UnixNano())
	_, err := s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	_, err = s3Client.DeleteObjects(context.TODO(), []string{objectKey, "abc"})
	suite.Require().NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().False(exists)
}

func (suite *s3ClientTestSuite) TestDoesObjectExist_ReturnsErrorWhenS3Errors() {
	mockClient := s3.MockClient{
		HeadObjectResult: nil,
		HeadObjectError:  errors.New("Something went wrong"),
	}

	s3Client := s3.NewTestS3Client(&mockClient, &s3.MockPresignClient{}, "hike-log")

	_, err := s3Client.DoesObjectExist(context.TODO(), "abc")
	suite.Require().Error(err)
}

func (suite *s3ClientTestSuite) TestDoesObjectExist_ReturnsFalseWhenObjectDoesNotExist() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	exists, err := s3Client.DoesObjectExist(context.TODO(), "abc")
	suite.Require().NoError(err)
	suite.Require().False(exists)
}

func (suite *s3ClientTestSuite) TestDoesObjectExist_ReturnsTrueWhenObjectExists() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	objectKey := fmt.Sprintf("does-object-exist-%d", time.Now().UnixNano())
	_, err := s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().True(exists)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.Require().NoError(err)
}

func (suite *s3ClientTestSuite) TestPutObject_ReturnsErrorWhenS3Errors() {
	mockClient := s3.MockClient{
		PutObjectResult: nil,
		PutObjectError:  errors.New("Something went wrong"),
	}

	s3Client := s3.NewTestS3Client(&mockClient, &s3.MockPresignClient{}, "hike-log")

	_, err := s3Client.PutObject(context.TODO(), "abc", strings.NewReader("content"), "text/plain")
	suite.Require().Error(err)
}

func (suite *s3ClientTestSuite) TestPutObject_InsertsNewObject() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	objectKey := fmt.Sprintf("put-object-%d", time.Now().UnixNano())
	_, err := s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().True(exists)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.Require().NoError(err)
}

func (suite *s3ClientTestSuite) TestPutObject_UpdatesExistingObject() {
	s3Client := testutils.NewS3Client(suite.T(), "", config.AWSConfig{})

	objectKey := fmt.Sprintf("put-object-%d", time.Now().UnixNano())
	_, err := s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().True(exists)

	_, err = s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.Require().NoError(err)

	exists, err = s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().True(exists)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.Require().NoError(err)
}
