package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type authTestSuite struct {
	suite.Suite
	auth   auth.Auth
	router *gin.Engine
}

func (suite *authTestSuite) SetupSuite() {
	suite.auth = testutils.NewAuth(suite.T(), testutils.ADMIN_PASSWORD_HASH, testutils.SESSION_SECRET)
	suite.router = testutils.NewAuthRouter(handler.NewAuthHandler(suite.auth))
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(authTestSuite))
}

func (suite *authTestSuite) TestCheckAuth_ReturnsUnauthenticated() {
	res := testutils.SendRequest(suite.router, http.MethodGet, "/auth", nil) // Session cookie missing

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.False(response.Authenticated)
}

func (suite *authTestSuite) TestCheckAuth_ReturnsAuthenticated() {
	res := testutils.SendRequest(
		suite.router,
		http.MethodGet,
		"/auth",
		nil,
		testutils.CreateSessionCookie(suite.auth),
	)

	suite.Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.NoError(err)

	suite.True(response.Authenticated)
}

func (suite *authTestSuite) TestLogin_Returns400ErrorWhenRequestBodyHasInvalidField() {
	body := map[string]any{
		"password": 123, // number instead of string
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/auth", reqBody)

	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *authTestSuite) TestLogin_Returns400ErrorWhenRequestBodyIsMissingField() {
	res := testutils.SendRequest(suite.router, http.MethodPost, "/auth", nil)
	suite.Equal(http.StatusBadRequest, res.Code)
}

func (suite *authTestSuite) TestLogin_Returns401ErrorWhenPasswordIsIncorrect() {
	body := map[string]any{
		"password": "incorrect password",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/auth", reqBody)

	suite.Equal(http.StatusUnauthorized, res.Code)
}

func (suite *authTestSuite) TestLogin_SetsSessionCookie() {
	body := map[string]any{
		"password": testutils.ADMIN_PASSWORD,
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.router, http.MethodPost, "/auth", reqBody)

	suite.Equal(http.StatusOK, res.Code)

	cookies := res.Result().Cookies()
	suite.Len(cookies, 1)
	suite.Equal(auth.SESSION_COOKIE_NAME, cookies[0].Name)
}
