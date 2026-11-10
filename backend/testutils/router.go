package testutils

import (
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/TylerWon/hike-log/backend/router"
	"github.com/gin-gonic/gin"
)

// Creates a test router with only auth routes.
func NewAuthRouter(h *handler.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterAuthRoutes(r.Group(""), h)
	return r
}

// Creates a test router with only hike routes.
func NewHikesRouter(h *handler.HikesHandler, m *middleware.Middleware) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterHikesRoutes(r.Group(""), h, m)
	return r
}

// Creates a test router with only photo routes.
func NewPhotosRouter(h *handler.PhotosHandler, m *middleware.Middleware) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterPhotosRoutes(r.Group(""), h, m)
	return r
}
