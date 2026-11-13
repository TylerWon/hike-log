package photostorage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/s3client"
	"github.com/google/uuid"
)

// PhotoStorage handles photo object operations in the S3 bucket for the app.
type PhotoStorage interface {
	CreatePhotoObjectKey(hikeID uint) string
	CreatePresignedUploadURL(
		ctx context.Context,
		hikeID uint,
		contentType string,
		contentLength int64,
	) (presignedURL string, objectKey string, err error)
	DeletePhoto(ctx context.Context, objectURL string) error
	DeletePhotosForHike(ctx context.Context, hikeID uint) error
	DoesPhotoExist(ctx context.Context, objectKey string) (bool, error)
	GetPhotoObjectKey(objectURL string) string
	GetPhotoObjectURL(objectKey string) string
}

// photoStorageImpl is an implementation of the PhotoStorage interface.
type photoStorageImpl struct {
	awsConfig config.AWSConfig
	env       string
	s3Client  s3client.S3Client
}

// Creates a new PhotoStorage.
func New(s3Client s3client.S3Client, env string, awsConfig config.AWSConfig) (PhotoStorage, error) {
	if env == "" {
		return nil, errors.New("env is empty")
	} else if awsConfig.S3.BucketName == "" {
		return nil, errors.New("bucketName is empty")
	} else if awsConfig.Region == "" {
		return nil, errors.New("region is empty")
	}

	if env == "local" {
		if awsConfig.ExternalEndpointURL == "" {
			return nil, errors.New("externalEndpointURL is empty in local environment")
		}
	}

	return &photoStorageImpl{
		awsConfig: awsConfig,
		env:       env,
		s3Client:  s3Client,
	}, nil
}

// Creates a key for a photo stored in the bucket. Format is "hikes/<hike_id>/photos/<uuid>".
func (p *photoStorageImpl) CreatePhotoObjectKey(hikeID uint) string {
	return fmt.Sprintf("hikes/%d/photos/%s", hikeID, uuid.New())
}

// Creates a presigned URL that can be used to upload a photo for a Hike. Also returns the object key that will be
// written to.
func (p *photoStorageImpl) CreatePresignedUploadURL(
	ctx context.Context,
	hikeID uint,
	contentType string,
	contentLength int64,
) (string, string, error) {
	objectKey := p.CreatePhotoObjectKey(hikeID)
	presignedReq, err := p.s3Client.CreatePresignedPutObjectRequest(ctx, objectKey, contentType, contentLength)
	if err != nil {
		return "", "", err
	}
	return presignedReq.URL, objectKey, nil
}

// Deletes a photo object identified by its object URL.
func (p *photoStorageImpl) DeletePhoto(ctx context.Context, objectURL string) error {
	objectKey := p.GetPhotoObjectKey(objectURL)
	_, err := p.s3Client.DeleteObject(ctx, objectKey)
	return err
}

// Deletes all photo objects for a Hike. No-op when there are no objects for that Hike.
func (p *photoStorageImpl) DeletePhotosForHike(ctx context.Context, hikeID uint) error {
	objectKeyPrefix := fmt.Sprintf("hikes/%d/photos/", hikeID)
	result, err := p.s3Client.ListObjects(ctx, objectKeyPrefix)
	if err != nil {
		return err
	} else if len(result.Contents) == 0 {
		return nil
	}

	var objectKeys []string
	for _, object := range result.Contents {
		objectKeys = append(objectKeys, *object.Key)
	}

	_, err = p.s3Client.DeleteObjects(ctx, objectKeys)
	return err
}

// Checks if a photo object with the given key exists in the bucket.
func (p *photoStorageImpl) DoesPhotoExist(ctx context.Context, objectKey string) (bool, error) {
	return p.s3Client.DoesObjectExist(ctx, objectKey)
}

// Returns the object key from a photo object URL.
func (p *photoStorageImpl) GetPhotoObjectKey(objectURL string) string {
	if p.env == "local" {
		// objectURL is in path-style for local environments
		prefix := fmt.Sprintf("%s/%s/", p.awsConfig.ExternalEndpointURL, p.awsConfig.S3.BucketName)
		return strings.TrimPrefix(objectURL, prefix)
	}
	// objectURL is in virtual-hosted-style for non-local environments
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", p.awsConfig.S3.BucketName, p.awsConfig.Region)
	return strings.TrimPrefix(objectURL, prefix)
}

// Returns the object URL for a photo object with the given key.
func (p *photoStorageImpl) GetPhotoObjectURL(objectKey string) string {
	if p.env == "local" {
		// objectURL should be path-style for local environments
		return fmt.Sprintf("%s/%s/%s", p.awsConfig.ExternalEndpointURL, p.awsConfig.S3.BucketName, objectKey)
	}
	// objectURL should be virtual-hosted-style for non-local environments
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.awsConfig.S3.BucketName, p.awsConfig.Region, objectKey)
}
