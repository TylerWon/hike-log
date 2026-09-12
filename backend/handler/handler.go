package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TylerWon/hike-log/backend/aws"
	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// A Handler handles the request-response lifecycle for API routes.
type Handler struct {
	store    store.Store
	s3Client aws.S3Client
}

// Creates a new Handler
func New(store store.Store, s3Client aws.S3Client) *Handler {
	_, err := registerCustomValidators()
	if err != nil {
		log.Fatal("Error while registering custom validators: ", err)
	}

	return &Handler{store, s3Client}
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
func (h *Handler) ListHikes(c *gin.Context) {
	hikes, err := h.store.ListHikes()
	if err != nil {
		log.Println("Failed to list hikes: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, hikes)
}

/*
Creates a Hike.

Request body: [createHikeRequest]

Returns:
 1. 201 Created and the [models.Hike] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreateHike(c *gin.Context) {
	var req createHikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Skip validation here since Data is validated when request body gets binded to [createHikeRequest]
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
	err := h.store.CreateModel(&hike)
	if err != nil {
		log.Println("Failed to create Hike: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusCreated, hike)
}

/*
Creates a presigned URL that can be used to upload a photo for a Hike to the S3 bucket.

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

	_, err := h.store.GetHikeByID(params.HikeID)
	if err != nil {
		handleHikeDoesNotExistError(c, err, params.HikeID)
		return
	}

	var req createPhotoUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectKey := fmt.Sprintf("hikes/%d/photos/%s", params.HikeID, uuid.New())
	presignedReq, err := h.s3Client.CreatePresignedPutObjectRequest(c, objectKey, req.ContentType, int64(req.ContentLength))

	if err != nil {
		log.Println("Failed to create presigned PutObject request: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	res := CreatePhotoUploadURLResponse{presignedReq.URL, objectKey}
	c.JSON(http.StatusOK, res)
}

/*
Creates a Photo for a Hike.

Path parameters: [createPhotoPathParams]

Request body: [createPhotoRequest]

Returns:
 1. 201 Created and the [models.Photo] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 404 Not Found and an error message when the hike does not exist
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreatePhoto(c *gin.Context) {
	var params createPhotoPathParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.store.GetHikeByID(params.HikeID)
	if err != nil {
		handleHikeDoesNotExistError(c, err, params.HikeID)
		return
	}

	var req createPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Skip validation here since ObjectKey is validated when request body gets binded to [createPhotoRequest]
	keyParts := strings.Split(req.ObjectKey, "/")
	hikeId, _ := strconv.ParseUint(keyParts[1], 10, 64)
	if hikeId != uint64(params.HikeID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hike ID in path and ObjectKey do not match"})
		return
	}

	exists, err := h.s3Client.DoesObjectExist(c, req.ObjectKey)
	if err != nil {
		log.Printf("Unable to verify if Photo (key=%s) exists in S3 bucket: %v", req.ObjectKey, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	} else if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo with provided ObjectKey does not exist in S3 bucket"})
		return
	}

	srcURL := h.s3Client.GetObjectURL(req.ObjectKey, os.Getenv("ENV"))
	photo := models.Photo{
		SrcUrl:  srcURL,
		Caption: req.Caption,
		HikeID:  params.HikeID,
	}
	err = h.store.CreateModel(&photo)
	if err != nil {
		log.Println("Failed to create Photo: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

// Handles error that occurs when a Hike unexpectedly does not exist.
func handleHikeDoesNotExistError(c *gin.Context, err error, hikeId uint) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	log.Printf("Failed to get hike (id=%d): %v", hikeId, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
}
