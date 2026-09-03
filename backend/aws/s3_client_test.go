package aws

import (
	"context"
	"errors"
	"net/url"
	"testing"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/stretchr/testify/suite"
)

type s3ClientTestSuite struct {
	suite.Suite
}

func (suite *s3ClientTestSuite) TestCreatePresignedPutObjectRequest_ReturnsErrorWhenRequestCannotBeCreated() {
	mockPresignClient := MockPresignClient{
		PresignPutObjectResult: &v4.PresignedHTTPRequest{},
		PresignPutObjectError:  errors.New("Simulated error"),
	}

	s3Client, err := NewTestS3Client(&mockPresignClient, "hike-log")
	suite.NoError(err)

	_, err = s3Client.CreatePresignedPutObjectRequest(context.TODO(), "abc", "application/json", 5000)
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

func TestS3ClientTestSuite(t *testing.T) {
	suite.Run(t, new(s3ClientTestSuite))
}
