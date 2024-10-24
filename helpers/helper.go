package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/raunak173/bms-go/initializers"
	"github.com/raunak173/bms-go/models"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
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

// UnreserveSeats unreserves the seats after the specified duration
func UnReserveSeats(seatIDs []uint, duration time.Duration) {
	time.Sleep(duration)

	tx := initializers.Db.Begin()

	for _, seatID := range seatIDs {
		var seat models.Seat
		if err := tx.First(&seat, seatID).Error; err != nil {
			tx.Rollback()
			return
		}

		// Unreserve if the seat is still reserved but not booked
		if seat.IsReserved && !seat.IsBooked {
			seat.IsReserved = false
			seat.IsAvailable = true
			tx.Save(&seat)
		}
	}

	tx.Commit()
}

func SendOtp(phone string) (string, error) {

	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	serviceId := os.Getenv("TWILIO_SERVICE_SID")

	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	to := "+91" + phone
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceId, params)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("otp failed to generate")
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
		return *resp.Sid, nil
	}
}

func CheckOtp(phone, code string) error {

	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	serviceId := os.Getenv("TWILIO_SERVICE_SID")

	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	to := "+91" + phone
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(serviceId, params)

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("invalid otp")
	} else if *resp.Status == "approved" {
		return nil
	} else {
		return errors.New("invalid otp")
	}
}
