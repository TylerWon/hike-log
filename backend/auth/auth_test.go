package auth_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
	auth auth.Auth
}

func (suite *authTestSuite) SetupSuite() {
	suite.auth = testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, testutils.SESSION_SECRET)
}

func (suite *authTestSuite) TestCreateSessionCookieValue_ReturnsValue() {
	value := suite.auth.CreateSessionCookieValue()

	parts := strings.Split(value, ".")
	suite.Len(parts, 2)

	expiryStr := parts[0]
	expiryInt, err := strconv.ParseInt(expiryStr, 10, 64)
	suite.NoError(err)
	suite.InDelta(time.Now().Add(24*time.Hour).Unix(), expiryInt, 5)

	signature := parts[1]
	suite.Equal(suite.auth.GenerateHmacSha256Signature(expiryStr), signature)
}

func (suite *authTestSuite) TestGenerateHmacSha256Signature_ReturnsSignature() {
	signature := suite.auth.GenerateHmacSha256Signature("abc")
	suite.Equal("545003ead8123a7a230d5e79249e903768f4a407ebebbd6b88db6c519b5fb222", signature)
}

func (suite *authTestSuite) TestGenerateHmacSha256Signature_ReturnsDifferentSignaturesForDifferentPayloads() {
	signature1 := suite.auth.GenerateHmacSha256Signature("abc")
	signature2 := suite.auth.GenerateHmacSha256Signature("cba")
	suite.NotEqual(signature1, signature2)
}

func (suite *authTestSuite) TestGenerateHmacSha256Signature_ReturnsDifferentSignaturesForDifferentSecrets() {
	auth1 := testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, "a")
	auth2 := testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, "b")

	payload := "abc"
	signature1 := auth1.GenerateHmacSha256Signature(payload)
	signature2 := auth2.GenerateHmacSha256Signature(payload)

	suite.NotEqual(signature1, signature2)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenCookieIsMissing() {
	c := testutils.NewGinContext(nil)

	res, err := suite.auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenCookieValueIsInvalid() {
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: "invalid value",
	})

	res, err := suite.auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenExpiryTimeIsNotAnInteger() {
	expiry := "abc"
	signature := suite.auth.GenerateHmacSha256Signature(expiry)
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	res, err := suite.auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenCookieIsExpired() {
	expiry := "1" // Expiry time is not within 24hrs
	signature := suite.auth.GenerateHmacSha256Signature(expiry)
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	res, err := suite.auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsFalseWhenSignatureIsInvalid() {
	expiry := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	signature := suite.auth.GenerateHmacSha256Signature(expiry)
	c := testutils.NewGinContext(&http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature + "123", // Add extra characters to signature to invalidate it
	})

	res, err := suite.auth.IsAuthenticated(c)

	suite.False(res)
	suite.Error(err)
}

func (suite *authTestSuite) TestIsAuthenticated_ReturnsTrue() {
	c := testutils.NewGinContext(testutils.CreateSessionCookie(suite.auth))

	res, err := suite.auth.IsAuthenticated(c)

	suite.True(res)
	suite.NoError(err)
}

func (suite *authTestSuite) TestSetSessionCookie_SetsCookie() {
	c := testutils.NewGinContext(nil)

	suite.auth.SetSessionCookie(c)

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
