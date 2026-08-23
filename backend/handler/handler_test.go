// By adding _test to the package name, this file belongs to its own package separate from handler. This ensures the
// tests can only access the public API exposed by the handler package.
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Serializes a JSON request body
func serializeJsonRequestBody(t *testing.T, body map[string]any) *bytes.Reader {
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(rawBody)
}

type handlerTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

func (suite *handlerTestSuite) SetupTest() {
	suite.db = testutils.SetupTestDB(suite.T())
	suite.router = testutils.SetupTestRouter(handler.New(suite.db))
}

func (suite *handlerTestSuite) TearDownTest() {
	testutils.TeardownTestDB(suite.T(), suite.db)
}

func (suite *handlerTestSuite) TestListHike_ReturnsNothingWhenThereAreNoHikes() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes", nil)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusOK, rec.Code)

	var hikes []models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &hikes)
	if err != nil {
		suite.T().Fatal("Failed to decode response: ", err)
	}

	suite.Len(hikes, 0)
}

func (suite *handlerTestSuite) TestListHike_ReturnsHikes() {
	var hikes []*models.Hike
	hike1 := models.Hike{
		TrailName:     "Trail 1",
		Date:          datatypes.Date(time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC)),
		Notes:         "Easy hike",
		Rating:        4.5,
		Difficulty:    3,
		Distance:      8.2,
		ElevationGain: 1200,
		AllTrailsUrl:  "https://www.alltrails.com/",
		Duration:      60,
		Photos:        []models.Photo{},
	}
	hikes = append(hikes, &hike1)

	photos := []models.Photo{
		{SrcUrl: "https://example.com/photo-1.jpg"},
	}
	hike2 := models.Hike{
		TrailName:     "Trail 2",
		Date:          datatypes.Date(time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)),
		Notes:         "Difficult hike",
		Rating:        3,
		Difficulty:    9.5,
		Distance:      21.2,
		ElevationGain: 1587,
		Duration:      127,
		AllTrailsUrl:  "https://www.alltrails.com/",
		Photos:        photos,
	}
	hikes = append(hikes, &hike2)

	result := suite.db.Create(&hikes)
	if result.Error != nil {
		suite.T().Fatal("Failed to create hikes: ", result.Error)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes", nil)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusOK, rec.Code)

	var response []models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatal("Failed to decode response: ", err)
	}

	suite.Len(response, 2)
	suite.Equal(hike1, response[0])
	suite.Equal(hike2, response[1])
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
	reqBody := serializeJsonRequestBody(suite.T(), body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hikes", reqBody)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
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
	reqBody := serializeJsonRequestBody(suite.T(), body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hikes", reqBody)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusBadRequest, rec.Code)
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
	reqBody := serializeJsonRequestBody(suite.T(), body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/hikes", reqBody)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusCreated, rec.Code)

	var response models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatal("Failed to decode response: ", err)
	}

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

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}
