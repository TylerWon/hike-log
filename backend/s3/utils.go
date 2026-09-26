package s3

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

// Creates a key for a photo stored in the bucket. Format is "hikes/<hike_id>/photos/<uuid>"
func CreatePhotoObjectKey(hikeId uint) string {
	return fmt.Sprintf("hikes/%d/photos/%s", hikeId, uuid.New())
}

// Creates a S3 URL for a photo object with the given key. Returns a path-style URL in local environment, otherwise
// returns a normal virtual-host-style object URL.
func CreatePhotoObjectURL(objectKey string, env string) string {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	if env == "local" {
		endpoint := os.Getenv("AWS_PUBLIC_ENDPOINT_URL")
		return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectKey)
	}

	region := os.Getenv("AWS_REGION")
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey)
}

// Returns the object key in the given photo object URL. Expects a path-style URL in local environment, otherwise
// expects a normal virtual-host-style object URL.
func GetPhotoObjectKey(objectURL string, env string) string {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	if env == "local" {
		endpoint := os.Getenv("AWS_PUBLIC_ENDPOINT_URL")
		prefix := fmt.Sprintf("%s/%s/", endpoint, bucketName)
		return strings.TrimPrefix(objectURL, prefix)
	}

	region := os.Getenv("AWS_REGION")
	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", bucketName, region)
	return strings.TrimPrefix(objectURL, prefix)
}
