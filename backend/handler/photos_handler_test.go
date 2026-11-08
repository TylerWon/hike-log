package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type photosHandlerTestSuite struct {
	suite.Suite
	auth       auth.Auth
	middleware *middleware.Middleware
	s3Client   s3.S3Client
	store      store.Store
	suiteDB    *testutils.TestSuiteDB
	router     *gin.Engine
}

func (suite *photosHandlerTestSuite) SetupSuite() {
	suite.suiteDB = testutils.NewTestSuiteDB(suite.T())
	suite.store = testutils.NewStore(suite.T(), suite.suiteDB)
	suite.s3Client = testutils.NewS3Client(suite.T())
	suite.auth = testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, testutils.SESSION_SECRET)
	suite.middleware = middleware.New(suite.auth)
	suite.router = testutils.NewPhotosRouter(handler.NewPhotosHandler(suite.store, suite.s3Client), suite.middleware)
}

func (suite *photosHandlerTestSuite) TearDownSuite() {
	testutils.TeardownStore(suite.T(), suite.store)
	suite.suiteDB.Teardown(suite.T())
}

func (suite *photosHandlerTestSuite) SetupTest() {
	suite.suiteDB.Reset(suite.T())
}

func TestPhotosTestSuite(t *testing.T) {
	suite.Run(t, new(photosHandlerTestSuite))
}

func (suite *photosHandlerTestSuite) TestDeletePhoto_Returns401ErrorWhenUnauthorized() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]
	objectKey := s3.GetPhotoObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(suite.router, http.MethodDelete, fmt.Sprintf("/photos/%d", photo.ID), nil) // Missing session cookie
	suite.Equal(http.StatusUnauthorized, res.Code)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *photosHandlerTestSuite) TestDeletePhoto_Returns400ErrorWhenPhotoIDIsInvalid() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodDelete,
		"/photos/abc",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *photosHandlerTestSuite) TestDeletePhoto_Returns404ErrorWhenPhotoDoesNotExist() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodDelete,
		"/photos/1",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *photosHandlerTestSuite) TestDeletePhoto_Returns500ErrorWhenDBErrors() {
	mockStore := store.MockStore{GetPhotoByIDError: errors.New("Something went wrong")}
	photosHandler := handler.NewPhotosHandler(&mockStore, suite.s3Client)
	router := testutils.NewPhotosRouter(photosHandler, suite.middleware)

	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	res := testutils.SendRequest(
		router,
		http.MethodDelete,
		fmt.Sprintf("/photos/%d", photo.ID),
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *photosHandlerTestSuite) TestDeletePhoto_DeletesPhotoWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		DeleteObjectResult: nil,
		DeleteObjectError:  errors.New("Something went wrong"),
	}
	photosHandler := handler.NewPhotosHandler(suite.store, &mockS3Client)
	router := testutils.NewPhotosRouter(photosHandler, suite.middleware)

	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]
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
		fmt.Sprintf("/photos/%d", photo.ID),
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *photosHandlerTestSuite) TestDeletePhoto_DeletesPhotoAndS3Object() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]
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
		fmt.Sprintf("/photos/%d", photo.ID),
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.False(exists)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_Returns401ErrorWhenUnauthorized() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/photos/%d", photo.ID), reqBody) // Missing session cookie

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_Returns400ErrorWhenPhotoIDIsInvalid() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		"/photos/abc",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_Returns404ErrorWhenPhotoDoesNotExist() {
	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		"/photos/1",
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_Returns400ErrorWhenRequestBodyHasInvalidFields() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 0, // must be > 0
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/photos/%d", photo.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_Returns400ErrorWhenRequestBodyIsMissingFields() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption": "Updated caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/photos/%d", photo.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_Returns500ErrorWhenDBErrors() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	mockStore := store.MockStore{
		GetPhotoByIDResult: &photo,
		UpdateRecordError:  errors.New("Something went wrong"),
	}
	photosHandler := handler.NewPhotosHandler(&mockStore, suite.s3Client)
	router := testutils.NewPhotosRouter(photosHandler, suite.middleware)

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		router,
		http.MethodPut,
		fmt.Sprintf("/photos/%d", photo.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_UpdatesPhoto() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/photos/%d", photo.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusOK, res.Code)

	updated, err := suite.store.GetPhotoByID(photo.ID)
	suite.NoError(err)

	expected := models.Photo{
		ID:           photo.ID,
		SrcUrl:       photo.SrcUrl,
		Caption:      "Updated caption",
		DisplayOrder: 2,
		HikeID:       photo.HikeID,
	}
	suite.Equal(expected, *updated)
}

func (suite *photosHandlerTestSuite) TestUpdatePhoto_UpdatesPhotoWhenOptionalFieldMissing() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(
		suite.router,
		http.MethodPut,
		fmt.Sprintf("/photos/%d", photo.ID),
		reqBody,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusOK, res.Code)

	updated, err := suite.store.GetPhotoByID(photo.ID)
	suite.NoError(err)

	expected := models.Photo{
		ID:           photo.ID,
		SrcUrl:       photo.SrcUrl,
		Caption:      photo.Caption,
		DisplayOrder: 2,
		HikeID:       photo.HikeID,
	}
	suite.Equal(expected, *updated)
}
