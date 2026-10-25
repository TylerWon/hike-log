package handler

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*
Checks session cookie to see if a user is authenticated.

Returns:
 1. 200 OK and [CheckAuthResponse] when successful.
 2. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CheckAuth(c *gin.Context) {
	authenticated, err := auth.IsAuthenticated(c)

	// Return 500 error if authentication fails for non-auth reason
	if err != nil {
		_, ok := errors.AsType[auth.AuthError](err)
		if !ok {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}
	}

	c.JSON(http.StatusOK, CheckAuthResponse{Authenticated: authenticated})
}

/*
Logs in to the admin account.

Request body: [loginRequest]

Returns:
 1. 200 OK and sets a session cookie when successful. The session cookie is used to authenticate requests to restricted
    API routes and is only valid for 24hrs after being issued.
 2. 400 Bad Request and an error message when input is bad
 3. 401 Unauthorized and an error message when authentication fails
 4. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := os.Getenv("ADMIN_PASSWORD_HASH")
	if hashedPassword == "" {
		log.Println("Admin password hash could not be retrieved from ADMIN_PASSWORD_HASH env var")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	err = auth.SetSessionCookie(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
