package models

import (
	"gorm.io/gorm"
)

type Venue struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Location string `json:"location" gorm:"not null"`

	//Many venues will have multiple movies
	Movies []Movie `gorm:"many2many:movie_venues;"`

	//One venue can have many show timings
	ShowTimes []ShowTime `json:"show_times"`
}

type ShowTime struct {
	gorm.Model
	Timing string `json:"timing"`

	MovieID uint  `json:"movie_id"`
	Movie   Movie `json:"movie"`

	VenueID uint  `json:"venue_id"`
	Venue   Venue `json:"venue"`
}
