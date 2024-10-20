package models

import "gorm.io/gorm"

type Venue struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Location string `json:"location" gorm:"not null"`

	//Many venues will have multiple movies
	Movies []Movie `gorm:"many2many:movie_venues;"`
}
