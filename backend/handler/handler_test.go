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
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestListHike_Empty(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)
	router := testutils.SetupTestRouter(handler.New(db))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var hikes []models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &hikes)
	if err != nil {
		t.Fatal("Failed to decode response: ", err)
	}

	assert.Len(t, hikes, 0)
}

func TestListHike_ReturnsHikes(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)
	router := testutils.SetupTestRouter(handler.New(db))

	var hikes []*models.Hike
	hike1 := models.Hike{
		TrailName:     "Trail 1",
		Date:          datatypes.Date(time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC)),
		Rating:        4.5,
		Difficulty:    3,
		Distance:      8.2,
		ElevationGain: 1200,
		TotalTime:     60,
	}
	hikes = append(hikes, &hike1)

	notes := "Difficult hike"
	allTrailsUrl := "https://www.alltrails.com/"
	photos := []models.Photo{
		{SrcUrl: "https://example.com/photo-1.jpg"},
	}
	hike2 := models.Hike{
		TrailName:     "Trail 2",
		Date:          datatypes.Date(time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)),
		Notes:         &notes,
		Rating:        3,
		Difficulty:    9.5,
		Distance:      21.2,
		ElevationGain: 1587,
		TotalTime:     127,
		AllTrailsUrl:  &allTrailsUrl,
		Photos:        photos,
		CoverPhoto:    &photos[0],
	}
	hikes = append(hikes, &hike2)

	result := db.Create(&hikes)
	if result.Error != nil {
		t.Fatal("Failed to create hikes: ", result.Error)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to decode response: ", err)
	}

	assert.Len(t, response, 2)
	assert.Equal(t, hike1, response[0])
	hike2.Photos = nil // Set Photos field to nil since it's excluded from the response by the List endpoint
	assert.Equal(t, hike2, response[1])
}

func TestRetrieveHike_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)
	router := testutils.SetupTestRouter(handler.New(db))

	photoCaption := "The peak"
	photos := []models.Photo{
		{SrcUrl: "https://example.com/photo-1.jpg"},
		{SrcUrl: "https://example.com/photo-2.jpg", Caption: &photoCaption},
	}
	hike := models.Hike{
		TrailName:     "Trail 1",
		Date:          datatypes.Date(time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)),
		Rating:        3,
		Difficulty:    9.5,
		Distance:      21.2,
		ElevationGain: 1587,
		TotalTime:     60,
		Photos:        photos,
		CoverPhoto:    &photos[0],
	}

	result := db.Create(&hike)
	if result.Error != nil {
		t.Fatal("Failed to create hike: ", result.Error)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Hike
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Failed to decode response: ", err)
	}

	hike.CoverPhoto = nil // Set CoverPhoto field to nil since it's excluded from the response by the Retrieve endpoint
	assert.Equal(t, hike, response)
}

func TestRetrieveHike_NotFound(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)
	router := testutils.SetupTestRouter(handler.New(db))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hikes/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
