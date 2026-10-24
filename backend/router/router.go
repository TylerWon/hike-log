package router

import (
	"os"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/TylerWon/hike-log/backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Creates a Gin router and registers middleware and API routes.
func New(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	// CORS middleware (must come before routes)
	corsConfig := cors.DefaultConfig()
	// Allow frontend to access API
	corsConfig.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	// Allow frontend to pass session cookie for auth
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Routes
	// /api
	api := router.Group("/api")
	{
		// /api/v1
		v1 := api.Group("/v1")
		{
			// /api/v1/hikes
			hikes := v1.Group("/hikes")
			{
				hikes.GET("", handler.ListHikes)
				hikes.POST("", middleware.RequireAuth, handler.CreateHike)

				// /api/v1/hikes/:hikeId
				hike := hikes.Group("/:hikeId")
				{
					hike.DELETE("", middleware.RequireAuth, handler.DeleteHike)
					hike.PUT("", middleware.RequireAuth, handler.UpdateHike)

					// /api/v1/hikes/:hikeId/photos
					photos := hike.Group("/photos")
					{
						photos.POST("", middleware.RequireAuth, handler.CreatePhotos)
						photos.POST("/presigned-urls", middleware.RequireAuth, handler.CreatePresignedURLs)
					}
				}
			}

			// /api/v1/photos
			photos := v1.Group("/photos")
			{
				// /api/v1/photos/:photoId
				photo := photos.Group("/:photoId")
				{
					photo.DELETE("", middleware.RequireAuth, handler.DeletePhoto)
					photo.PUT("", middleware.RequireAuth, handler.UpdatePhoto)
				}
			}

			// /api/v1/auth
			auth := v1.Group("/auth")
			{
				auth.POST("", handler.Login)
				auth.GET("", handler.CheckAuth)
			}
		}
	}

	// /health
	router.GET("/health", handler.HealthCheck)

	return router
}
