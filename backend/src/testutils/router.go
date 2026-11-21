package testutils

import (
	"testing"

	"github.com/TylerWon/hike-log/backend/src/handler"
	"github.com/TylerWon/hike-log/backend/src/middleware"
	"github.com/TylerWon/hike-log/backend/src/router"
	"github.com/gin-gonic/gin"
)

// Creates a test router with only auth routes.
func NewAuthRouter(t *testing.T, h *handler.AuthHandler, m *middleware.Middleware) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterAuthRoutes(r.Group(""), h, m)
	return r
}

// Creates a test router with only hike routes.
func NewHikesRouter(t *testing.T, h *handler.HikesHandler, m *middleware.Middleware) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterHikesRoutes(r.Group(""), h, m)
	return r
}

// Creates a test router with only photo routes.
func NewPhotosRouter(t *testing.T, h *handler.PhotosHandler, m *middleware.Middleware) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterPhotosRoutes(r.Group(""), h, m)
	return r
}

// Creates a test router with health and ready routes.
func NewHealthRouter(t *testing.T, h *handler.HealthHandler) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterHealthRoutes(r, h)
	return r
}
