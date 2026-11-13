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

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/models/types"
	"github.com/TylerWon/hike-log/backend/photostorage"
	"github.com/TylerWon/hike-log/backend/s3client"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type hikesHandlerTestSuite struct {
	suite.Suite
	auth         auth.Auth
	middleware   *middleware.Middleware
	photoStorage photostorage.PhotoStorage
	s3Client     s3client.S3Client
	store        store.Store
	suiteDB      *testutils.TestSuiteDB
	router       *gin.Engine
}

func (suite *hikesHandlerTestSuite) SetupSuite() {
	suite.suiteDB = testutils.NewTestSuiteDB(suite.T())
	suite.store = testutils.NewStore(suite.T(), suite.suiteDB.TestDB)
	suite.s3Client = testutils.NewS3Client(suite.T(), "", config.AWSConfig{})
	suite.photoStorage = testutils.NewPhotoStorage(suite.T(), suite.s3Client, "", config.AWSConfig{})
	suite.auth = testutils.NewAuth(suite.T(), config.AuthConfig{})
	suite.middleware = middleware.New(suite.auth)
	suite.router = testutils.NewHikesRouter(suite.T(),
		testutils.NewHikesHandler(suite.T(), suite.store, suite.photoStorage),
		suite.middleware,
	)
}

func (suite *hikesHandlerTestSuite) TearDownSuite() {
	suite.suiteDB.Teardown(suite.T())
}

func (suite *hikesHandlerTestSuite) SetupTest() {
	suite.suiteDB.Reset(suite.T())
}

func TestHikesTestSuite(t *testing.T) {
	suite.Run(t, new(hikesHandlerTestSuite))
}

