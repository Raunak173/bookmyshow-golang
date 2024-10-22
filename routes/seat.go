package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/controllers"
	"github.com/raunak173/bms-go/middleware"
)

func SeatRoutes(c *gin.Engine) {
	Seat := c.Group("/seats")
	{
		Seat.GET("/showtime/:id", controllers.GetSeatLayout)
		Seat.POST("/showtime/reserve", middleware.RequireAuth, controllers.ReserveSeats)
		Seat.POST("/showtime/book", middleware.RequireAuth, controllers.BookSeats)
	}
}
