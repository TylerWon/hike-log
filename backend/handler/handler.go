package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// A Handler handles the request-response lifecycle for API routes.
type Handler struct {
	db *gorm.DB
}

// Creates a new Handler with access to the given database.
func New(db *gorm.DB) *Handler {
	_, err := registerCustomValidators()
	if err != nil {
		log.Fatal("Error while registering custom validators: ", err)
	}

	return &Handler{db}
}

// Returns a 200 OK. Used for application health checks.
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

/*
Returns all Hikes in reverse chronological order by Date.

Input: None

Returns:
1. 200 OK when successful
- Response body: list of Hikes

2. 500 Internal Server Error when there is an unexpected error
- Response body: error message
*/
func (h *Handler) ListHike(c *gin.Context) {
	var hikes []models.Hike

	result := h.db.Preload("Photos").Order("date desc, trail_name").Find(&hikes)
	if result.Error != nil {
		log.Println("Failed to list hikes: ", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, hikes)
}

/*
Creates a Hike.

Input: CreateHikeRequest

Returns:
1. 201 Created when successful
- Response body: the Hike

2. 400 Bad Request when input is bad
- Response body: error message

3. 500 Internal Server Error when there is an unexpected error
- Response body: error message
*/
func (h *Handler) CreateHike(c *gin.Context) {
	var req createHikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsed, _ := time.Parse("2006-01-02", req.Date) // convert date string to datatypes.Date

	hike := models.Hike{
		TrailName:     req.TrailName,
		Date:          datatypes.Date(parsed),
		Notes:         req.Notes,
		Rating:        req.Rating,
		Difficulty:    req.Difficulty,
		Distance:      req.Distance,
		ElevationGain: req.ElevationGain,
		Duration:      req.Duration,
		AllTrailsUrl:  req.AllTrailsUrl,
	}

	result := h.db.Create(&hike)
	if result.Error != nil {
		log.Println("Failed to create hike: ", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusCreated, hike)
}
