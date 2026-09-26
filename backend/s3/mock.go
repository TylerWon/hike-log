package s3

import (
	"context"
	"io"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3Sdk "github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockS3Client struct {
	CreatePhotoObjectKeyResult            string
	CreatePresignedPutObjectRequestResult *v4.PresignedHTTPRequest
	CreatePresignedPutObjectRequestError  error
	DeleteObjectResult                    *s3Sdk.DeleteObjectOutput
	DeleteObjectError                     error
	DeleteObjectsResult                   *s3Sdk.DeleteObjectsOutput
	DeleteObjectsError                    error
	DoesObjectExistResult                 bool
	DoesObjectExistError                  error
	ListObjectsResult                     *s3Sdk.ListObjectsV2Output
	ListObjectsError                      error
	PutObjectResult                       *s3Sdk.PutObjectOutput
	PutObjectError                        error
}

func (m *MockS3Client) CreatePhotoObjectKey(hikeId uint) string {
	return m.CreatePhotoObjectKeyResult
}

func (m *MockS3Client) CreatePresignedPutObjectRequest(
	ctx context.Context,
	objectKey string,
	contentType string,
	contentLength int64,
) (*v4.PresignedHTTPRequest, error) {
	return m.CreatePresignedPutObjectRequestResult, m.CreatePresignedPutObjectRequestError
}

func (m *MockS3Client) DeleteObject(ctx context.Context, objectKey string) (*s3Sdk.DeleteObjectOutput, error) {
	return m.DeleteObjectResult, m.DeleteObjectError
}

func (m *MockS3Client) DeleteObjects(ctx context.Context, objectKeys []string) (*s3Sdk.DeleteObjectsOutput, error) {
	return m.DeleteObjectsResult, m.DeleteObjectsError
}

func (m *MockS3Client) DoesObjectExist(ctx context.Context, objectKey string) (bool, error) {
	return m.DoesObjectExistResult, m.DoesObjectExistError
}

func (m *MockS3Client) ListObjects(ctx context.Context, prefix string) (*s3Sdk.ListObjectsV2Output, error) {
	return m.ListObjectsResult, m.ListObjectsError
}

func (m *MockS3Client) PutObject(
	ctx context.Context,
	objectKey string,
	body io.Reader,
	contentType string,
) (*s3Sdk.PutObjectOutput, error) {
	return m.PutObjectResult, m.PutObjectError
}

type MockClient struct {
	DeleteObjectResult  *s3Sdk.DeleteObjectOutput
	DeleteObjectError   error
	DeleteObjectsResult *s3Sdk.DeleteObjectsOutput
	DeleteObjectsError  error
	HeadObjectResult    *s3Sdk.HeadObjectOutput
	HeadObjectError     error
	ListObjectsResult   *s3Sdk.ListObjectsV2Output
	ListObjectsError    error
	PutObjectResult     *s3Sdk.PutObjectOutput
	PutObjectError      error
}

func (m *MockClient) DeleteObject(
	ctx context.Context,
	params *s3Sdk.DeleteObjectInput,
	optFns ...func(*s3Sdk.Options),
) (*s3Sdk.DeleteObjectOutput, error) {
	return m.DeleteObjectResult, m.DeleteObjectError
}

func (m *MockClient) DeleteObjects(
	ctx context.Context,
	params *s3Sdk.DeleteObjectsInput,
	optFns ...func(*s3Sdk.Options),
) (*s3Sdk.DeleteObjectsOutput, error) {
	return m.DeleteObjectsResult, m.DeleteObjectsError
}

func (m *MockClient) HeadObject(
	ctx context.Context,
	params *s3Sdk.HeadObjectInput,
	optFns ...func(*s3Sdk.Options),
) (*s3Sdk.HeadObjectOutput, error) {
	return m.HeadObjectResult, m.HeadObjectError
}

func (m *MockClient) ListObjectsV2(
	ctx context.Context,
	params *s3Sdk.ListObjectsV2Input,
	optFns ...func(*s3Sdk.Options),
) (*s3Sdk.ListObjectsV2Output, error) {
	return m.ListObjectsResult, m.ListObjectsError
}

func (m *MockClient) PutObject(
	ctx context.Context,
	params *s3Sdk.PutObjectInput,
	optFns ...func(*s3Sdk.Options),
) (*s3Sdk.PutObjectOutput, error) {
	return m.PutObjectResult, m.PutObjectError
}

type MockPresignClient struct {
	PresignPutObjectResult *v4.PresignedHTTPRequest
	PresignPutObjectError  error
}

func (m *MockPresignClient) PresignPutObject(
	ctx context.Context,
	params *s3Sdk.PutObjectInput,
	optFns ...func(*s3Sdk.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	return m.PresignPutObjectResult, m.PresignPutObjectError
}
