package helpers

import (
	"fmt"
	"time"

	"github.com/raunak173/bms-go/models"
)

func FormatShowTime(t time.Time) string {
	return t.Format("15:04") // 24-hour format (16:38)
	// For 12-hour format with AM/PM: return t.Format("03:04 PM")
}

// GenerateSeatsForShowTime generates the default seat layout for a showtime
func GenerateSeatsForShowTime(showtimeID uint) []models.Seat {
	// Define the default seat layout (e.g., 5 rows with 10 seats each)
	rows := []string{"A", "B", "C", "D", "E"} // Example rows
	seatsPerRow := 10                         // Example: 10 seats per row

	var seats []models.Seat
	for _, row := range rows {
		for seatNumber := 1; seatNumber <= seatsPerRow; seatNumber++ {
			seat := models.Seat{
				SeatNumber:  fmt.Sprintf("%s%d", row, seatNumber), // e.g., A1, A2, B1, etc.
				IsAvailable: true,                                 // All seats are available initially
				IsReserved:  false,
				IsBooked:    false,
				Price:       250,
				ShowTimeID:  showtimeID,
			}
			seats = append(seats, seat)
		}
	}
	return seats
}

// [
//   { SeatNumber: "A1", IsAvailable: true, IsReserved: false, IsBooked: false, Price: 250, ShowTimeID: 1 },
//   { SeatNumber: "A2", IsAvailable: true, IsReserved: false, IsBooked: false, Price: 250, ShowTimeID: 1 },
//   ...
//   { SeatNumber: "E10", IsAvailable: true, IsReserved: false, IsBooked: false, Price: 250, ShowTimeID: 1 },
// ]

func CreateSeatMatrix(seats []models.Seat) map[string][]map[string]interface{} {
	seatMatrix := make(map[string][]map[string]interface{})

	for _, seat := range seats {
		row := string(seat.SeatNumber[0]) // Assume the first character represents the row
		seatData := map[string]interface{}{
			"seat_number":  seat.SeatNumber,
			"is_reserved":  seat.IsReserved,
			"is_booked":    seat.IsBooked,
			"is_available": seat.IsAvailable,
			"price":        seat.Price,
		}
		// Append seat data to the appropriate row
		seatMatrix[row] = append(seatMatrix[row], seatData)
	}

	return seatMatrix
}

// {
// 	"A": [
// 	  {"seat_number": "A1", "is_reserved": false, "is_booked": false, "is_available": true, "price": 250},
// 	  {"seat_number": "A2", "is_reserved": false, "is_booked": false, "is_available": true, "price": 250},
// 	],
// 	"B": [
// 	  {"seat_number": "B1", "is_reserved": false, "is_booked": false, "is_available": true, "price": 250},
// 	]
//   }
