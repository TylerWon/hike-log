package handler

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/gin-gonic/gin"
)

/*
Creates a Photo for a Hike.

The photo should already be uploaded to S3. This endpoint is just responsible for creating the Photo model in the DB.

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
		handleRecordDoesNotExistError(c, err, params.HikeID)
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
		log.Printf("Unable to verify if photo object (key=%s) exists in S3 bucket: %v", req.ObjectKey, err)
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
	err = h.store.CreateRecord(&photo)
	if err != nil {
		log.Println("Failed to create Photo: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusCreated, photo)
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
Deletes a Photo.

Also deletes the object stored in the S3 bucket. S3 clean-up is best-effort and does not result in an error on failure.

Path parameters: [photoIDPathParam]

Returns:
 1. 200 OK when successful
 2. 404 Not Found and an error message when the Photo does not exist
 3. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) DeletePhoto(c *gin.Context) {
	var params photoIDPathParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo, err := h.store.GetPhotoByID(params.PhotoID)
	if err != nil {
		handleRecordDoesNotExistError(c, err, params.PhotoID)
		return
	}

	err = h.store.DeleteRecord(&models.Photo{ID: params.PhotoID})
	if err != nil {
		log.Printf("Failed to delete Photo (id=%d): %v", params.PhotoID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	objectKey := h.s3Client.GetObjectKey(photo.SrcUrl, os.Getenv("ENV"))
	_, err = h.s3Client.DeleteObject(c, objectKey)
	if err != nil {
		log.Printf("Failed to delete S3 photo object for Photo (id=%d): %v", params.PhotoID, err)
		c.JSON(http.StatusOK, gin.H{"message": "ok"}) // return OK since Photo deleted, only orphan S3 object
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

/*
Updates a Photo.

Path parameters: [photoIDPathParam]

Request body: [updatePhotoRequest]

Returns:
 1. 200 OK when successful
 2. 404 Not Found and an error message when the Photo does not exist
 3. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) UpdatePhoto(c *gin.Context) {
	var params photoIDPathParam
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo, err := h.store.GetPhotoByID(params.PhotoID)
	if err != nil {
		handleRecordDoesNotExistError(c, err, params.PhotoID)
		return
	}

	var req updatePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.DisplayOrder = req.DisplayOrder
	if req.Caption != "" {
		photo.Caption = req.Caption
	}

	err = h.store.UpdateRecord(photo)
	if err != nil {
		log.Printf("Failed to update Photo (id=%d): %v", photo.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
