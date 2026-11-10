package router

import (
	"errors"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const API_BASE_PATH = "/api/v1"

// Creates a Gin router and registers middleware and API routes.
func New(h *handler.Handler, m *middleware.Middleware, frontendURL string) (*gin.Engine, error) {
	if frontendURL == "" {
		return nil, errors.New("frontendURL is empty")
	}

	router := gin.Default()

	configureCORSMiddleware(router, frontendURL)

	apiRouteGroup := router.Group(API_BASE_PATH)
	RegisterAuthRoutes(apiRouteGroup, h.Auth)
	RegisterHikesRoutes(apiRouteGroup, h.Hikes, m)
	RegisterPhotosRoutes(apiRouteGroup, h.Photos, m)

	router.GET("/health", handler.HealthCheck)

	return router, nil
}

// Configures CORS middleware on the router.
func configureCORSMiddleware(router *gin.Engine, frontendURL string) {
	cfg := cors.DefaultConfig()

	// Allow frontend to access API
	cfg.AllowOrigins = []string{frontendURL}
	// Allow frontend to pass session cookie for auth
	cfg.AllowCredentials = true

	router.Use(cors.New(cfg))
}

// Registers authentication API routes under the given apiRouteGroup.
func RegisterAuthRoutes(apiRouteGroup *gin.RouterGroup, h *handler.AuthHandler) {
	auth := apiRouteGroup.Group("/auth")
	{
		auth.POST("", h.Login)
		auth.GET("", h.CheckAuth)
	}
}

// Registers hike API routes under the given apiRouteGroup.
func RegisterHikesRoutes(apiRouteGroup *gin.RouterGroup, h *handler.HikesHandler, m *middleware.Middleware) {
	hikes := apiRouteGroup.Group("/hikes")
	{
		hikes.GET("", h.ListHikes)
		hikes.POST("", m.RequireAuth, h.CreateHike)

		hike := hikes.Group("/:hikeId")
		{
			hike.DELETE("", m.RequireAuth, h.DeleteHike)
			hike.PUT("", m.RequireAuth, h.UpdateHike)

			photos := hike.Group("/photos")
			{
				photos.POST("", m.RequireAuth, h.CreatePhotos)
				photos.POST("/presigned-urls", m.RequireAuth, h.CreatePresignedURLs)
			}
		}
	}
}

// Registers photo API routes under the given apiRouteGroup.
func RegisterPhotosRoutes(apiRouteGroup *gin.RouterGroup, h *handler.PhotosHandler, m *middleware.Middleware) {
	photos := apiRouteGroup.Group("/photos")
	{
		photo := photos.Group("/:photoId")
		{
			photo.DELETE("", m.RequireAuth, h.DeletePhoto)
			photo.PUT("", m.RequireAuth, h.UpdatePhoto)
		}
	}
}
