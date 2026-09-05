package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/gin-gonic/gin"
)

// Creates a router to use during testing.
func NewRouter(t *testing.T, handler *handler.Handler) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	return router.New(handler)
}
