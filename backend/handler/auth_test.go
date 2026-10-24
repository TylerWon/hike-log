package handler_test

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/testutils"
)

var sessionSecret = os.Getenv("SESSION_SECRET")

func (suite *handlerTestSuite) TestCheckAuth_ReturnsUnauthenticated() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil) // Session cookie missing

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *handlerTestSuite) TestCheckAuth_ReturnsAuthenticated() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/api/v1/auth", nil, &http.Cookie{
		Name:  auth.SESSION_COOKIE_NAME,
		Value: auth.CreateSessionCookieValue(sessionSecret),
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
	suite.Equal(auth.SESSION_COOKIE_NAME, cookies[0].Name)
}
