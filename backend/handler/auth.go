package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	ONE_DAY_IN_SECONDS  = 24 * 60 * 60
	SESSION_COOKIE_NAME = "hike-log-session"
)

/*
Checks if a user is authenticated with the admin account.

A user is authenticated if they have a session cookie and its value satisfies the following criteria:
1. The expiry time in the value is valid, and
2. A HMAC-SHA256 signature of the expiry time matches the signature in the value.

Returns:
 1. 200 OK and [CheckAuthResponse] when successful.
 2. 500 Internal Server Error and an error message when there is an unexpected error
*/
func (h *Handler) CheckAuth(c *gin.Context) {
	res := CheckAuthResponse{Authenticated: false}

	value, err := c.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		// Session cookie is missing
		c.JSON(http.StatusOK, res)
		return
	}

	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		// Session cookie value has incorrect format
		c.JSON(http.StatusOK, res)
		return
	}

	expiryStr := parts[0]
	expiryInt, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		// Expiry time is not an int
		c.JSON(http.StatusOK, res)
		return
	}

	if expiryInt <= time.Now().Unix() {
		// Cookie is expired
		c.JSON(http.StatusOK, res)
		return
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Println("Session secret could not be retrieved from SESSION_SECRET env var")
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	signature := parts[1]
	if !hmac.Equal([]byte(signature), []byte(generateHmacSha256Signature(expiryStr, sessionSecret))) {
		// Signature mismatch
		c.JSON(http.StatusOK, res)
		return
	}

	res.Authenticated = true
	c.JSON(http.StatusOK, res)
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
	expiryInt := time.Now().Add(24 * time.Hour).Unix()
	expiryStr := strconv.FormatInt(expiryInt, 10)
	signature := generateHmacSha256Signature(expiryStr, sessionSecret)
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

/*
Generates a HMAC-SHA256 signature of the given payload using the provided key as the hash key. Returns it as a hex
string.
*/
func generateHmacSha256Signature(payload string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload))
	digest := mac.Sum(nil)
	signature := hex.EncodeToString(digest)
	return signature
}
