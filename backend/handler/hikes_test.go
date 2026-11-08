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
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/models/types"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	s3Sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type hikesTestSuite struct {
	suite.Suite
	auth       auth.Auth
	middleware *middleware.Middleware
	s3Client   s3.S3Client
	store      store.Store
	suiteDB    *testutils.TestSuiteDB
	router     *gin.Engine
}

func (suite *hikesTestSuite) SetupSuite() {
	suite.suiteDB = testutils.NewTestSuiteDB(suite.T())
	suite.store = testutils.NewStore(suite.T(), suite.suiteDB)
	suite.s3Client = testutils.NewS3Client(suite.T())
	suite.auth = testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, testutils.SESSION_SECRET)
	suite.middleware = middleware.New(suite.auth)
	suite.router = testutils.NewHikesRouter(handler.NewHikesHandler(suite.store, suite.s3Client), suite.middleware)
}

func (suite *hikesTestSuite) TearDownSuite() {
	testutils.TeardownStore(suite.T(), suite.store)
	suite.suiteDB.Teardown(suite.T())
}

func (suite *hikesTestSuite) SetupTest() {
	suite.suiteDB.Reset(suite.T())
}

func TestHikesTestSuite(t *testing.T) {
	suite.Run(t, new(hikesTestSuite))
}

