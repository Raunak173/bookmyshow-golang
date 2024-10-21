package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/helpers"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
)

func GetSeatLayout(c *gin.Context) {
	showtimeID := c.Param("id")
	// Fetch the showtime to ensure it exists
	var showTime models.ShowTime
	if err := initializers.Db.Preload("Seats").First(&showTime, showtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShowTime not found"})
		return
	}
	// If seats are fetched successfully, convert to a matrix format based on row and seat number
	seatMatrix := helpers.CreateSeatMatrix(showTime.Seats)

	c.JSON(http.StatusOK, gin.H{
		"showtime": showTime.Timing,
		"venue":    showTime.VenueID,
		"seats":    seatMatrix,
	})
}
