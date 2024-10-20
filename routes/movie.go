package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/controllers"
	"github.com/raunak173/bms-go/middleware"
)

func MovieRoutes(c *gin.Engine) {
	Movie := c.Group("/movies")
	{
		Movie.GET("/", controllers.GetAllMovies)
		Movie.POST("/", middleware.RequireAuth, controllers.CreateMovie)
		Movie.GET("/:id", controllers.GetMovieByID)
	}
}
