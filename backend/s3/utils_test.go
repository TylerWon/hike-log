package s3

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type utilsTestSuite struct {
	suite.Suite
}

func (suite *utilsTestSuite) TestCreatePhotoObjectKey_ReturnsKeyWithExpectedFormat() {
	hikeId := uint(42)
	objectKey := CreatePhotoObjectKey(hikeId)

	suite.Regexp(fmt.Sprintf(`^hikes/%d/photos/[0-9a-f-]{36}$`, hikeId), objectKey)
}

func (suite *utilsTestSuite) TestCreatePhotoObjectKey_ReturnsUniqueKeys() {
	key1 := CreatePhotoObjectKey(1)
	key2 := CreatePhotoObjectKey(1)

	suite.NotEqual(key1, key2)
}

func (suite *utilsTestSuite) TestCreatePhotoObjectURL_ReturnsPathStyleURLInLocalEnvironment() {
	objectKey := "test"
	awsEndpointURL := os.Getenv("AWS_PUBLIC_ENDPOINT_URL")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	url := CreatePhotoObjectURL(objectKey, "local")
	suite.Equal(fmt.Sprintf("%s/%s/%s", awsEndpointURL, bucketName, objectKey), url)
}

func (suite *utilsTestSuite) TestCreatePhotoObjectURL_ReturnsVirtualHostStyleURLInNonLocalEnvironment() {
	objectKey := "test"
	awsRegion := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	url := CreatePhotoObjectURL(objectKey, "production")
	suite.Equal(fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, objectKey), url)
}

func (suite *utilsTestSuite) TestGetPhotoObjectKey_ReturnsObjectKeyFromPathStyleURLInLocalEnvironment() {
	objectKey := "test"

	objectURL := CreatePhotoObjectURL(objectKey, "local")
	actual := GetPhotoObjectKey(objectURL, "local")
	suite.Equal(objectKey, actual)
}

func (suite *utilsTestSuite) TestGetPhotoObjectKey_ReturnsObjectKeyFromVirtualHostStyleURLInNonLocalEnvironment() {
	objectKey := "test"

	objectURL := CreatePhotoObjectURL(objectKey, "production")
	actual := GetPhotoObjectKey(objectURL, "production")
	suite.Equal(objectKey, actual)
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(utilsTestSuite))
}
