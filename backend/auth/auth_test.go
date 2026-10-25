package auth_test

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

var sessionSecret = os.Getenv("SESSION_SECRET")

type authTestSuite struct {
	suite.Suite
}

func (suite *authTestSuite) TestCreateSessionCookieValue_ReturnsValue() {
	value := auth.CreateSessionCookieValue(sessionSecret)

	parts := strings.Split(value, ".")
	suite.Len(parts, 2)

	expiryStr := parts[0]
	expiryInt, err := strconv.ParseInt(expiryStr, 10, 64)
	suite.NoError(err)
	suite.InDelta(time.Now().Add(24*time.Hour).Unix(), expiryInt, 5)

	signature := parts[1]
	suite.Equal(auth.GenerateHmacSha256Signature(expiryStr, sessionSecret), signature)
}

func (suite *authTestSuite) TestCreateSessionCookieValue_ReturnsDifferentSignaturesForDifferentSecrets() {
	// Use synctest.Test to freeze time so two CreateSessionCookieValue() calls use the same expiry time to create the signature
	synctest.Test(suite.T(), func(t *testing.T) {
		value1 := auth.CreateSessionCookieValue("a")
		value2 := auth.CreateSessionCookieValue("b")

		value1Signature := strings.Split(value1, ".")[1]
		value2Signature := strings.Split(value2, ".")[1]

		suite.NotEqual(value1Signature, value2Signature)
	})
}

func (suite *authTestSuite) TestGenerateHmacSha256Signature_ReturnsSignature() {
	payload := "abc"
	key := "a"
	signature := auth.GenerateHmacSha256Signature(payload, key)
	suite.Equal("b4ffe67d267d9eb19f81ccd6ffcde1443294d00892d5b648f03e590b56250ced", signature)
}

func (suite *authTestSuite) TestGenerateHmacSha256Signature_ReturnsDifferentSignaturesForDifferentPayloads() {
	payload1 := "abc"
	payload2 := "cba"
	key := "a"
	signature1 := auth.GenerateHmacSha256Signature(payload1, key)
	signature2 := auth.GenerateHmacSha256Signature(payload2, key)
	suite.NotEqual(signature1, signature2)
}

func (suite *authTestSuite) TestGenerateHmacSha256Signature_ReturnsDifferentSignaturesForDifferentKeys() {
	payload := "abc"
	key1 := "a"
	key2 := "b"
	signature1 := auth.GenerateHmacSha256Signature(payload, key1)
	signature2 := auth.GenerateHmacSha256Signature(payload, key2)
	suite.NotEqual(signature1, signature2)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenCookieIsMissing() {
	c := testutils.NewGinContext(nil)

	res, err := auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
	suite.ErrorAs(err, &auth.AuthError{})
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenCookieValueIsInvalid() {
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: "invalid value",
	})

	res, err := auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
	suite.ErrorAs(err, &auth.AuthError{})
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenExpiryTimeIsNotAnInteger() {
	expiry := "abc"
	signature := auth.GenerateHmacSha256Signature(expiry, sessionSecret)
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	res, err := auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
	suite.ErrorAs(err, &auth.AuthError{})
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenCookieIsExpired() {
	expiry := "1" // Expiry time is not within 24hrs
	signature := auth.GenerateHmacSha256Signature(expiry, sessionSecret)
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	res, err := auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
	suite.ErrorAs(err, &auth.AuthError{})
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenSessionSecretIsMissing() {
	c := testutils.NewGinContext(testutils.CreateValidSessionCookie())

	suite.T().Setenv("SESSION_SECRET", "") // Set secret value to empty string

	res, err := auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
	suite.ErrorAs(err, &auth.ConfigurationError{})
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenSignatureIsInvalid() {
	expiry := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	signature := auth.GenerateHmacSha256Signature(expiry, "abc") // Do not use session secret to generate signature
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	res, err := auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
	suite.ErrorAs(err, &auth.AuthError{})
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsTrue() {
	c := testutils.NewGinContext(testutils.CreateValidSessionCookie())

	res, err := auth.IsAuthenticated(c)

	suite.True(res)
	suite.NoError(err)
}

func (suite *authTestSuite) TestSetSessionCookie_ReturnsErrorWhenSessionSecretIsMissing() {
	c := testutils.NewGinContext(nil)
	suite.T().Setenv("SESSION_SECRET", "") // Set secret value to empty string
	err := auth.SetSessionCookie(c)
	suite.Error(err)
}

func (suite *authTestSuite) TestSetSessionCookie_SetsCookie() {
	c := testutils.NewGinContext(nil)

	err := auth.SetSessionCookie(c)

	suite.NoError(err)

	raw := c.Writer.Header().Get("Set-Cookie")
	sessionCookie, err := http.ParseSetCookie(raw)
	suite.NoError(err)

	suite.Equal(http.SameSiteLaxMode, sessionCookie.SameSite)
	suite.Equal(auth.SESSION_COOKIE_NAME, sessionCookie.Name)
	suite.NotEmpty(sessionCookie.Value)
	suite.Equal(auth.ONE_DAY_IN_SECONDS, sessionCookie.MaxAge)
	suite.Equal("/", sessionCookie.Path)
	suite.Empty(sessionCookie.Domain)
	suite.True(sessionCookie.Secure)
	suite.True(sessionCookie.HttpOnly)
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}
