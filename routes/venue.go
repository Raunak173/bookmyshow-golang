package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/controllers"
	"github.com/raunak173/bms-go/middleware"
)

func VenueRoutes(c *gin.Engine) {
	Venue := c.Group("/venues")
	{
		Venue.GET("/", controllers.GetAllVenues)
		Venue.POST("/", middleware.RequireAuth, controllers.CreateVenue)
		Venue.POST("/:id/movies/add", middleware.RequireAuth, controllers.AddMoviesInVenue)
		Venue.GET("/:id", controllers.GetVenueByID)
	}
}