func (suite *hikesHandlerTestSuite) TestCreateHike_Returns401ErrorWhenUnauthorized() {
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
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodPost, "/hikes", reqBody) // Missing session cookie

	suite.Require().Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreateHike_Returns400ErrorWhenRequestBodyHasInvalidFields() {
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
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreateHike_Returns400ErrorWhenRequestBodyIsMissingFields() {
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
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreateHike_Returns500ErrorWhenDBErrors() {
	mockStore := store.MockStore{CreateRecordError: errors.New("Something went wrong")}
	hikesHandler := testutils.NewHikesHandler(suite.T(), &mockStore, suite.photoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

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
	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreateHike_CreatesAndReturnsHike() {
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
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	expected := models.Hike{
		ID:            response.ID,
		TrailName:     "Trail 1",
		Date:          types.Date("2026-02-05"),
		Notes:         "Easy hike",
		Rating:        4.5,
		Difficulty:    3,
		Distance:      8.2,
		ElevationGain: 1200,
		Duration:      60,
		AllTrailsUrl:  "https://www.alltrails.com/",
		Photos:        []models.Photo{},
	}
	suite.Require().Equal(expected, response)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_Returns401ErrorWhenUnauthorized() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	objectKey1 := suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID)
	_, err := suite.s3Client.PutObject(context.TODO(), objectKey1, strings.NewReader("content"), "image/png")
	suite.Require().NoError(err)

	objectKey2 := suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID)
	_, err = suite.s3Client.PutObject(context.TODO(), objectKey2, strings.NewReader("content"), "image/jpeg")
	suite.Require().NoError(err)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    objectKey1,
			"caption":      "Caption 1",
			"displayOrder": 1,
		},
		{
			"index":        1,
			"objectKey":    objectKey2,
			"caption":      "Caption 2",
			"displayOrder": 2,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
	) // Missing session cookie

	suite.Require().Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_Returns400ErrorWhenHikeIDIsInvalid() {
	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes/abc/photos",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_Returns404ErrorWhenHikeDoesNotExist() {
	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes/1/photos",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_Returns400ErrorWhenRequestItemIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"caption":      "Caption",
			"displayOrder": 1,
		},
		{
			"index":        1,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(2),
			"caption":      "Caption",
			"displayOrder": 2,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_Returns400ErrorWhenObjectKeyIsInvalid() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    "acde070d-8c4c-4f0d-9d8a-162843c10333",
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_Returns400ErrorWhenPathAndObjectKeyHikeIDMismatch() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID + 1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenPhotoDoesNotExistInS3() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Len(response, 1)
	suite.Require().Equal(uint(0), response[0].Index)
	suite.Require().False(response[0].Success)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenS3Errors() {
	mockPhotoStorage := photostorage.MockPhotoStorage{
		DoesPhotoExistResult: false,
		DoesPhotoExistError:  errors.New("Something went wrong"),
	}
	hikesHandler := testutils.NewHikesHandler(suite.T(), suite.store, &mockPhotoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":        2,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Len(response, 1)
	suite.Require().Equal(uint(2), response[0].Index)
	suite.Require().False(response[0].Success)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenDBErrors() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	mockStore := store.MockStore{
		GetHikeByIDResult: &hikes[0],
		CreateRecordError: errors.New("Something went wrong"),
	}
	mockPhotoStorage := photostorage.MockPhotoStorage{DoesPhotoExistResult: true}
	hikesHandler := testutils.NewHikesHandler(suite.T(), &mockStore, &mockPhotoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	body := [](map[string]any){
		{
			"index":        4,
			"objectKey":    suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Len(response, 1)
	suite.Require().Equal(uint(4), response[0].Index)
	suite.Require().False(response[0].Success)
}

func (suite *hikesHandlerTestSuite) TestCreatePhotos_CreatesPhotos() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	objectKey1 := suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID)
	_, err := suite.s3Client.PutObject(context.TODO(), objectKey1, strings.NewReader("content"), "image/png")
	suite.Require().NoError(err)

	objectKey2 := suite.photoStorage.CreatePhotoObjectKey(hikes[0].ID)
	_, err = suite.s3Client.PutObject(context.TODO(), objectKey2, strings.NewReader("content"), "image/jpeg")
	suite.Require().NoError(err)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    objectKey1,
			"caption":      "Caption 1",
			"displayOrder": 1,
		},
		{
			"index":        1,
			"objectKey":    objectKey2,
			"caption":      "Caption 2",
			"displayOrder": 2,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err = json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Len(response, 2)
	objectKeys := []string{objectKey1, objectKey2}
	for idx, item := range response {
		suite.Require().Equal(uint(idx), item.Index)
		suite.Require().True(item.Success)
		suite.Require().Empty(item.Error)
		suite.Require().NotEmpty(item.Result)

		actual := item.Result
		expected := models.Photo{
			ID:           actual.ID,
			SrcUrl:       suite.photoStorage.GetPhotoObjectURL(objectKeys[idx]),
			Caption:      fmt.Sprintf("Caption %d", idx+1),
			DisplayOrder: uint(idx + 1),
			HikeID:       hikes[0].ID,
		}
		suite.Require().Equal(expected, *actual)
	}

	_, err = suite.s3Client.DeleteObject(context.TODO(), objectKey1)
	suite.Require().NoError(err)

	_, err = suite.s3Client.DeleteObject(context.TODO(), objectKey2)
	suite.Require().NoError(err)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns401ErrorWhenUnauthorized() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
		{
			"index":         1,
			"contentType":   "image/png",
			"contentLength": 200,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
	) // Missing session cookie

	suite.Require().Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenHikeIDIsInvalid() {
	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes/abc/photos/presigned-urls",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns404ErrorWhenHikeDoesNotExist() {
	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		"/hikes/1/photos/presigned-urls",
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenRequestItemIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenRequestItemHasAnInvalidImageType() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "text/html",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenContentLengthIsOutsideBounds() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": -1, // Below lower bound
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)

	body = [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 10485761, // Above upper bound
		},
	}
	reqBody = testutils.SerializeJSONRequestBody(suite.T(), body)
	res = testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_Returns500ErrorWhenDBErrors() {
	mockStore := store.MockStore{
		GetHikeByIDResult: nil,
		GetHikeByIDError:  errors.New("Something went wrong"),
	}
	hikesHandler := testutils.NewHikesHandler(suite.T(), &mockStore, suite.photoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_ReturnsPartialErrorWhenS3Errors() {
	mockPhotoStorage := photostorage.MockPhotoStorage{
		CreatePresignedUploadURLError: errors.New("Something went wrong"),
	}
	hikesHandler := testutils.NewHikesHandler(suite.T(), suite.store, &mockPhotoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         2,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response handler.CreatePresignedURLsResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Equal(1, len(response))

	item := response[0]
	suite.Require().Equal(uint(2), item.Index)
	suite.Require().False(item.Success)
	suite.Require().Empty(item.Result)
	suite.Require().NotEmpty(item.Error)
}

func (suite *hikesHandlerTestSuite) TestCreatePresignedURLs_ReturnsPresignedURLs() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
		{
			"index":         1,
			"contentType":   "image/png",
			"contentLength": 200,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusCreated, res.Code)

	var response handler.CreatePresignedURLsResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Equal(2, len(response))
	for idx, item := range response {
		suite.Require().Equal(uint(idx), item.Index)
		suite.Require().True(item.Success)
		suite.Require().Empty(item.Error)
		suite.Require().NotEmpty(item.Result)

		result := item.Result
		suite.Require().Regexp(fmt.Sprintf("^hikes/%d/photos/[0-9a-f-]{36}$", hikes[0].ID), result.ObjectKey)
		parsedPresignedURL, err := url.Parse(result.PresignedURL)
		suite.Require().NoError(err)
		suite.Require().Contains(parsedPresignedURL.Path, result.ObjectKey)
	}
}

func (suite *hikesHandlerTestSuite) TestDeleteHike_Returns401ErrorWhenUnauthorized() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, true, true)[0]

	photo := hike.Photos[0]
	objectKey := suite.photoStorage.GetPhotoObjectKey(photo.SrcUrl)
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.Require().NoError(err)

	res := testutils.SendRequest(suite.T(), suite.router, http.MethodDelete, fmt.Sprintf("/hikes/%d", hike.ID), nil) // Missing session cookie
	suite.Require().Equal(http.StatusUnauthorized, res.Code)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *hikesHandlerTestSuite) TestDeleteHike_Returns400ErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodDelete,
		"/hikes/abc",
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestDeleteHike_Returns404ErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodDelete,
		"/hikes/1",
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesHandlerTestSuite) TestDeleteHike_Returns500ErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, true, true)[0]

	mockStore := store.MockStore{DeleteRecordError: errors.New("Something went wrong")}
	hikesHandler := testutils.NewHikesHandler(suite.T(), &mockStore, suite.photoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodDelete,
		fmt.Sprintf("/hikes/%d", hike.ID),
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesHandlerTestSuite) TestDeleteHike_DeletesHikeAndPhotosWhenS3Errors() {
	mockPhotoStorage := photostorage.MockPhotoStorage{
		DeletePhotosForHikeError: errors.New("Something went wrong"),
	}
	hikesHandler := testutils.NewHikesHandler(suite.T(), suite.store, &mockPhotoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, true, true)[0]

	photo := hike.Photos[0]
	objectKey := suite.photoStorage.GetPhotoObjectKey(photo.SrcUrl)
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.Require().NoError(err)

	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodDelete,
		fmt.Sprintf("/hikes/%d", hike.ID),
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetHikeByID(hike.ID)
	suite.Require().Error(err)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Require().Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().True(exists)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *hikesHandlerTestSuite) TestDeleteHike_DeletesHikeAndPhotosAndS3Objects() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, true, true)[0]

	photo := hike.Photos[0]
	objectKey := suite.photoStorage.GetPhotoObjectKey(photo.SrcUrl)
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.Require().NoError(err)

	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodDelete,
		fmt.Sprintf("/hikes/%d", hike.ID),
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetHikeByID(hike.ID)
	suite.Require().Error(err)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Require().Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.Require().NoError(err)
	suite.Require().True(exists)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *hikesHandlerTestSuite) TestListHike_Returns500ErrorWhenDBErrors() {

	mockStore := store.MockStore{
		ListHikesResult: nil,
		ListHikesError:  errors.New("Something went wrong"),
	}
	hikesHandler := testutils.NewHikesHandler(suite.T(), &mockStore, suite.photoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

	testutils.ConstructHikes(suite.T(), 2, suite.store, suite.photoStorage, true, true)

	res := testutils.SendRequest(suite.T(), router, http.MethodGet, "/hikes", nil)

	suite.Require().Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesHandlerTestSuite) TestListHike_ReturnsNothingWhenThereAreNoHikes() {
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodGet, "/hikes", nil)

	suite.Require().Equal(http.StatusOK, res.Code)

	var hikes []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &hikes)
	suite.Require().NoError(err)

	suite.Require().Len(hikes, 0)
}

