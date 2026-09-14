package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	aws_sdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type handlerTestSuite struct {
	suite.Suite
	router   *gin.Engine
	suiteDB  *testutils.TestSuiteDB
	store    store.Store
	s3Client aws.S3Client
}

func (suite *handlerTestSuite) SetupSuite() {
	suite.suiteDB = testutils.NewTestSuiteDB(suite.T())
	suite.store = testutils.NewStore(suite.T(), suite.suiteDB)
}

func (suite *handlerTestSuite) TearDownSuite() {
	testutils.TeardownStore(suite.T(), suite.store)
	suite.suiteDB.Teardown(suite.T())
}

func (suite *handlerTestSuite) SetupTest() {
	suite.suiteDB.Reset(suite.T())
	suite.s3Client = testutils.NewS3Client(suite.T())
	handler := handler.New(suite.store, suite.s3Client)
	suite.router = testutils.NewRouter(suite.T(), handler)
}

func (suite *handlerTestSuite) TestListHike_ReturnsErrorWhenDBErrors() {

	mockStore := store.MockStore{
		ListHikesResult: nil,
		ListHikesError:  errors.New("Something went wrong"),
	}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)

	res := testutils.SendRequest(router, http.MethodGet, "/api/v1/hikes", nil)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestListHike_ReturnsNothingWhenThereAreNoHikes() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/hikes", nil)

	suite.Equal(http.StatusOK, res.Code)

	var hikes []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &hikes)
	suite.NoError(err)

	suite.Len(hikes, 0)
}

func (suite *handlerTestSuite) TestListHike_ReturnsHikes() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)

	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/hikes", nil)

	suite.Equal(http.StatusOK, res.Code)

	var response []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 2)
	suite.Equal(hikes[1], response[0])
	suite.Equal(hikes[0], response[1])
}

