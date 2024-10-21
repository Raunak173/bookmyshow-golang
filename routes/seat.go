package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/controllers"
)

func SeatRoutes(c *gin.Engine) {
	Seat := c.Group("/seats")
	{
		Seat.GET("/showtime/:id", controllers.GetSeatLayout)
	}
}
