package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	DeleteObject(ctx context.Context, objectKey string) (*s3.DeleteObjectOutput, error)
	DeleteObjects(ctx context.Context, objectKeys []string) (*s3.DeleteObjectsOutput, error)
	DoesObjectExist(ctx context.Context, objectKey string) (bool, error)
	GetObjectKey(objectURL string, env string) string
	GetObjectURL(objectKey string, env string) string
	ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error)
	PutObject(ctx context.Context, objectKey string, body io.Reader, contentType string) (*s3.PutObjectOutput, error)
}

// s3ClientImpl is an implementation of the S3Client interface.
type s3ClientImpl struct {
	client        client
	presignClient presignClient
	bucketName    string
}

// presignClient is a thin wrapper around the s3.PresignClient type
type presignClient interface {
	PresignPutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.PresignOptions),
	) (*v4.PresignedHTTPRequest, error)
}

// client is a thin wrapper around the s3.Client type.
type client interface {
	DeleteObject(
		ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options),
	) (*s3.DeleteObjectOutput, error)
	DeleteObjects(
		ctx context.Context,
		params *s3.DeleteObjectsInput,
		optFns ...func(*s3.Options),
	) (*s3.DeleteObjectsOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	ListObjectsV2(
		ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options),
	) (*s3.ListObjectsV2Output, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// Creates a new S3Client.
func NewS3Client() (S3Client, error) {
	// LoadDefaultConfig loads configuration from all the SDK's supported sources (env vars, ~/.aws/config,
	// ~/.aws/credentials) and resolves the credentials using the SDK's default credential chain
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	var client *s3.Client
	if os.Getenv("ENV") == "local" {
		// Docker DNS cannot resolve default virtual-hosted-style addresses so use path-style instead
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	} else {
		client = s3.NewFromConfig(cfg)
	}

	presignClient := s3.NewPresignClient(client)

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
	output, err := s3Client.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s3Client.bucketName),
		Key:           aws.String(objectKey),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	})

	return output, err
}

// Deletes an object from the bucket. Deleting a missing object succeeds.
func (s3Client *s3ClientImpl) DeleteObject(ctx context.Context, objectKey string) (*s3.DeleteObjectOutput, error) {
	output, err := s3Client.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s3Client.bucketName),
		Key:    aws.String(objectKey),
	})

	return output, err
}

// Deletes objects from the bucket. Deleting a missing object does not result in an error.
func (s3Client *s3ClientImpl) DeleteObjects(ctx context.Context, objectKeys []string) (*s3.DeleteObjectsOutput, error) {
	var objectsToDelete []types.ObjectIdentifier
	for _, objectKey := range objectKeys {
		objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{Key: aws.String(objectKey)})
	}

	output, err := s3Client.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s3Client.bucketName),
		Delete: &types.Delete{Objects: objectsToDelete},
	})

	return output, err
}

// Checks if an object with the given key exists in the bucket.
func (s3Client *s3ClientImpl) DoesObjectExist(ctx context.Context, objectKey string) (bool, error) {
	_, err := s3Client.client.HeadObject(ctx, &s3.HeadObjectInput{
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

// Returns the object key from the given S3 object URL. Expects a path-style URL in local environment, otherwise expects
// a normal virtual-host-style object URL.
func (s3Client *s3ClientImpl) GetObjectKey(objectURL string, env string) string {
	if env == "local" {
		endpoint := os.Getenv("AWS_ENDPOINT_URL")
		prefix := fmt.Sprintf("%s/%s/", endpoint, s3Client.bucketName)
		return strings.TrimPrefix(objectURL, prefix)
	}

	region := os.Getenv("AWS_REGION")
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s3Client.bucketName, region)
	return strings.TrimPrefix(objectURL, prefix)
}

// Returns the S3 URL of the object with the given key. Returns a path-style URL in local environment, otherwise returns
// a normal virtual-host-style object URL.
func (s3Client *s3ClientImpl) GetObjectURL(objectKey string, env string) string {
	if env == "local" {
		endpoint := os.Getenv("AWS_ENDPOINT_URL")
		return fmt.Sprintf("%s/%s/%s", endpoint, s3Client.bucketName, objectKey)
	}

	region := os.Getenv("AWS_REGION")
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Client.bucketName, region, objectKey)
}

// Lists all objects in the bucket which start with the given prefix.
func (s3Client *s3ClientImpl) ListObjects(ctx context.Context, prefix string) (*s3.ListObjectsV2Output, error) {
	output, err := s3Client.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
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
) (*s3.PutObjectOutput, error) {
	output, err := s3Client.client.PutObject(ctx, &s3.PutObjectInput{
		Body:        body,
		Bucket:      aws.String(s3Client.bucketName),
		ContentType: aws.String(contentType),
		Key:         aws.String(objectKey),
	})

	return output, err
}