func (suite *hikesTestSuite) TestCreateHike_Returns401ErrorWhenUnauthorized() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, "/hikes", reqBody) // Missing session cookie

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesTestSuite) TestCreateHike_Returns400ErrorWhenRequestBodyHasInvalidFields() {
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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreateHike_Returns400ErrorWhenRequestBodyIsMissingFields() {
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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreateHike_Returns500ErrorWhenDBErrors() {
	mockStore := store.MockStore{CreateRecordError: errors.New("Something went wrong")}
	hikesHandler := handler.NewHikesHandler(&mockStore, suite.s3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

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
	res := testutils.SendRequest(
		router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesTestSuite) TestCreateHike_CreatesAndReturnsHike() {
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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

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
	suite.Equal(expected, response)
}

func (suite *hikesTestSuite) TestCreatePhotos_Returns401ErrorWhenUnauthorized() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	objectKey1 := s3.CreatePhotoObjectKey(hikes[0].ID)
	_, err := suite.s3Client.PutObject(context.TODO(), objectKey1, strings.NewReader("content"), "image/png")
	suite.NoError(err)

	objectKey2 := s3.CreatePhotoObjectKey(hikes[0].ID)
	_, err = suite.s3Client.PutObject(context.TODO(), objectKey2, strings.NewReader("content"), "image/jpeg")
	suite.NoError(err)

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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
	) // Missing session cookie

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesTestSuite) TestCreatePhotos_Returns400ErrorWhenHikeIDIsInvalid() {
	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    s3.CreatePhotoObjectKey(1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes/abc/photos",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePhotos_Returns404ErrorWhenHikeDoesNotExist() {
	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    s3.CreatePhotoObjectKey(1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes/1/photos",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesTestSuite) TestCreatePhotos_Returns400ErrorWhenRequestItemIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"caption":      "Caption",
			"displayOrder": 1,
		},
		{
			"index":        1,
			"objectKey":    s3.CreatePhotoObjectKey(2),
			"caption":      "Caption",
			"displayOrder": 2,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePhotos_Returns400ErrorWhenObjectKeyIsInvalid() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    "acde070d-8c4c-4f0d-9d8a-162843c10333",
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePhotos_Returns400ErrorWhenPathAndObjectKeyHikeIDMismatch() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    s3.CreatePhotoObjectKey(hikes[0].ID + 1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenPhotoDoesNotExistInS3() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    s3.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 1)
	suite.Equal(uint(0), response[0].Index)
	suite.False(response[0].Success)
}

func (suite *hikesTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		DoesObjectExistResult: false,
		DoesObjectExistError:  errors.New("Something went wrong"),
	}
	hikesHandler := handler.NewHikesHandler(suite.store, &mockS3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":        2,
			"objectKey":    s3.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 1)
	suite.Equal(uint(2), response[0].Index)
	suite.False(response[0].Success)
}

func (suite *hikesTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenDBErrors() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	mockStore := store.MockStore{
		GetHikeByIDResult: &hikes[0],
		CreateRecordError: errors.New("Something went wrong"),
	}
	mockS3Client := s3.MockS3Client{DoesObjectExistResult: true}
	hikesHandler := handler.NewHikesHandler(&mockStore, &mockS3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	body := [](map[string]any){
		{
			"index":        4,
			"objectKey":    s3.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 1)
	suite.Equal(uint(4), response[0].Index)
	suite.False(response[0].Success)
}

func (suite *hikesTestSuite) TestCreatePhotos_CreatesPhotos() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	objectKey1 := s3.CreatePhotoObjectKey(hikes[0].ID)
	_, err := suite.s3Client.PutObject(context.TODO(), objectKey1, strings.NewReader("content"), "image/png")
	suite.NoError(err)

	objectKey2 := s3.CreatePhotoObjectKey(hikes[0].ID)
	_, err = suite.s3Client.PutObject(context.TODO(), objectKey2, strings.NewReader("content"), "image/jpeg")
	suite.NoError(err)

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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err = json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 2)
	objectKeys := []string{objectKey1, objectKey2}
	for idx, item := range response {
		suite.Equal(uint(idx), item.Index)
		suite.True(item.Success)
		suite.Empty(item.Error)
		suite.NotEmpty(item.Result)

		actual := item.Result
		expected := models.Photo{
			ID:           actual.ID,
			SrcUrl:       s3.CreatePhotoObjectURL(objectKeys[idx], "local"),
			Caption:      fmt.Sprintf("Caption %d", idx+1),
			DisplayOrder: uint(idx + 1),
			HikeID:       hikes[0].ID,
		}
		suite.Equal(expected, *actual)
	}

	_, err = suite.s3Client.DeleteObject(context.TODO(), objectKey1)
	suite.NoError(err)

	_, err = suite.s3Client.DeleteObject(context.TODO(), objectKey2)
	suite.NoError(err)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns401ErrorWhenUnauthorized() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
	) // Missing session cookie

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenHikeIDIsInvalid() {
	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes/abc/photos/presigned-urls",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns404ErrorWhenHikeDoesNotExist() {
	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		"/hikes/1/photos/presigned-urls",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenRequestItemIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenRequestItemHasAnInvalidImageType() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "text/html",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns400ErrorWhenContentLengthIsOutsideBounds() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": -1, // Below lower bound
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)

	body = [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 10485761, // Above upper bound
		},
	}
	reqBody = testutils.SerializeJSONRequestBody(suite.T(), body)
	res = testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_Returns500ErrorWhenDBErrors() {
	mockStore := store.MockStore{
		GetHikeByIDResult: nil,
		GetHikeByIDError:  errors.New("Something went wrong"),
	}
	hikesHandler := handler.NewHikesHandler(&mockStore, suite.s3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_ReturnsPartialErrorWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		CreatePresignedPutObjectRequestResult: nil,
		CreatePresignedPutObjectRequestError:  errors.New("Something went wrong"),
	}
	hikesHandler := handler.NewHikesHandler(suite.store, &mockS3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         2,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePresignedURLsResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal(1, len(response))

	item := response[0]
	suite.Equal(uint(2), item.Index)
	suite.False(item.Success)
	suite.Empty(item.Result)
	suite.NotEmpty(item.Error)
}

func (suite *hikesTestSuite) TestCreatePresignedURLs_ReturnsPresignedURLs() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPost,
		fmt.Sprintf("/hikes/%d/photos/presigned-urls", hikes[0].ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePresignedURLsResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Equal(2, len(response))
	for idx, item := range response {
		suite.Equal(uint(idx), item.Index)
		suite.True(item.Success)
		suite.Empty(item.Error)
		suite.NotEmpty(item.Result)

		result := item.Result
		suite.Regexp(fmt.Sprintf("^hikes/%d/photos/[0-9a-f-]{36}$", hikes[0].ID), result.ObjectKey)
		parsedPresignedURL, err := url.Parse(result.PresignedURL)
		suite.NoError(err)
		suite.Contains(parsedPresignedURL.Path, result.ObjectKey)
	}
}

func (suite *hikesTestSuite) TestDeleteHike_Returns401ErrorWhenUnauthorized() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	photo := hike.Photos[0]
	objectKey := s3.GetPhotoObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(suite.router, http.MethodDelete, fmt.Sprintf("/hikes/%d", hike.ID), nil) // Missing session cookie
	suite.Equal(http.StatusUnauthorized, res.Code)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *hikesTestSuite) TestDeleteHike_Returns400ErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodDelete,
		"/hikes/abc",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestDeleteHike_Returns404ErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodDelete,
		"/hikes/1",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesTestSuite) TestDeleteHike_Returns500ErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	mockStore := store.MockStore{DeleteRecordError: errors.New("Something went wrong")}
	hikesHandler := handler.NewHikesHandler(&mockStore, suite.s3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	res := testutils.SendRequest(
		router,
		http.MethodDelete,
		fmt.Sprintf("/hikes/%d", hike.ID),
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesTestSuite) TestDeleteHike_DeletesHikeAndPhotosWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		ListObjectsResult:   &s3Sdk.ListObjectsV2Output{Contents: []s3types.Object{{Key: awsSdk.String("test")}}},
		DeleteObjectsResult: nil,
		DeleteObjectsError:  errors.New("Something went wrong"),
	}
	hikesHandler := handler.NewHikesHandler(suite.store, &mockS3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	photo := hike.Photos[0]
	objectKey := s3.GetPhotoObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(
		router,
		http.MethodDelete,
		fmt.Sprintf("/hikes/%d", hike.ID),
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
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

func (suite *hikesTestSuite) TestDeleteHike_DeletesHikeAndPhotosAndS3Objects() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	photo := hike.Photos[0]
	objectKey := s3.GetPhotoObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(
		suite.router,
		http.MethodDelete,
		fmt.Sprintf("/hikes/%d", hike.ID),
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
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

func (suite *hikesTestSuite) TestListHike_Returns500ErrorWhenDBErrors() {

	mockStore := store.MockStore{
		ListHikesResult: nil,
		ListHikesError:  errors.New("Something went wrong"),
	}
	hikesHandler := handler.NewHikesHandler(&mockStore, suite.s3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

	testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)

	res := testutils.SendRequest(router, http.MethodGet, "/hikes", nil)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesTestSuite) TestListHike_ReturnsNothingWhenThereAreNoHikes() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/hikes", nil)

	suite.Equal(http.StatusOK, res.Code)

	var hikes []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &hikes)
	suite.NoError(err)

	suite.Len(hikes, 0)
}

func (suite *hikesTestSuite) TestListHike_ReturnsHikes() {
	hikes := testutils.ConstructHikes(suite.T(), 2, suite.store, true, true)

	res := testutils.SendRequest(suite.router, http.MethodGet, "/hikes", nil)

	suite.Equal(http.StatusOK, res.Code)

	var response []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 2)
	suite.Equal(hikes[1], response[0])
	suite.Equal(hikes[0], response[1])
}

func (suite *hikesTestSuite) TestUpdateHike_Returns401ErrorWhenUnauthorized() {
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
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/hikes/%d", hike.ID), reqBody) // Missing session cookie

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *hikesTestSuite) TestUpdateHike_Returns400ErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		"/hikes/abc",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestUpdateHike_Returns404ErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		"/hikes/1",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *hikesTestSuite) TestUpdateHike_Returns400ErrorWhenRequestBodyHasInvalidFields() {
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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestUpdateHike_Returns400ErrorWhenRequestBodyIsMissingFields() {
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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *hikesTestSuite) TestUpdateHike_Returns500ErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	mockStore := store.MockStore{
		GetHikeByIDResult: &hike,
		UpdateRecordError: errors.New("Something went wrong"),
	}
	hikesHandler := handler.NewHikesHandler(&mockStore, suite.s3Client)
	router := testutils.NewHikesRouter(hikesHandler, suite.middleware)

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
	res := testutils.SendRequest(
		router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *hikesTestSuite) TestUpdateHike_UpdatesHike() {
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
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/hikes/%d", hike.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusOK, res.Code)

	updated, err := suite.store.GetHikeByID(hike.ID)
	suite.NoError(err)

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
	suite.Equal(expected, *updated)
}
