package models

import "gorm.io/gorm"

type Seat struct {
	gorm.Model

	SeatNumber  string  `json:"seat_number" gorm:"not null"`
	IsReserved  bool    `json:"isReserved"`
	IsBooked    bool    `json:"isBooked"`
	IsAvailable bool    `json:"isAvailable"`
	Price       float32 `json:"price"`

	//One seat belongs to one showtime
	ShowTimeID uint     `json:"showtime_id"`
	ShowTime   ShowTime `json:"showtime" gorm:"foreignKey:ShowTimeID"`
}
