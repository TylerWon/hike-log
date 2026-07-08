package router

import (
	"github.com/TylerWon/hike-log/backend/handler"
	"github.com/gin-gonic/gin"
)

// Creates a Gin router and registers API routes.
func New(handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			hikes := v1.Group("/hikes")
			hikes.GET("", handler.ListHike)
			hikes.GET("/:id", handler.RetrieveHike)
		}
	}

	router.GET("/health", handler.HealthCheck)

	return router
}
