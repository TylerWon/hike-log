package aws

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type s3ClientTestSuite struct {
	suite.Suite
}

func (suite *s3ClientTestSuite) TestCreatePresignedPutObjectRequest_ReturnsErrorWhenS3Errors() {
	mockPresignClient := MockPresignClient{
		PresignPutObjectResult: nil,
		PresignPutObjectError:  errors.New("Something went wrong"),
	}

	s3Client := NewTestS3Client(&MockClient{}, &mockPresignClient, "hike-log")

	_, err := s3Client.CreatePresignedPutObjectRequest(context.TODO(), "abc", "application/json", 5000)
	suite.Error(err)
}

func (suite *s3ClientTestSuite) TestCreatePresignedPutObjectRequest_ReturnsPresignedRequest() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	request, err := s3Client.CreatePresignedPutObjectRequest(context.TODO(), "abc", "application/json", 5000)
	suite.NoError(err)

	suite.Equal("PUT", request.Method)

	suite.Equal("application/json", request.SignedHeader.Get("Content-Type"))
	suite.Equal("5000", request.SignedHeader.Get("Content-Length"))

	// The presigned URL has dynamic values so only verify certain parts of it
	parsedURL, err := url.Parse(request.URL)
	suite.NoError(err)
	suite.Contains(parsedURL.Path, "abc")

	parsedQueryParams, err := url.ParseQuery(parsedURL.RawQuery)
	suite.NoError(err)
	suite.Contains(parsedQueryParams, "X-Amz-SignedHeaders")
	suite.Contains(parsedQueryParams, "X-Amz-Expires")

	signedHeadersParam := parsedQueryParams.Get("X-Amz-SignedHeaders")
	suite.Contains(signedHeadersParam, "content-type")
	suite.Contains(signedHeadersParam, "content-length")

	expiresParam := parsedQueryParams.Get("X-Amz-Expires")
	suite.Contains(expiresParam, "900")
}

func (suite *s3ClientTestSuite) TestDeleteObject_ReturnsErrorWhenS3Errors() {
	mockClient := MockClient{
		DeleteObjectResult: nil,
		DeleteObjectError:  errors.New("Something went wrong"),
	}

	s3Client := NewTestS3Client(&mockClient, &MockPresignClient{}, "hike-log")

	_, err := s3Client.DeleteObject(context.TODO(), "abc")
	suite.Error(err)
}

func (suite *s3ClientTestSuite) TestDeleteObject_DeletesObject() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	objectKey := fmt.Sprintf("delete-object-%d", time.Now().UnixNano())
	_, err = s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.NoError(err)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.False(exists)
}

func (suite *s3ClientTestSuite) TestDeleteObject_NoErrorWhenObjectDoesNotExist() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	_, err = s3Client.DeleteObject(context.TODO(), "test")
	suite.NoError(err)
}

func (suite *s3ClientTestSuite) TestDoesObjectExist_ReturnsErrorWhenS3Errors() {
	mockClient := MockClient{
		HeadObjectResult: nil,
		HeadObjectError:  errors.New("Something went wrong"),
	}

	s3Client := NewTestS3Client(&mockClient, &MockPresignClient{}, "hike-log")

	_, err := s3Client.DoesObjectExist(context.TODO(), "abc")
	suite.Error(err)
}

func (suite *s3ClientTestSuite) TestDoesObjectExist_ReturnsFalseWhenObjectDoesNotExist() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), "abc")
	suite.NoError(err)
	suite.False(exists)
}

func (suite *s3ClientTestSuite) TestDoesObjectExist_ReturnsTrueWhenObjectExists() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	objectKey := fmt.Sprintf("does-object-exist-%d", time.Now().UnixNano())
	_, err = s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.NoError(err)
}

func (suite *s3ClientTestSuite) TestGetObjectURL_ReturnsPathStyleURLInLocalEnvironment() {
	awsEndpointURL := "http://localstack:4566"
	bucketName := "hike-log"
	objectKey := "test"

	suite.T().Setenv("AWS_ENDPOINT_URL", awsEndpointURL)

	s3Client := NewTestS3Client(&MockClient{}, &MockPresignClient{}, bucketName)
	url := s3Client.GetObjectURL(objectKey, "local")
	suite.Equal(fmt.Sprintf("%s/%s/%s", awsEndpointURL, bucketName, objectKey), url)
}

func (suite *s3ClientTestSuite) TestGetObjectURL_ReturnsVirtualHostStyleURLInNonLocalEnvironment() {
	awsRegion := "us-east-1"
	bucketName := "hike-log"
	objectKey := "test"

	suite.T().Setenv("AWS_REGION", awsRegion)

	s3Client := NewTestS3Client(&MockClient{}, &MockPresignClient{}, bucketName)
	url := s3Client.GetObjectURL(objectKey, "production")
	suite.Equal(fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, objectKey), url)
}

func (suite *s3ClientTestSuite) TestPutObject_ReturnsErrorWhenS3Errors() {
	mockClient := MockClient{
		PutObjectResult: nil,
		PutObjectError:  errors.New("Something went wrong"),
	}

	s3Client := NewTestS3Client(&mockClient, &MockPresignClient{}, "hike-log")

	_, err := s3Client.PutObject(context.TODO(), "abc", strings.NewReader("content"), "text/plain")
	suite.Error(err)
}

func (suite *s3ClientTestSuite) TestPutObject_InsertsNewObject() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	objectKey := fmt.Sprintf("put-object-%d", time.Now().UnixNano())
	_, err = s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.NoError(err)
}

func (suite *s3ClientTestSuite) TestPutObject_UpdatesExistingObject() {
	s3Client, err := NewS3Client()
	suite.NoError(err)

	objectKey := fmt.Sprintf("put-object-%d", time.Now().UnixNano())
	_, err = s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.NoError(err)

	exists, err := s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	_, err = s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "text/plain")
	suite.NoError(err)

	exists, err = s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	_, err = s3Client.DeleteObject(context.TODO(), objectKey)
	suite.NoError(err)
}

func TestS3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(s3ClientTestSuite))
}
