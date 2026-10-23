package handler_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/testutils"
)

var sessionSecret = os.Getenv("SESSION_SECRET")

func (suite *handlerTestSuite) TestCheckAuth_ReturnsUnauthenticatedWhenCookieIsMissing() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil)

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *handlerTestSuite) TestCheckAuth_ReturnsUnauthenticatedWhenCookieValueIsInvalid() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  handler.SESSION_COOKIE_NAME,
		Value: "invalid value",
	})

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *handlerTestSuite) TestCheckAuth_ReturnsUnauthenticatedWhenExpiryTimeIsNotAnInteger() {
	expiry := "abc"
	signature := handler.GenerateHmacSha256Signature(expiry, sessionSecret)
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  handler.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *handlerTestSuite) TestCheckAuth_ReturnsUnauthenticatedWhenCookieIsExpired() {
	expiry := strconv.FormatInt(time.Now().Unix(), 10) // Set expiry to now so it is invalid by request time
	signature := handler.GenerateHmacSha256Signature(expiry, sessionSecret)
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  handler.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *handlerTestSuite) TestCheckAuth_ReturnsUnauthenticatedWhenSignatureDoesNotMatch() {
	expiry := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	signature := handler.GenerateHmacSha256Signature(expiry, "abc") // Do not use session secret to generate signature
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  handler.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *handlerTestSuite) TestCheckAuth_Returns500ErrorWhenSecretIsMissing() {
	expiry := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	signature := handler.GenerateHmacSha256Signature(expiry, sessionSecret)

	suite.T().Setenv("SESSION_SECRET", "") // Mock missing secret value

	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  handler.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestCheckAuth_ReturnsAuthenticated() {
	expiry := strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10)
	signature := handler.GenerateHmacSha256Signature(expiry, sessionSecret)
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  handler.SESSION_COOKIE_NAME,
		Value: expiry + "." + signature,
	})

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.True(response.Authenticated)
}

func (suite *handlerTestSuite) TestLogin_Returns400ErrorWhenRequestBodyHasInvalidField() {
	body := map[string]any{
		"password": 123, // number instead of string
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/auth", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestLogin_Returns400ErrorWhenRequestBodyIsMissingField() {
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/auth", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *handlerTestSuite) TestLogin_Returns401ErrorWhenPasswordIsIncorrect() {
	body := map[string]any{
		"password": "incorrect_password",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/auth", reqBody)

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *handlerTestSuite) TestLogin_Returns500ErrorWhenSecretIsMissing() {
	suite.T().Setenv("ADMIN_PASSWORD_HASH", "")

	body := map[string]any{
		"password": "admin",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/auth", reqBody)

	suite.Equal(http.StatusInternalServerError, res.Code)
}

func (suite *handlerTestSuite) TestLogin_SetsSessionCookie() {
	body := map[string]any{
		"password": "admin",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/api/v1/auth", reqBody)

	suite.Equal(http.StatusOK, res.Code)

	cookies := res.Result().Cookies()
	suite.Len(cookies, 1)

	sessionCookie := cookies[0]

	suite.Equal(http.SameSiteLaxMode, sessionCookie.SameSite)
	suite.Equal(handler.SESSION_COOKIE_NAME, sessionCookie.Name)
	suite.Equal(handler.ONE_DAY_IN_SECONDS, sessionCookie.MaxAge)
	suite.Equal("/", sessionCookie.Path)
	suite.Empty(sessionCookie.Domain)
	suite.True(sessionCookie.Secure)
	suite.True(sessionCookie.HttpOnly)

	parts := strings.Split(sessionCookie.Value, ".")
	suite.Len(parts, 2)

	expiry, err := strconv.ParseInt(parts[0], 10, 64)
	suite.NoError(err)
	suite.InDelta(time.Now().Add(24*time.Hour).Unix(), expiry, 5)

	mac := hmac.New(sha256.New, []byte(os.Getenv("SESSION_SECRET")))
	mac.Write([]byte(parts[0]))
	suite.Equal(hex.EncodeToString(mac.Sum(nil)), parts[1])
}
