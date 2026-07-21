package handler_test

import (
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

type HandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.db = testutils.SetupTestDB(suite.T())
	suite.router = testutils.SetupTestRouter(handler.New(suite.db))
}

func (suite *HandlerTestSuite) TearDownTest() {
	testutils.TeardownTestDB(suite.T(), suite.db)
}

func (suite *HandlerTestSuite) TestListHike_Empty() {
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

func (suite *HandlerTestSuite) TestListHike_ReturnsHikes() {
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
		Photos: 	   []models.Photo{},
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

func (suite *HandlerTestSuite) TestRetrieveHike_Success() {
	photos := []models.Photo{
		{SrcUrl: "https://example.com/photo-1.jpg"},
		{SrcUrl: "https://example.com/photo-2.jpg", Caption: "The peak"},
	}
	hike := models.Hike{
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

	result := suite.db.Create(&hike)
	if result.Error != nil {
		suite.T().Fatal("Failed to create hike: ", result.Error)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes/1", nil)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusOK, rec.Code)

	var response models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		suite.T().Fatal("Failed to decode response: ", err)
	}

	suite.Equal(hike, response)
}

func (suite *HandlerTestSuite) TestRetrieveHike_NotFound() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes/999", nil)
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
