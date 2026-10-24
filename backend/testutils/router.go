package testutils

import (
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/gin-gonic/gin"
)

// Creates a router to use during testing.
func NewRouter(handler *handler.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.New(handler)
}
