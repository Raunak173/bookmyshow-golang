package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null" validate:"required,min=2,max=50"`
	Email       string `json:"email" gorm:"not null;unique" validate:"email,required"`
	Password    string `json:"password" gorm:"not null"  validate:"required"`
	PhoneNumber string `json:"phone_number" gorm:"not null" validate:"required"`
	Otp         string `json:"otp"`
	IsAdmin     bool   `json:"isAdmin"`
}