func (suite *handlerTestSuite) TestCreateHike_ReturnsErrorWhenRequestBodyHasInvalidFields() {
	body := map[string]any{
		"trailName":     "Trail 1",
		"date":          "05-02-2026", // DD-MM-YYYY instead of YYYY-MM-DD
		"notes":         "Easy hike",
		"rating":        4.5,
		"difficulty":    3,
		"distance":      "8.2 km", // string instead of float
		"elevationGain": 1200,
		"duration":      60.5, // float instead of uint
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreateHike_ReturnsErrorWhenRequestBodyIsMissingFields() {
	body := map[string]any{
		"date":          "2026-02-05",
		"notes":         "Easy hike",
		"rating":        4.5,
		"distance":      8.2,
		"elevationGain": 1200,
		"duration":      60,
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreateHike_ReturnsErrorWhenDBErrors() {
	mockStore := store.MockStore{CreateModelError: errors.New("Something went wrong")}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	body := map[string]any{
		"trailName":     "Trail 1",
		"date":          "2026-02-05",
		"notes":         "Easy hike",
		"rating":        4.5,
		"difficulty":    3,
		"distance":      8.2,
		"elevationGain": 1200,
		"duration":      60,
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, "/api/v1/hikes", reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreateHike_CreatesAndReturnsHike() {
	body := map[string]any{
		"trailName":     "Trail 1",
		"date":          "2026-02-05",
		"notes":         "Easy hike",
		"rating":        4.5,
		"difficulty":    3,
		"distance":      8.2,
		"elevationGain": 1200,
		"duration":      60,
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes", reqBody)

	suite.Equal(http.StatusCreated, res.Code)

	var response models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	expected := models.Hike{
		ID:            response.ID,
		TrailName:     "Trail 1",
		Date:          datatypes.Date(time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC)),
		Notes:         "Easy hike",
		Rating:        4.5,
		Difficulty:    3,
		Distance:      8.2,
		ElevationGain: 1200,
		Duration:      60,
		AllTrailsUrl:  "https://www.alltrails.com/",
	}
	suite.Equal(expected, response)
}

func (suite *handlerTestSuite) TestDeleteHike_ReturnsErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodDelete, "/api/v1/hikes/abc/", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestDeleteHike_ReturnsErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(suite.router, http.MethodDelete, "/api/v1/hikes/1/", nil)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestDeleteHike_ReturnsErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	mockStore := store.MockStore{DeleteModelError: errors.New("Something went wrong")}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	res := testutils.SendRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), nil)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestDeleteHike_DeletesHikeAndPhotosWhenS3Errors() {
	mockS3Client := aws.MockS3Client{
		ListObjectsResult:   &s3.ListObjectsV2Output{Contents: []types.Object{{Key: aws_sdk.String("test")}}},
		DeleteObjectsResult: nil,
		DeleteObjectsError:  errors.New("Something went wrong"),
	}
	handler := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), handler)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	photo := hike.Photos[0]
	objectKey := suite.s3Client.GetObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), nil)
	suite.Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetHikeByID(hike.ID)
	suite.Error(err)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *handlerTestSuite) TestDeleteHike_DeletesHikeAndPhotosAndS3Objects() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	photo := hike.Photos[0]
	objectKey := suite.s3Client.GetObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(suite.router, http.MethodDelete, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), nil)
	suite.Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetHikeByID(hike.ID)
	suite.Error(err)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodPut, "/api/v1/hikes/abc/", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(suite.router, http.MethodPut, "/api/v1/hikes/1/", nil)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenRequestBodyHasInvalidFields() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	body := map[string]any{
		"trailName":     "Trail 1",
		"date":          "05-02-2026", // DD-MM-YYYY instead of YYYY-MM-DD
		"notes":         "Easy hike",
		"rating":        4.5,
		"difficulty":    3,
		"distance":      "8.2 km", // string instead of float
		"elevationGain": 1200,
		"duration":      60.5, // float instead of uint
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenRequestBodyIsMissingFields() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	body := map[string]any{
		"date":          "2026-02-05",
		"notes":         "Easy hike",
		"rating":        4.5,
		"distance":      8.2,
		"elevationGain": 1200,
		"duration":      60,
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	mockStore := store.MockStore{
		GetHikeByIDResult: &hike,
		UpdateModelError:  errors.New("Something went wrong"),
	}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	body := map[string]any{
		"trailName":     "Updated Trail",
		"date":          "2026-02-05",
		"notes":         "Updated notes",
		"rating":        4.5,
		"difficulty":    3,
		"distance":      8.2,
		"elevationGain": 1200,
		"duration":      60,
		"allTrailsUrl":  "https://www.alltrails.com/",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_UpdatesHike() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	body := map[string]any{
		"trailName":     "Updated Trail",
		"date":          "2026-03-15",
		"notes":         "Updated notes",
		"rating":        4.5,
		"difficulty":    3,
		"distance":      8.2,
		"elevationGain": 1200,
		"duration":      60,
		"allTrailsUrl":  "https://www.alltrails.com/updated",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d/", hike.ID), reqBody)

	suite.Equal(http.StatusOK, res.Code)

	updated, err := suite.store.GetHikeByID(hike.ID)
	suite.NoError(err)

	expected := models.Hike{
		ID:            hike.ID,
		TrailName:     "Updated Trail",
		Date:          datatypes.Date(time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)),
		Notes:         "Updated notes",
		Rating:        4.5,
		Difficulty:    3,
		Distance:      8.2,
		ElevationGain: 1200,
		Duration:      60,
		AllTrailsUrl:  "https://www.alltrails.com/updated",
	}
	suite.Equal(expected, *updated)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenHikeIDIsInvalid() {
	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/abc/photos/upload-url", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenHikeDoesNotExist() {
	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/1/photos/upload-url", reqBody)

	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenRequestBodyIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenContentTypeIsAnInvalidImageType() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"contentType":   "text/html",
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenContentLengthIsOutsideBounds() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": -1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)

	body = map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 10485761,
	}
	reqBody = testutils.SerializeJSONRequestBody(suite.T(), body)
	res = testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenS3Errors() {
	mockS3Client := aws.MockS3Client{
		CreatePresignedPutObjectRequestResult: nil,
		CreatePresignedPutObjectRequestError:  errors.New("Something went wrong"),
	}
	handler := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), handler)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenDBErrors() {
	mockStore := store.MockStore{
		GetHikeByIDResult: nil,
		GetHikeByIDError:  errors.New("Something went wrong"),
	}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsUploadURL() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CreatePhotoUploadURLResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Regexp(fmt.Sprintf("^hikes/%d/photos/[0-9a-f-]{36}$", hikes[0].ID), response.ObjectKey)

	parsedUploadURL, err := url.Parse(response.UploadURL)
	suite.NoError(err)
	suite.Contains(parsedUploadURL.Path, response.ObjectKey)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenHikeIDIsInvalid() {
	body := map[string]any{
		"objectKey":    suite.s3Client.CreatePhotoObjectKey(1),
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/abc/photos/", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenHikeDoesNotExist() {
	body := map[string]any{
		"objectKey":    suite.s3Client.CreatePhotoObjectKey(1),
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/1/photos/", reqBody)

	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenRequestBodyIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenObjectKeyIsInvalid() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey":    "acde070d-8c4c-4f0d-9d8a-162843c10333",
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenPathAndObjectKeyHikeIDMismatch() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey":    suite.s3Client.CreatePhotoObjectKey(hikes[0].ID + 1),
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenPhotoDoesNotExistInS3() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey":    suite.s3Client.CreatePhotoObjectKey(hikes[0].ID),
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenS3Errors() {
	mockS3Client := aws.MockS3Client{
		DoesObjectExistResult: false,
		DoesObjectExistError:  errors.New("Something went wrong"),
	}
	handler := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), handler)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey":    suite.s3Client.CreatePhotoObjectKey(hikes[0].ID),
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenDBErrors() {
	mockStore := store.MockStore{
		GetHikeByIDResult: nil,
		GetHikeByIDError:  errors.New("Something went wrong"),
	}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey":    suite.s3Client.CreatePhotoObjectKey(hikes[0].ID),
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_CreatesPhoto() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	objectKey := suite.s3Client.CreatePhotoObjectKey(hikes[0].ID)
	_, err := suite.s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "image/png")
	suite.NoError(err)

	body := map[string]any{
		"objectKey":    objectKey,
		"caption":      "Caption",
		"displayOrder": 1,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusCreated, res.Code)

	var response models.Photo
	err = json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	expected := models.Photo{
		ID:           response.ID,
		SrcUrl:       fmt.Sprintf("http://localstack:4566/hike-log/%s", objectKey),
		Caption:      "Caption",
		DisplayOrder: 1,
		HikeID:       hikes[0].ID,
	}
	suite.Equal(expected, response)

	_, err = suite.s3Client.DeleteObject(context.TODO(), objectKey)
	suite.NoError(err)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
