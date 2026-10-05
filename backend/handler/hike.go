package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/TylerWon/hike-log/backend/models/types"
	"github.com/TylerWon/hike-log/backend/s3"
	"github.com/gin-gonic/gin"
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

	hike := models.Hike{
		TrailName:     req.TrailName,
		Date:          types.Date(req.Date),
		Notes:         req.Notes,
		Rating:        req.Rating,
		Difficulty:    req.Difficulty,
		Distance:      req.Distance,
		ElevationGain: req.ElevationGain,
		Duration:      req.Duration,
		AllTrailsUrl:  req.AllTrailsUrl,
		Photos:        []models.Photo{},
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
Creates Photos for a Hike.

The photos should already be uploaded to S3. This endpoint is just responsible for creating the Photo model in the DB.

Path parameters: [hikeIDPathParam]

Request body: [createPhotosRequest]

Returns:
 1. 201 Created and [CreatePhotosResponse] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 404 Not Found and an error message when the Hike does not exist
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreatePhotos(c *gin.Context) {
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

	var req createPhotosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Note: ShouldBindJSON unmarshals and validates the entire request array in one shot so the handler returns early
		// if ANY of the items in the array are invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res CreatePhotosResponse
	for idx, item := range req {
		keyParts := strings.Split(item.ObjectKey, "/")
		hikeId, _ := strconv.ParseUint(keyParts[1], 10, 64)
		if hikeId != uint64(params.HikeID) {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": fmt.Sprintf("Hike ID in path does not match with the ObjectKey of item %d", idx)},
			)
			return
		}

		exists, err := h.s3Client.DoesObjectExist(c, item.ObjectKey)
		if err != nil {
			log.Printf("Unable to verify if photo %d (ObjectKey=%s) exists in the S3 bucket: %v", idx, item.ObjectKey, err)
			item := createPhotosResponseItem{
				Index:   *item.Index,
				Success: false,
				Error:   fmt.Sprintf("Unable to verify if a photo with the ObjectKey exists in the S3 bucket: %v", err),
			}
			res = append(res, item)
			continue
		} else if !exists {
			item := createPhotosResponseItem{
				Index:   *item.Index,
				Success: false,
				Error:   "Photo with the ObjectKey does not exist in the S3 bucket",
			}
			res = append(res, item)
			continue
		}

		srcURL := s3.CreatePhotoObjectURL(item.ObjectKey, os.Getenv("ENV"))
		photo := models.Photo{
			SrcUrl:       srcURL,
			Caption:      item.Caption,
			DisplayOrder: item.DisplayOrder,
			HikeID:       params.HikeID,
		}
		err = h.store.CreateRecord(&photo)
		if err != nil {
			log.Printf("Failed to create Photo %d: %v", idx, err)
			item := createPhotosResponseItem{Index: *item.Index, Success: false, Error: "Failed to create photo"}
			res = append(res, item)
			continue
		}

		res = append(res, createPhotosResponseItem{Index: *item.Index, Success: true, Result: &photo})
	}

	c.JSON(http.StatusCreated, res)
}

/*
Creates presigned URLs that can be used to upload Photos for a Hike to the S3 bucket.

The request that uses the presigned URL must include the same headers that were provided to generate the URL (i.e.
Content-Type and Content-Length).

Path parameters: [hikeIDPathParam]

Request body: [createPresignedURLsRequest]

Returns:
 1. 201 Created and [CreatePresignedURLsResponse] when successful
 2. 400 Bad Request and an error message when input is bad
 3. 404 Not Found and an error message when the Hike does not exist
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CreatePresignedURLs(c *gin.Context) {
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

	var req createPresignedURLsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Note: ShouldBindJSON unmarshals and validates the entire request array in one shot so the handler returns early
		// if ANY of the items in the array are invalid
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var res CreatePresignedURLsResponse
	for idx, item := range req {
		objectKey := s3.CreatePhotoObjectKey(params.HikeID)

		presignedReq, err := h.s3Client.CreatePresignedPutObjectRequest(
			c,
			objectKey,
			item.ContentType,
			int64(item.ContentLength),
		)
		if err != nil {
			log.Printf("Failed to create presigned URL for item %d: %v", idx, err)
			item := createPresignedURLsResponseItem{Index: *item.Index, Success: false, Error: "Failed to create presigned URL"}
			res = append(res, item)
			continue
		}

		item := createPresignedURLsResponseItem{
			Index:   *item.Index,
			Success: true,
			Result: &createPresignedURLsResponseItemResult{
				PresignedURL: presignedReq.URL,
				ObjectKey:    objectKey,
			},
		}
		res = append(res, item)
	}

	c.JSON(http.StatusCreated, res)
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

	hike := models.Hike{
		ID:            params.HikeID,
		TrailName:     req.TrailName,
		Date:          types.Date(req.Date),
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
