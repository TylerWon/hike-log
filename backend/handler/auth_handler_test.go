package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/TylerWon/hike-log/backend/auth"
	"github.com/TylerWon/hike-log/backend/config"
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type authHandlerTestSuite struct {
	suite.Suite
	auth   auth.Auth
	router *gin.Engine
}

func (suite *authHandlerTestSuite) SetupSuite() {
	suite.auth = testutils.NewAuth(suite.T(), config.AuthConfig{})
	suite.router = testutils.NewAuthRouter(suite.T(), testutils.NewAuthHandler(suite.T(), suite.auth))
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(authHandlerTestSuite))
}

func (suite *authHandlerTestSuite) TestNewAuthHandler_CreatesAuthHandler() {
	h := handler.NewAuthHandler(suite.auth)
	suite.Require().NotNil(h)
}

func (suite *authHandlerTestSuite) TestCheckAuth_ReturnsUnauthenticated() {
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodGet, "/auth", nil) // Session cookie missing

	suite.Require().Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().False(response.Authenticated)
}

func (suite *authHandlerTestSuite) TestCheckAuth_ReturnsAuthenticated() {
	res := testutils.SendRequest(suite.T(),
		suite.router,
		http.MethodGet,
		"/auth",
		nil,
		testutils.CreateSessionCookie(suite.T(), suite.auth),
	)

	suite.Require().Equal(http.StatusOK, res.Code)

	var response handler.CheckAuthResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	suite.Require().NoError(err)

	suite.Require().True(response.Authenticated)
}

func (suite *authHandlerTestSuite) TestLogin_Returns400ErrorWhenRequestBodyHasInvalidField() {
	body := map[string]any{
		"password": 123, // number instead of string
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodPost, "/auth", reqBody)

	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *authHandlerTestSuite) TestLogin_Returns400ErrorWhenRequestBodyIsMissingField() {
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodPost, "/auth", nil)
	suite.Require().Equal(http.StatusBadRequest, res.Code)
}

func (suite *authHandlerTestSuite) TestLogin_Returns401ErrorWhenPasswordIsIncorrect() {
	body := map[string]any{
		"password": "incorrect password",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodPost, "/auth", reqBody)

	suite.Require().Equal(http.StatusUnauthorized, res.Code)
}

func (suite *authHandlerTestSuite) TestLogin_SetsSessionCookie() {
	body := map[string]any{
		"password": "admin",
	}
	reqBody := testutils.SerializeJSONRequestBody(suite.T(), body)
	res := testutils.SendRequest(suite.T(), suite.router, http.MethodPost, "/auth", reqBody)

	suite.Require().Equal(http.StatusOK, res.Code)

	cookies := res.Result().Cookies()
	suite.Require().Len(cookies, 1)
	suite.Require().Equal(auth.SessionCookieName, cookies[0].Name)
}
