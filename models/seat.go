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

type Order struct {
	gorm.Model
	UserID     uint    `json:"user_id"`
	ShowTimeID uint    `json:"showtime_id"`
	TotalPrice float32 `json:"total_price"`

	// One order can have multiple seats
	Seats []Seat `json:"seats" gorm:"many2many:order_seats;"`
}