func (suite *hikesHandlerTestSuite) TestListHike_ReturnsHikes() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, suite.photoStorage, true, true)

	res := testutils.SendRequest(suite.T(), suite.router, http.MethodGet, "/hikes", nil)

	suite.Require().Equal(http.StatusOK, res.Code)

	var response []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().Len(response, 2)
	suite.Require().Equal(hikes[1], response[0])
	suite.Require().Equal(hikes[0], response[1])
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_Returns401ErrorWhenUnauthorized() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)[0]

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
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodPut, fmt.Sprintf("/hikes/%d", hike.ID), reqBody) // Missing session cookie

	suite.Require().Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_Returns400ErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPut,
		"/hikes/abc",
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_Returns404ErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPut,
		"/hikes/1",
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)
	suite.Require().Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_Returns400ErrorWhenRequestBodyHasInvalidFields() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)[0]

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
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_Returns400ErrorWhenRequestBodyIsMissingFields() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)[0]

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
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_Returns500ErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)[0]

	mockStore := store.MockStore{
		GetHikeByIDResult: &hike,
		UpdateRecordError: errors.New("Something went wrong"),
	}
	hikesHandler := testutils.NewHikesHandler(suite.T(), &mockStore, suite.photoStorage)
	router := testutils.NewHikesRouter(suite.T(), hikesHandler, suite.middleware)

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
	res := testutils.SendRequest(suite.T(),
		router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesHandlerTestSuite) TestUpdateHike_UpdatesHike() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, suite.photoStorage, false, true)[0]

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
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusOK, res.Code)

	updated, err := suite.store.GetHikeByID(hike.ID)
	suite.Require().NoError(err)

	expected := models.Hike{
		ID:            hike.ID,
		TrailName:     "Updated Trail",
		Date:          types.Date("2026-03-15"),
		Notes:         "Updated notes",
		Rating:        4.5,
		Difficulty:    3,
		Distance:      8.2,
		ElevationGain: 1200,
		Duration:      60,
		AllTrailsUrl:  "https://www.alltrails.com/updated",
	}
	suite.Require().Equal(expected, *updated)
}
