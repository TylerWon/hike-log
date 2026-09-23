package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
)

func (suite *handlerTestSuite) TestDeletePhoto_ReturnsErrorWhenPhotoIDIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodDelete, "/api/v1/photos/abc/", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestDeletePhoto_ReturnsErrorWhenPhotoDoesNotExist() {
	res := testutils.SendRequest(suite.router, http.MethodDelete, "/api/v1/photos/1/", nil)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestDeletePhoto_ReturnsErrorWhenDBErrors() {
	mockStore := store.MockStore{GetPhotoByIDError: errors.New("Something went wrong")}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	res := testutils.SendRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), nil)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestDeletePhoto_DeletesPhotoWhenS3Errors() {
	mockS3Client := aws.MockS3Client{
		DeleteObjectResult: nil,
		DeleteObjectError:  errors.New("Something went wrong"),
	}
	handler := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), handler)

	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]
	objectKey := suite.s3Client.GetObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), nil)
	suite.Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.True(exists)

	suite.s3Client.DeleteObject(context.TODO(), objectKey)
}

func (suite *handlerTestSuite) TestDeletePhoto_DeletesPhotoAndS3Object() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]
	objectKey := suite.s3Client.GetObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(suite.router, http.MethodDelete, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), nil)
	suite.Equal(http.StatusOK, res.Code)

	_, err = suite.store.GetPhotoByID(photo.ID)
	suite.Error(err)

	exists, err := suite.s3Client.DoesObjectExist(context.TODO(), objectKey)
	suite.NoError(err)
	suite.False(exists)
}

func (suite *handlerTestSuite) TestUpdatePhoto_ReturnsErrorWhenPhotoIDIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodPut, "/api/v1/photos/abc/", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdatePhoto_ReturnsErrorWhenPhotoDoesNotExist() {
	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, "/api/v1/photos/1/", reqBody)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestUpdatePhoto_ReturnsErrorWhenRequestBodyHasInvalidFields() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 0, // must be > 0
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdatePhoto_ReturnsErrorWhenRequestBodyIsMissingFields() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption": "Updated caption",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdatePhoto_ReturnsErrorWhenDBErrors() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	mockStore := store.MockStore{
		GetPhotoByIDResult: &photo,
		UpdateRecordError:  errors.New("Something went wrong"),
	}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestUpdatePhoto_UpdatesPhoto() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"caption":      "Updated caption",
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), reqBody)

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

func (suite *handlerTestSuite) TestUpdatePhoto_UpdatesPhotoWhenOptionalFieldMissing() {
	photo := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0].Photos[0]

	body := map[string]any{
		"displayOrder": 2,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/photos/%d/", photo.ID), reqBody)

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
