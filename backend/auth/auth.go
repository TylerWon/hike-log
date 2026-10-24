package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ONE_DAY_IN_SECONDS  = 24 * 60 * 60
	SESSION_COOKIE_NAME = "hike-log-session"
)

/*
Creates the value for the session cookie.

The value is "<cookie expiry time>.<HMAC-SHA256 hex signature of the expiry time>". The signature is created using the
given sessionSecret. It is included so that the API can verify the authenticity of the cookie when it is provided to
restricted routes since signatures will only match when they are created with the same secret.
*/
func CreateSessionCookieValue(sessionSecret string) string {
	expiryInt := time.Now().Add(24 * time.Hour).Unix()
	expiryStr := strconv.FormatInt(expiryInt, 10)
	signature := GenerateHmacSha256Signature(expiryStr, sessionSecret)
	return expiryStr + "." + signature
}

/*
Generates a HMAC-SHA256 signature of the given payload using the provided key as the hash key. Returns it as a hex
string.
*/
func GenerateHmacSha256Signature(payload string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload))
	digest := mac.Sum(nil)
	signature := hex.EncodeToString(digest)
	return signature
}

/*
Checks that a request is by an authenticated user.

A request is authenticated if it has a session cookie and the cookie's value satisfies the following criteria:
1. The expiry time in the value is valid, and
2. A HMAC-SHA256 signature of the expiry time matches the signature in the value.
*/
func IsAuthenticated(c *gin.Context) (bool, error) {
	value, err := c.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return false, AuthError{reason: "Session cookie is missing from request."}
	}

	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return false, AuthError{reason: "Session cookie value has invalid format."}
	}

	expiryStr := parts[0]
	expiryInt, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return false, AuthError{reason: "Expiry time in cookie value is not an integer."}
	}

	if expiryInt <= time.Now().Unix() {
		return false, AuthError{reason: "Cookie is expired."}
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		return false, ConfigurationError{reason: "Session secret could not be retrieved from SESSION_SECRET env var"}
	}

	signature := parts[1]
	if !hmac.Equal([]byte(signature), []byte(GenerateHmacSha256Signature(expiryStr, sessionSecret))) {
		return false, AuthError{reason: "Signature in cookie value is invalid."}
	}

	return true, nil
}

/*
Sets a session cookie on the response.

The session cookie is used to authenticate requests to restricted API routes and is only valid for 24hrs after being
issued. The cookie is named [SESSION_COOKIE_NAME].
*/
func SetSessionCookie(c *gin.Context) error {
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		return ConfigurationError{reason: "Session secret could not be retrieved from SESSION_SECRET env var"}
	}

	cookieValue := CreateSessionCookieValue(sessionSecret)

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

	return nil
}
