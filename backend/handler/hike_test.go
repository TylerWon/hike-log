package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/models/types"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/TylerWon/hike-log/backend/testutils"
	awsSdk "github.com/aws/aws-sdk-go-v2/aws"
	s3Sdk "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

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
	mockStore := store.MockStore{CreateRecordError: errors.New("Something went wrong")}
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

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsErrorWhenHikeIDIsInvalid() {
	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    s3.CreatePhotoObjectKey(1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/abc/photos", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsErrorWhenHikeDoesNotExist() {
	body := [](map[string]any){
		{
			"index":        0,
			"objectKey":    s3.CreatePhotoObjectKey(1),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/1/photos", reqBody)

	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsErrorWhenRequestItemIsMissingFields() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsErrorWhenObjectKeyIsInvalid() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsErrorWhenPathAndObjectKeyHikeIDMismatch() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenPhotoDoesNotExistInS3() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 1)
	suite.Equal(uint(0), response[0].Index)
	suite.False(response[0].Success)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		DoesObjectExistResult: false,
		DoesObjectExistError:  errors.New("Something went wrong"),
	}
	h := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), h)

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
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 1)
	suite.Equal(uint(2), response[0].Index)
	suite.False(response[0].Success)
}

func (suite *handlerTestSuite) TestCreatePhotos_ReturnsPartialErrorWhenDBErrors() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	mockStore := store.MockStore{
		GetHikeByIDResult: &hikes[0],
		CreateRecordError: errors.New("Something went wrong"),
	}
	mockS3Client := s3.MockS3Client{DoesObjectExistResult: true}
	h := handler.New(&mockStore, &mockS3Client)
	router := testutils.NewRouter(suite.T(), h)

	body := [](map[string]any){
		{
			"index":        4,
			"objectKey":    s3.CreatePhotoObjectKey(hikes[0].ID),
			"caption":      "Caption",
			"displayOrder": 1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

	suite.Equal(http.StatusCreated, res.Code)

	var response handler.CreatePhotosResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Len(response, 1)
	suite.Equal(uint(4), response[0].Index)
	suite.False(response[0].Success)
}

func (suite *handlerTestSuite) TestCreatePhotos_CreatesPhotos() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos", hikes[0].ID), reqBody)

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

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsErrorWhenHikeIDIsInvalid() {
	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/abc/photos/presigned-urls", reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsErrorWhenHikeDoesNotExist() {
	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/hikes/1/photos/presigned-urls", reqBody)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsErrorWhenRequestItemIsMissingFields() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsErrorWhenRequestItemHasAnInvalidImageType() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "text/html",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsErrorWhenContentLengthIsOutsideBounds() {
	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": -1,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)

	body = [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 10485761,
		},
	}
	reqBody = testutils.SerializeJSONRequestBody(suite.T(), body)
	res = testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsErrorWhenDBErrors() {
	mockStore := store.MockStore{
		GetHikeByIDResult: nil,
		GetHikeByIDError:  errors.New("Something went wrong"),
	}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         0,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsPartialErrorWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		CreatePresignedPutObjectRequestResult: nil,
		CreatePresignedPutObjectRequestError:  errors.New("Something went wrong"),
	}
	h := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), h)

	hikes := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)

	body := [](map[string]any){
		{
			"index":         2,
			"contentType":   "image/jpeg",
			"contentLength": 100,
		},
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)

	suite.Equal(http.StatusOK, res.Code)

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

func (suite *handlerTestSuite) TestCreatePresignedURLs_ReturnsPresignedURLs() {
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
	res := testutils.SendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/presigned-urls", hikes[0].ID), reqBody)

	suite.Equal(http.StatusOK, res.Code)

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

func (suite *handlerTestSuite) TestDeleteHike_ReturnsErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodDelete, "/api/v1/hikes/abc", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestDeleteHike_ReturnsErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(suite.router, http.MethodDelete, "/api/v1/hikes/1", nil)
	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestDeleteHike_ReturnsErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, true, true)[0]

	mockStore := store.MockStore{DeleteRecordError: errors.New("Something went wrong")}
	handler := handler.New(&mockStore, suite.s3Client)
	router := testutils.NewRouter(suite.T(), handler)

	res := testutils.SendRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), nil)
	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestDeleteHike_DeletesHikeAndPhotosWhenS3Errors() {
	mockS3Client := s3.MockS3Client{
		ListObjectsResult:   &s3Sdk.ListObjectsV2Output{Contents: []s3types.Object{{Key: awsSdk.String("test")}}},
		DeleteObjectsResult: nil,
		DeleteObjectsError:  errors.New("Something went wrong"),
	}
	handler := handler.New(suite.store, &mockS3Client)
	router := testutils.NewRouter(suite.T(), handler)

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

	res := testutils.SendRequest(router, http.MethodDelete, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), nil)
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
	objectKey := s3.GetPhotoObjectKey(photo.SrcUrl, "local")
	_, err := suite.s3Client.PutObject(
		context.TODO(),
		objectKey,
		strings.NewReader("content"),
		"image/png",
	)
	suite.NoError(err)

	res := testutils.SendRequest(suite.router, http.MethodDelete, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), nil)
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

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenHikeIDIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodPut, "/api/v1/hikes/abc", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenHikeDoesNotExist() {
	res := testutils.SendRequest(suite.router, http.MethodPut, "/api/v1/hikes/1", nil)
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
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), reqBody)

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
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestUpdateHike_ReturnsErrorWhenDBErrors() {
	hike := testutils.ConstructHikes(suite.T(), 1, suite.store, false, true)[0]

	mockStore := store.MockStore{
		GetHikeByIDResult: &hike,
		UpdateRecordError: errors.New("Something went wrong"),
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
	res := testutils.SendRequest(router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), reqBody)

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
	res := testutils.SendRequest(suite.router, http.MethodPut, fmt.Sprintf("/api/v1/hikes/%d", hike.ID), reqBody)

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
