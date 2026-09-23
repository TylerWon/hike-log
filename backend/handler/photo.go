package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/TylerWon/hike-log/backend/models"
	"github.com/gin-gonic/gin"
)

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
