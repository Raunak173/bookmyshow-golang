package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null" validate:"required,min=2,max=50"`
	Description string `json:"desc" gorm:"not null"`
	Duration    string `json:"duration" gorm:"not null"` //in hours
}
