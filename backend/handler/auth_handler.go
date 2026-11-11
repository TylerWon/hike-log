package handler

import (
	"net/http"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication API routes.
type AuthHandler struct {
	auth auth.Auth
}

// Creates a new AuthHandler.
func NewAuthHandler(auth auth.Auth) (*AuthHandler, error) {
	_, err := registerCustomValidators()
	if err != nil {
		return nil, err
	}
	return &AuthHandler{auth: auth}, nil
}

/*
Checks session cookie to see if a user is authenticated.

Returns: 200 OK and [CheckAuthResponse].
*/
func (h *AuthHandler) CheckAuth(c *gin.Context) {
	ok, _ := h.auth.IsAuthenticated(c)
	c.JSON(http.StatusOK, CheckAuthResponse{Authenticated: ok})
}

/*
Logs in to the admin account.

Request body: [loginRequest]

Returns:
 1. 200 OK and sets a session cookie when successful. The session cookie is used to authenticate requests to restricted
    API routes and is only valid for 24hrs after being issued.
 2. 400 Bad Request and an error message when input is bad
 3. 401 Unauthorized and an error message when authentication fails
*/
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok := h.auth.VerifyPassword(req.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	h.auth.SetSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
