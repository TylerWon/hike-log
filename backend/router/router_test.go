package router_test

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type routerTestSuite struct {
	suite.Suite
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(routerTestSuite))
}

func (suite *routerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *routerTestSuite) TestNew_ReturnsErrorWhenFrontendURLIsEmpty() {
	h := &handler.Handler{}
	m := &middleware.Middleware{}
	r, err := router.New(h, m, "")
	suite.Require().Error(err)
	suite.Require().Nil(r)
}

func (suite *routerTestSuite) TestNew_CreatesRouter() {
	h := &handler.Handler{}
	m := &middleware.Middleware{}
	r, err := router.New(h, m, "http://localhost:5173")
	suite.Require().NoError(err)
	suite.Require().NotNil(r)

	var paths []string
	for _, route := range r.Routes() {
		paths = append(paths, route.Method+" "+route.Path)
	}

	// auth
	suite.Require().Contains(paths, "POST /api/v1/auth")
	suite.Require().Contains(paths, "GET /api/v1/auth")

	// hikes
	suite.Require().Contains(paths, "GET /api/v1/hikes")
	suite.Require().Contains(paths, "POST /api/v1/hikes")
	suite.Require().Contains(paths, "DELETE /api/v1/hikes/:hikeId")
	suite.Require().Contains(paths, "PUT /api/v1/hikes/:hikeId")
	suite.Require().Contains(paths, "POST /api/v1/hikes/:hikeId/photos")
	suite.Require().Contains(paths, "POST /api/v1/hikes/:hikeId/photos/presigned-urls")

	// photos
	suite.Require().Contains(paths, "DELETE /api/v1/photos/:photoId")
	suite.Require().Contains(paths, "PUT /api/v1/photos/:photoId")

	// health
	suite.Require().Contains(paths, "GET /health")
	suite.Require().Contains(paths, "GET /ready")
}
