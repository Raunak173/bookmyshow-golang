package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/helpers"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
	"gorm.io/gorm/clause"
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

func ReserveSeats(c *gin.Context) {
	var request struct {
		ShowID uint   `json:"show_id"`
		Seats  []uint `json:"seat_ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Start a GORM transaction
	tx := initializers.Db.Begin()

	// Reserve seats
	for _, seatID := range request.Seats {
		var seat models.Seat
		// Use FOR UPDATE to lock the seat row until transaction is complete
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND show_time_id = ?", seatID, request.ShowID).
			First(&seat).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
			return
		}
		// Check if the seat is available
		if !seat.IsAvailable || seat.IsReserved || seat.IsBooked {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Seat is already booked or reserved"})
			return
		}

		// Reserve the seat
		seat.IsReserved = true
		seat.IsAvailable = false
		if err := tx.Save(&seat).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reserve seat"})
			return
		}
	}

	// Commit the transaction
	tx.Commit()

	// Schedule a job to unreserve the seats after 5 minutes
	go helpers.UnReserveSeats(request.Seats, 5*time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": "Seats reserved successfully for 5 minutes"})
}

// BookSeats function to book reserved seats
func BookSeats(c *gin.Context) {
	var request struct {
		ShowID uint   `json:"show_id"`
		Seats  []uint `json:"seat_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	//We are checking we are authorized or not
	user, _ := c.Get("user")
	//We get userDetails, because we need to check that we are admin or not
	userDetails := user.(models.User)
	userId := userDetails.ID

	// Start a GORM transaction
	tx := initializers.Db.Begin()

	// Check if the seats are reserved and still valid
	var totalPrice float32
	var seats []models.Seat

	for _, seatID := range request.Seats {
		var seat models.Seat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND show_time_id = ?", seatID, request.ShowID).
			First(&seat).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
			return
		}

		// Check if the seat is reserved and not yet booked
		if !seat.IsReserved || seat.IsBooked {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Seat is not reserved or reservation expired"})
			return
		}

		// Add seat price to total
		totalPrice += seat.Price

		// Append seat to the list
		seats = append(seats, seat)
	}

	// Create an order
	order := models.Order{
		UserID:     userId,
		ShowTimeID: request.ShowID,
		TotalPrice: totalPrice,
		Seats:      seats,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create order"})
		return
	}

	// Mark the seats as booked and no longer reserved
	for i := range seats {
		seats[i].IsBooked = true
		seats[i].IsReserved = false
		if err := tx.Save(&seats[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to book seat"})
			return
		}
	}

	// Commit the transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Seats booked successfully",
		"order_id":    order.ID,
		"total_price": totalPrice,
	})
}
