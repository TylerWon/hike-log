// By adding _test to the package name, this file belongs to its own package separate from handler. This ensures the
// tests can only access the public API exposed by the handler package.
package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Serializes a request body to JSON
func serializeRequestBodyToJSON(t *testing.T, body map[string]any) *bytes.Reader {
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(rawBody)
}

// Sends a request to the provided endpoint and returns the response.
func sendRequest(router *gin.Engine, method string, endpoint string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, endpoint, body)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

// Creates n Hikes, saves them to the database, and returns them.
func createHikes(t *testing.T, db *gorm.DB, n int) []models.Hike {
	var hikes []models.Hike
	for i := range n {
		photos := []models.Photo{
			{SrcUrl: "https://example.com/photo-1.jpg"},
		}
		hike := models.Hike{
			TrailName:     fmt.Sprintf("Trail %d", i),
			Date:          datatypes.Date(time.Date(2026, 1, i, 0, 0, 0, 0, time.UTC)),
			Notes:         "Hike notes",
			Rating:        3,
			Difficulty:    9,
			Distance:      10,
			ElevationGain: 1000,
			Duration:      120,
			AllTrailsUrl:  "https://www.alltrails.com/",
			Photos:        photos,
		}
		hikes = append(hikes, hike)
	}

	result := db.Create(&hikes)
	if result.Error != nil {
		t.Fatal("Failed to create hikes: ", result.Error)
	}

	return hikes
}

type handlerTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

func (suite *handlerTestSuite) SetupTest() {
	suite.db = testutils.SetupTestDB(suite.T())
	s3Client, err := aws.NewS3Client()
	if err != nil {
		suite.T().Fatal("Failed to setup S3 client: ", err)
	}
	handler := handler.New(suite.db, s3Client)
	suite.router = testutils.SetupTestRouter(handler)
}

func (suite *handlerTestSuite) TearDownTest() {
	testutils.TeardownTestDB(suite.T(), suite.db)
}

func (suite *handlerTestSuite) TestListHike_ReturnsNothingWhenThereAreNoHikes() {
	res := sendRequest(suite.router, http.MethodGet, "/api/v1/hikes", nil)

	suite.Equal(http.StatusOK, res.Code)

	var hikes []models.Hike
	err := json.Unmarshal(res.Body.Bytes(), &hikes)
	suite.NoError(err)

	suite.Len(hikes, 0)
}

func (suite *handlerTestSuite) TestListHike_ReturnsHikes() {
	hikes := createHikes(suite.T(), suite.db, 2)

	res := sendRequest(suite.router, http.MethodGet, "/api/v1/hikes", nil)

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
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, "/api/v1/hikes", reqBody)

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
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, "/api/v1/hikes", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
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
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, "/api/v1/hikes", reqBody)

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
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, "/api/v1/hikes/abc/photos/upload-url", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenHikeDoesNotExist() {
	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, "/api/v1/hikes/1/photos/upload-url", reqBody)

	suite.Equal(http.StatusNotFound, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenRequestBodyIsMissingFields() {
	hikes := createHikes(suite.T(), suite.db, 1)

	body := map[string]any{
		"contentLength": 100,
	}
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenContentTypeIsNotAnImageType() {
	hikes := createHikes(suite.T(), suite.db, 1)

	body := map[string]any{
		"contentType":   "text/html",
		"contentLength": 100,
	}
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsErrorWhenContentLengthIsOutsideBounds() {
	hikes := createHikes(suite.T(), suite.db, 1)

	body := map[string]any{
		"contentType":   "text/html",
		"contentLength": -1,
	}
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)

	body = map[string]any{
		"contentType":   "text/html",
		"contentLength": 10485761,
	}
	reqBody = serializeRequestBodyToJSON(suite.T(), body)
	res = sendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestCreatePhotoUploadURL_ReturnsUploadURL() {
	hikes := createHikes(suite.T(), suite.db, 1)

	body := map[string]any{
		"contentType":   "image/jpeg",
		"contentLength": 100,
	}
	reqBody := serializeRequestBodyToJSON(suite.T(), body)
	res := sendRequest(suite.router, http.MethodPost, fmt.Sprintf("/api/v1/hikes/%d/photos/upload-url", hikes[0].ID), reqBody)

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CreatePhotoUploadURLResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.Regexp(fmt.Sprintf("^hikes/%d/[0-9a-f-]{36}$", hikes[0].ID), response.ObjectKey)

	parsedUploadURL, err := url.Parse(response.UploadURL)
	suite.NoError(err)
	suite.Contains(parsedUploadURL.Path, response.ObjectKey)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
