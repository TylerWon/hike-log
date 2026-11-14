package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TylerWon/hike-log/backend/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	OneDayInSeconds   = 24 * 60 * 60
	SessionCookieName = "hike-log-session"
)

// Auth handles authentication for the app.
type Auth interface {
	CreateSessionCookieValue() string
	GenerateHmacSha256Signature(payload string) string
	IsAuthenticated(c *gin.Context) (bool, error)
	SetSessionCookie(c *gin.Context)
	VerifyPassword(password string) bool
}

// authImpl is an implementation of the Auth interface.
type authImpl struct {
	authConfig config.AuthConfig
}

// Creates a new Auth.
func New(authConfig config.AuthConfig) (Auth, error) {
	if authConfig.AdminPasswordHash == "" {
		return nil, errors.New("adminPasswordHash is empty")
	} else if authConfig.SessionSecret == "" {
		return nil, errors.New("sessionSecret is empty")
	}

	return &authImpl{authConfig: authConfig}, nil
}

/*
Creates the value for the session cookie.

The value is "<cookie expiry time>.<HMAC-SHA256 hex signature of the expiry time>". The signature is created using the
sessionSecret. It is included so that the API can verify the authenticity of the cookie when it is provided to
restricted routes since signatures will only match when they are created with the same secret.
*/
func (a *authImpl) CreateSessionCookieValue() string {
	expiryInt := time.Now().Add(24 * time.Hour).Unix()
	expiryStr := strconv.FormatInt(expiryInt, 10)
	signature := a.GenerateHmacSha256Signature(expiryStr)
	return expiryStr + "." + signature
}

/*
Generates a HMAC-SHA256 signature of the given payload using the sessionSecret as the hash key. Returns it as a hex
string.
*/
func (a *authImpl) GenerateHmacSha256Signature(payload string) string {
	mac := hmac.New(sha256.New, []byte(a.authConfig.SessionSecret))
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
func (a *authImpl) IsAuthenticated(c *gin.Context) (bool, error) {
	value, err := c.Cookie(SessionCookieName)
	if err != nil {
		return false, errors.New("Session cookie is missing from request.")
	}

	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return false, errors.New("Session cookie value has invalid format.")
	}

	expiryStr := parts[0]
	expiryInt, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return false, errors.New("Expiry time in cookie value is not an integer.")
	}

	if expiryInt <= time.Now().Unix() {
		return false, errors.New("Cookie is expired.")
	}

	signature := parts[1]
	if !hmac.Equal([]byte(signature), []byte(a.GenerateHmacSha256Signature(expiryStr))) {
		return false, errors.New("Signature in cookie value is invalid.")
	}

	return true, nil
}

/*
Sets a session cookie on the response.

The session cookie is used to authenticate requests to restricted API routes and is only valid for 24hrs after being
issued. The cookie is named [SessionCookieName].
*/
func (a *authImpl) SetSessionCookie(c *gin.Context) {
	cookieValue := a.CreateSessionCookieValue()

	// Only allow cookie to be sent between same sites (same site = scheme + host are the same)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,
		cookieValue,
		OneDayInSeconds,
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
Checks if the given password is correct.
*/
func (a *authImpl) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.authConfig.AdminPasswordHash), []byte(password))
	return err == nil
}
