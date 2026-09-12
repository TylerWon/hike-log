// By adding _test to the package name, this file belongs to its own package separate from handler. This ensures the
// tests can only access the public API exposed by the handler package.
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
		"objectKey": "hikes/1/photos/acde070d-8c4c-4f0d-9d8a-162843c10333",
		"caption":   "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/abc/photos/", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenHikeDoesNotExist() {
	body := map[string]any{
		"objectKey": "hikes/1/photos/acde070d-8c4c-4f0d-9d8a-162843c10333",
		"caption":   "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/1/photos/", reqBody)

	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenRequestBodyIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"caption": "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenObjectKeyIsInvalid() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey": "acde070d-8c4c-4f0d-9d8a-162843c10333",
		"caption":   "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenPathAndObjectKeyHikeIDMismatch() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey": fmt.Sprintf("hikes/%d/photos/acde070d-8c4c-4f0d-9d8a-162843c10333", hikes[0].ID+1),
		"caption":   "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_ReturnsErrorWhenPhotoDoesNotExistInS3() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := map[string]any{
		"objectKey": fmt.Sprintf("hikes/%d/photos/acde070d-8c4c-4f0d-9d8a-162843c10333", hikes[0].ID),
		"caption":   "Caption",
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
		"objectKey": fmt.Sprintf("hikes/%d/photos/acde070d-8c4c-4f0d-9d8a-162843c10333", hikes[0].ID),
		"caption":   "Caption",
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
		"objectKey": fmt.Sprintf("hikes/%d/photos/acde070d-8c4c-4f0d-9d8a-162843c10333", hikes[0].ID),
		"caption":   "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhoto_CreatesPhoto() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	objectKey := fmt.Sprintf("hikes/%d/photos/acde070d-8c4c-4f0d-9d8a-162843c10333", hikes[0].ID)
	_, err := suite.s3Client.PutObject(context.TODO(), objectKey, strings.NewReader("content"), "image/png")
	suite.NoError(err)

	body := map[string]any{
		"objectKey": objectKey,
		"caption":   "Caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/", hikes[0].ID), reqBody)

	suite.Equal(http.StatusCreated, res.Code)

	var response models.Photo
	err = json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	expected := models.Photo{
		ID:      response.ID,
		SrcUrl:  fmt.Sprintf("http://localstack:4566/hike-log/%s", objectKey),
		Caption: "Caption",
		HikeID:  hikes[0].ID,
	}
	suite.Equal(expected, response)

	_, err = suite.s3Client.DeleteObject(context.TODO(), objectKey)
	suite.NoError(err)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
