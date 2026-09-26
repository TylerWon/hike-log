package s3

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	s3Sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Client handles interactions with the AWS S3 bucket for the app.
type S3Client interface {
	CreatePresignedPutObjectRequest(
		ctx context.Context,
		objectKey string,
		contentType string,
		contentLength int64,
	) (*v4.PresignedHTTPRequest, error)
	DeleteObject(ctx context.Context, objectKey string) (*s3Sdk.DeleteObjectOutput, error)
	DeleteObjects(ctx context.Context, objectKeys []string) (*s3Sdk.DeleteObjectsOutput, error)
	DoesObjectExist(ctx context.Context, objectKey string) (bool, error)
	ListObjects(ctx context.Context, prefix string) (*s3Sdk.ListObjectsV2Output, error)
	PutObject(ctx context.Context, objectKey string, body io.Reader, contentType string) (*s3Sdk.PutObjectOutput, error)
}

// s3ClientImpl is an implementation of the S3Client interface.
type s3ClientImpl struct {
	client        client
	presignClient presignClient
	bucketName    string
}

// presignClient is a thin wrapper around the s3Sdk.PresignClient type
type presignClient interface {
	PresignPutObject(
		ctx context.Context,
		params *s3Sdk.PutObjectInput,
		optFns ...func(*s3Sdk.PresignOptions),
	) (*v4.PresignedHTTPRequest, error)
}

// client is a thin wrapper around the s3Sdk.Client type.
type client interface {
	DeleteObject(
		ctx context.Context,
		params *s3Sdk.DeleteObjectInput,
		optFns ...func(*s3Sdk.Options),
	) (*s3Sdk.DeleteObjectOutput, error)
	DeleteObjects(
		ctx context.Context,
		params *s3Sdk.DeleteObjectsInput,
		optFns ...func(*s3Sdk.Options),
	) (*s3Sdk.DeleteObjectsOutput, error)
	HeadObject(
		ctx context.Context,
		params *s3Sdk.HeadObjectInput,
		optFns ...func(*s3Sdk.Options),
	) (*s3Sdk.HeadObjectOutput, error)
	ListObjectsV2(
		ctx context.Context,
		params *s3Sdk.ListObjectsV2Input,
		optFns ...func(*s3Sdk.Options),
	) (*s3Sdk.ListObjectsV2Output, error)
	PutObject(
		ctx context.Context,
		params *s3Sdk.PutObjectInput,
		optFns ...func(*s3Sdk.Options),
	) (*s3Sdk.PutObjectOutput, error)
}

// Creates a new S3Client.
func NewS3Client() (S3Client, error) {
	// Resolve AWS credentials using the SDK's default credential chain (i.e. (env vars, then ~/.aws/config, then
	// ~/.aws/credentials)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	var client *s3Sdk.Client
	var presignClient *s3Sdk.PresignClient
	if os.Getenv("ENV") == "local" {
		client = s3Sdk.NewFromConfig(cfg, func(o *s3Sdk.Options) {
			// Docker DNS cannot resolve default virtual-hosted-style addresses so use path-style instead
			o.UsePathStyle = true
		})

		publicEndpoint := os.Getenv("AWS_PUBLIC_ENDPOINT_URL")
		presignS3Client := s3Sdk.NewFromConfig(cfg, func(o *s3Sdk.Options) {
			// Docker DNS cannot resolve default virtual-hosted-style addresses so use path-style instead
			o.UsePathStyle = true

			// Ensures presign URLs use a publicly accessible endpoint (i.e. "localhost:4566" instead of "localstack:4566").
			// Only services running in Docker know how to resolve the "localstack" hostname. Presign URLs are used in the
			// browser (i.e.) outside of Docker so they cannot contain this hostname. As a result, "localhost:4566" must be
			// used instead as it is accessible by the host machine and maps to the Localstack Docker service via port-
			// forwarding.
			o.BaseEndpoint = aws.String(publicEndpoint)
		})
		presignClient = s3Sdk.NewPresignClient(presignS3Client)
	} else {
		client = s3Sdk.NewFromConfig(cfg)
		presignClient = s3Sdk.NewPresignClient(client)
	}

	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucketName == "" {
		return nil, errors.New("Bucket name could not be retrieved from AWS_S3_BUCKET_NAME env var")
	}

	return &s3ClientImpl{client, presignClient, bucketName}, nil
}

// Creates a new S3Client. Allows injection of internal dependencies to allow for mocking.
func NewTestS3Client(client client, presignClient presignClient, bucketName string) S3Client {
	return &s3ClientImpl{client, presignClient, bucketName}
}

// Creates a presigned request that can be used to put an object in the bucket. The request expires after 900 seconds.
func (s3Client *s3ClientImpl) CreatePresignedPutObjectRequest(
	ctx context.Context,
	objectKey string,
	contentType string,
	contentLength int64,
) (*v4.PresignedHTTPRequest, error) {
	output, err := s3Client.presignClient.PresignPutObject(ctx, &s3Sdk.PutObjectInput{
		Bucket:        aws.String(s3Client.bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	})

	return output, err
}

// Deletes an object from the bucket. Deleting a missing object succeeds.
func (s3Client *s3ClientImpl) DeleteObject(ctx context.Context, objectKey string) (*s3Sdk.DeleteObjectOutput, error) {
	output, err := s3Client.client.DeleteObject(ctx, &s3Sdk.DeleteObjectInput{
		Bucket: aws.String(s3Client.bucketName),
		Key:    aws.String(objectKey),
	})

	return output, err
}

// Deletes objects from the bucket. Deleting a missing object does not result in an error.
func (s3Client *s3ClientImpl) DeleteObjects(ctx context.Context, objectKeys []string) (*s3Sdk.DeleteObjectsOutput, error) {
	var objectsToDelete []types.ObjectIdentifier
	for _, objectKey := range objectKeys {
		objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{Key: aws.String(objectKey)})
	}

	output, err := s3Client.client.DeleteObjects(ctx, &s3Sdk.DeleteObjectsInput{
		Bucket: aws.String(s3Client.bucketName),
		Delete: &types.Delete{Objects: objectsToDelete},
	})

	return output, err
}

// Checks if an object with the given key exists in the bucket.
func (s3Client *s3ClientImpl) DoesObjectExist(ctx context.Context, objectKey string) (bool, error) {
	_, err := s3Client.client.HeadObject(ctx, &s3Sdk.HeadObjectInput{
		Bucket: aws.String(s3Client.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Lists all objects in the bucket which start with the given prefix.
func (s3Client *s3ClientImpl) ListObjects(ctx context.Context, prefix string) (*s3Sdk.ListObjectsV2Output, error) {
	output, err := s3Client.client.ListObjectsV2(ctx, &s3Sdk.ListObjectsV2Input{
		Bucket: aws.String(s3Client.bucketName),
		Prefix: aws.String(prefix),
	})

	return output, err
}

// Upserts an object into the bucket.
func (s3Client *s3ClientImpl) PutObject(
	ctx context.Context,
	objectKey string,
	body io.Reader,
	contentType string,
) (*s3Sdk.PutObjectOutput, error) {
	output, err := s3Client.client.PutObject(ctx, &s3Sdk.PutObjectInput{
		Body:        body,
		Bucket:      aws.String(s3Client.bucketName),
		ContentType: aws.String(contentType),
		Key:         aws.String(objectKey),
	})

	return output, err
}
