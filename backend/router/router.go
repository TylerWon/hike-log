package router

import (
	"os"

	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Creates a Gin router and registers middleware and API routes.
func New(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	// Middleware (must come before routes)
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
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
				hikes.POST("", handler.CreateHike)

				// /api/v1/hikes/:hikeId
				hike := hikes.Group("/:hikeId")
				{
					// /api/v1/hikes/:hikeId/photos
					photos := hike.Group("/photos")
					{
						photos.POST("/upload-url", handler.CreatePhotoUploadURL)
					}
				}
			}
		}
	}

	router.GET("/health", handler.HealthCheck)

	return router
}
