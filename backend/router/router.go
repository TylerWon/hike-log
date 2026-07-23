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
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			hikes := v1.Group("/hikes")
			hikes.GET("", handler.ListHike)
		}
	}

	router.GET("/health", handler.HealthCheck)

	return router
}
