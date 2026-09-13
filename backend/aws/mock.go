package aws

import (
	"context"
	"io"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockS3Client struct {
	CreatePresignedPutObjectRequestResult *v4.PresignedHTTPRequest
	CreatePresignedPutObjectRequestError  error
	DeleteObjectResult                    *s3.DeleteObjectOutput
	DeleteObjectError                     error
	DeleteObjectsResult                   *s3.DeleteObjectsOutput
	DeleteObjectsError                    error
	DoesObjectExistResult                 bool
	DoesObjectExistError                  error
	GetObjectKeyResult                    string
	GetObjectURLResult                    string
	ListObjectsResult                     *s3.ListObjectsV2Output
	ListObjectsError                      error
	PutObjectResult                       *s3.PutObjectOutput
	PutObjectError                        error
}

func (m *MockS3Client) CreatePresignedPutObjectRequest(
	ctx context.Context,
	objectKey string,
	contentType string,
	contentLength int64,
) (*v4.PresignedHTTPRequest, error) {
	return m.CreatePresignedPutObjectRequestResult, m.CreatePresignedPutObjectRequestError
}

func (m *MockS3Client) DeleteObject(ctx context.Context, objectKey string) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectResult, m.DeleteObjectError
}

func (m *MockS3Client) DeleteObjects(ctx context.Context, objectKeys []string) (*s3.DeleteObjectsOutput, error) {
	return m.DeleteObjectsResult, m.DeleteObjectsError
}

func (m *MockS3Client) DoesObjectExist(ctx context.Context, objectKey string) (bool, error) {
	return m.DoesObjectExistResult, m.DoesObjectExistError
}

func (m *MockS3Client) GetObjectKey(objectURL string, env string) string {
	return m.GetObjectKeyResult
}

func (m *MockS3Client) GetObjectURL(objectKey string, env string) string {
	return m.GetObjectURLResult
}

func (m *MockS3Client) ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error) {
	return m.ListObjectsResult, m.ListObjectsError
}

func (m *MockS3Client) PutObject(
	ctx context.Context,
	objectKey string,
	body io.Reader,
	contentType string,
) (*s3.PutObjectOutput, error) {
	return m.PutObjectResult, m.PutObjectError
}

type MockClient struct {
	DeleteObjectResult  *s3.DeleteObjectOutput
	DeleteObjectError   error
	DeleteObjectsResult *s3.DeleteObjectsOutput
	DeleteObjectsError  error
	HeadObjectResult    *s3.HeadObjectOutput
	HeadObjectError     error
	ListObjectsResult   *s3.ListObjectsV2Output
	ListObjectsError    error
	PutObjectResult     *s3.PutObjectOutput
	PutObjectError      error
}

func (m *MockClient) DeleteObject(
	ctx context.Context,
	params *s3.DeleteObjectInput,
	optFns ...func(*s3.Options),
) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectResult, m.DeleteObjectError
}

func (m *MockClient) DeleteObjects(
	ctx context.Context,
	params *s3.DeleteObjectsInput,
	optFns ...func(*s3.Options),
) (*s3.DeleteObjectsOutput, error) {
	return m.DeleteObjectsResult, m.DeleteObjectsError
}

func (m *MockClient) HeadObject(
	ctx context.Context,
	params *s3.HeadObjectInput,
	optFns ...func(*s3.Options),
) (*s3.HeadObjectOutput, error) {
	return m.HeadObjectResult, m.HeadObjectError
}

func (m *MockClient) ListObjectsV2(
	ctx context.Context,
	params *s3.ListObjectsV2Input,
	optFns ...func(*s3.Options),
) (*s3.ListObjectsV2Output, error) {
	return m.ListObjectsResult, m.ListObjectsError
}

func (m *MockClient) PutObject(
	ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.Options),
) (*s3.PutObjectOutput, error) {
	return m.PutObjectResult, m.PutObjectError
}

type MockPresignClient struct {
	PresignPutObjectResult *v4.PresignedHTTPRequest
	PresignPutObjectError  error
}

func (m *MockPresignClient) PresignPutObject(
	ctx context.Context,
	params *s3.PutObjectInput,
	optFns ...func(*s3.PresignOptions),
) (*v4.PresignedHTTPRequest, error) {
	return m.PresignPutObjectResult, m.PresignPutObjectError
}
