package testutils

import (
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockPresignClient struct {
	PresignPutObjectResult *v4.PresignedHTTPRequest
	PresignPutObjectError  error
}

func (m MockPresignClient) PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	return m.PresignPutObjectResult, m.PresignPutObjectError
}
