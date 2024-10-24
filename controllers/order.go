package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
)

type OrderResponse struct {
	ID         uint     `json:"id"`
	TotalPrice float32  `json:"total_price"`
	Seats      []string `json:"seats"`
}

func GetOrders(c *gin.Context) {
	user, _ := c.Get("user")
	userDetails := user.(models.User)
	userId := userDetails.ID

	var orders []models.Order
	if err := initializers.Db.Where("user_id= ?", userId).Preload("Seats").Find(&orders).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No order found for the given id"})
		return
	}

	var orderResponses []OrderResponse
	for _, order := range orders {
		var seatNumbers []string
		for _, seat := range order.Seats {
			seatNumbers = append(seatNumbers, seat.SeatNumber) // Append seat number directly
		}
		// Append the transformed order to the response slice
		orderResponses = append(orderResponses, OrderResponse{
			ID:         order.ID,
			TotalPrice: order.TotalPrice,
			Seats:      seatNumbers,
		})
	}

	c.JSON(200, gin.H{
		"orders": orderResponses,
	})
}
