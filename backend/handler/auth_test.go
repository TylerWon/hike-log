package handler_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/testutils"
)

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

func (suite *handlerTestSuite) TestLogin_Returns500ErrorWhenConfigurationIsMissing() {
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
