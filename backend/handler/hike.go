package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

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
	err := h.store.CreateRecord(&hike)
	if err != nil {
		log.Println("Failed to create Hike: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusCreated, hike)
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
		handleRecordDoesNotExistError(c, err, params.HikeID)
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
Deletes a Hike.

Also deletes any Photos associated with the Hike and the photo objects stored in the S3 bucket. S3 clean-up is best-
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
		handleRecordDoesNotExistError(c, err, params.HikeID)
		return
	}

	// Note 1: Delete cascades to Photos
	// Note 2: Deletion order matters here. Delete DB models first then S3 objects. This avoids a dangling pointer
	// when S3 objects are deleted first but models fail to delete.
	err = h.store.DeleteRecord(&models.Hike{ID: params.HikeID})
	if err != nil {
		log.Printf("Failed to delete Hike (id=%d): %v", params.HikeID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	objectKeyPrefix := fmt.Sprintf("hikes/%d/photos/", params.HikeID)
	result, err := h.s3Client.ListObjects(c, objectKeyPrefix)
	if err != nil {
		log.Printf("Failed to list S3 Photo objects for Hike (id=%d): %v", params.HikeID, err)
		c.JSON(http.StatusOK, gin.H{"message": "ok"}) // return OK since Hike deleted, only orphan S3 objects
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
			c.JSON(http.StatusOK, gin.H{"message": "ok"}) // return OK since Hike deleted, only orphan S3 objects
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
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
		handleRecordDoesNotExistError(c, err, params.HikeID)
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
	err = h.store.UpdateRecord(&hike)
	if err != nil {
		log.Printf("Failed to update Hike (id=%d): %v", hike.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
