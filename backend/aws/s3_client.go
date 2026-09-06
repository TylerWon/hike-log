package aws

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client handles interactions with the AWS S3 bucket for the app.
type S3Client interface {
	CreatePresignedPutObjectRequest(ctx context.Context, objectKey string, contentType string, contentLength int64) (*v4.PresignedHTTPRequest, error)
}

// s3ClientImpl is an implementation of the S3Client interface.
type s3ClientImpl struct {
	presignClient presignClient
	bucketName    string
}

// presignClient is a thin wrapper around the s3.PresignClient type
type presignClient interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// Creates a new S3Client.
func NewS3Client() (S3Client, error) {
	// LoadDefaultConfig loads configuration from all the SDK's supported sources (env vars, ~/.aws/config,
	// ~/.aws/credentials) and resolves the credentials using the SDK's default credential chain
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)

	presignClient := s3.NewPresignClient(client)

	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucketName == "" {
		return nil, errors.New("Bucket name could not be retrieved from AWS_S3_BUCKET_NAME env var")
	}

	return &s3ClientImpl{presignClient, bucketName}, nil
}

// Creates a new S3Client. Allows injection of internal dependencies to allow for mocking.
func NewTestS3Client(presignClient presignClient, bucketName string) S3Client {
	return &s3ClientImpl{presignClient, bucketName}
}

// Creates a presigned request that can be used to put an object in the bucket. The request expires after 900 seconds.
func (s3Client *s3ClientImpl) CreatePresignedPutObjectRequest(
	ctx context.Context,
	objectKey string,
	contentType string,
	contentLength int64,
) (*v4.PresignedHTTPRequest, error) {
	request, err := s3Client.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s3Client.bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	})

	if err != nil {
		return nil, err
	}

	return request, err
}
