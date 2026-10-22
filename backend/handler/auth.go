package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	ONE_DAY_IN_SECONDS  = 24 * 60 * 60
	SESSION_COOKIE_NAME = "hike-log-session"
)

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

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Println("Session secret could not be retrieved from SESSION_SECRET env var")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	setSessionCookie(c, sessionSecret)

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

/*
Sets a session cookie on the response.

The session cookie is used to authenticate requests to restricted API routes and is only valid for 24hrs after being
issued.

The cookie is named [SESSION_COOKIE_NAME] and its value is "<cookie expiry time>.<HMAC-SHA256 hex signature of the
expiry time>". The signature is included so that the API can verify the authenticity of the cookie when it is provided
to restricted routes.
*/
func setSessionCookie(c *gin.Context, sessionSecret string) {
	expiry := time.Now().Add(24 * time.Hour).Unix()
	expiryStr := strconv.FormatInt(expiry, 10)

	mac := hmac.New(sha256.New, []byte(sessionSecret))
	mac.Write([]byte(expiryStr))
	digest := mac.Sum(nil)
	signature := hex.EncodeToString(digest)

	cookieValue := expiryStr + "." + signature

	// Only allow cookie to be sent between same sites (same site = scheme + host are the same)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SESSION_COOKIE_NAME,
		cookieValue,
		ONE_DAY_IN_SECONDS,
		// Cookie can be sent to any path
		"/",
		// Cookie can only be sent to same host that issued it
		"",
		// Cookie must be sent over HTTPS
		true,
		// Cookie cannot be accessed by the browser
		true,
	)
}
