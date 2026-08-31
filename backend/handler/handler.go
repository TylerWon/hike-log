package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// A Handler handles the request-response lifecycle for API routes.
type Handler struct {
	db       *gorm.DB
	s3Client *aws.S3Client
}

// Creates a new Handler
func New(db *gorm.DB, s3Client *aws.S3Client) *Handler {
	_, err := registerCustomValidators()
	if err != nil {
		log.Fatal("Error while registering custom validators: ", err)
	}

	return &Handler{db, s3Client}
}

// Returns a 200 OK. Used for application health checks.
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "ok"})
}

/*
Returns all Hikes in reverse chronological order by Date.

Returns:
 1. 200 OK and a list of [models.Hike] when successful
 2. 500 Internal Server Error and an error message when there is an unexpected error
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

Request body: [createHikeRequest]

Returns:
 1. 201 Created and [models.Hike] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 500 Internal Server Error and an error message when there is an unexpected error
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

/*
Creates a presigned URL that can be used to upload a photo to the S3 bucket.

The request that uses the presigned URL must include the same headers that were provided to generate the URL (i.e.
Content-Type and Content-Length).

Path parameters: [createPhotoUploadURLPathParams]

Request body: [createPhotoUploadURLRequest]

Returns:
 1. 200 OK and [createPhotoUploadURLResponse] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 404 Not Found and an error message when the hike does not exist
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreatePhotoUploadURL(c *gin.Context) {
	var params createPhotoUploadURLPathParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hike models.Hike
	result := h.db.First(&hike, params.HikeID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": http.StatusText(http.StatusNotFound)})
			return
		}
		log.Printf("Failed to get hike (id=%d): %v", params.HikeID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var req createPhotoUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectKey := fmt.Sprintf("hikes/%d/%s", params.HikeID, uuid.New())
	presignedReq, err := h.s3Client.CreatePresignedPutObjectRequest(c, objectKey, req.ContentType, int64(req.ContentLength))

	if err != nil {
		log.Println("Failed to create presigned PutObject request: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	res := CreatePhotoUploadURLResponse{presignedReq.URL, objectKey}
	c.JSON(http.StatusOK, res)
}
