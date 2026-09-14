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
Returns all Hikes and their Photos in reverse chronological order by Date. The Photos are sorted by DisplayOrder.

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

	// Skip validation here since Date is validated when request body gets binded to [createHikeRequest]
	parsed, _ := time.Parse("2006-01-02", req.Date) // convert date string to time.Time

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
Deletes a Hike.

Also deletes any Photos associated with the Hike, and the photo objects stored in the S3 bucket. S3 clean-up is best-
effort and does not result in an error on failure.

Path parameters: [hikeIDPathParam]

Returns:
 1. 200 OK when successful
 2. 404 Not Found and an error message when the Hike does not exist
 3. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) DeleteHike(c *gin.Context) {
	var params hikeIDPathParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.store.GetHikeByID(params.HikeID)
	if err != nil {
		handleHikeDoesNotExistError(c, err, params.HikeID)
		return
	}

	// Note 1: Delete cascades to Photos
	// Note 2: Deletion order matters here. Delete DB models first then S3 objects. This avoids a dangling pointer
	// when S3 objects are deleted first but models fail to delete.
	err = h.store.DeleteModel(&models.Hike{ID: params.HikeID})
	if err != nil {
		log.Printf("Failed to delete Hike (id=%d): %v", params.HikeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	objectKeyPrefix := fmt.Sprintf("hikes/%d/photos/", params.HikeID)
	result, err := h.s3Client.ListObjects(c, objectKeyPrefix)
	if err != nil {
		log.Printf("Failed to list S3 Photo objects for Hike (id=%d): %v", params.HikeID, err)
		c.JSON(http.StatusOK, gin.H{"message": "ok"}) // return OK since Hike technically deleted, only orphan S3 objects
		return
	}

	if len(result.Contents) > 0 {
		var objectKeys []string
		for _, object := range result.Contents {
			objectKeys = append(objectKeys, *object.Key)
		}

		_, err := h.s3Client.DeleteObjects(c, objectKeys)
		if err != nil {
			log.Printf("Failed to delete S3 Photo objects for Hike (id=%d): %v", params.HikeID, err)
			c.JSON(http.StatusOK, gin.H{"message": "ok"}) // return OK since Hike technically deleted, only orphan S3 objects
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

/*
Updates a Hike.

Path parameters: [hikeIDPathParam]

Request body: [updateHikeRequest]

Returns:
 1. 200 OK when successful
 2. 404 Not Found and an error message when the Hike does not exist
 3. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) UpdateHike(c *gin.Context) {
	var params hikeIDPathParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.store.GetHikeByID(params.HikeID)
	if err != nil {
		handleHikeDoesNotExistError(c, err, params.HikeID)
		return
	}

	var req updateHikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Skip validation here since Date is validated when request body gets binded to [updateHikeRequest]
	parsed, _ := time.Parse("2006-01-02", req.Date) // convert date string to time.Time

	hike := models.Hike{
		ID:            params.HikeID,
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
	err = h.store.UpdateModel(&hike)
	if err != nil {
		log.Println("Failed to update Hike: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

/*
Creates a presigned URL that can be used to upload a Photo for a Hike to the S3 bucket.

The request that uses the presigned URL must include the same headers that were provided to generate the URL (i.e.
Content-Type and Content-Length).

Path parameters: [hikeIDPathParam]

Request body: [createPhotoUploadURLRequest]

Returns:
 1. 200 OK and [createPhotoUploadURLResponse] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 404 Not Found and an error message when the Photo does not exist
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreatePhotoUploadURL(c *gin.Context) {
	var params hikeIDPathParam
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

	objectKey := h.s3Client.CreatePhotoObjectKey(params.HikeID)
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
 3. 404 Not Found and an error message when the Photo does not exist
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreatePhoto(c *gin.Context) {
	var params hikeIDPathParam
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
		SrcUrl:       srcURL,
		Caption:      req.Caption,
		DisplayOrder: req.DisplayOrder,
		HikeID:       params.HikeID,
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
