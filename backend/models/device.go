package model

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	DeviceName string `gorm:"type:varchar(50);not null"`
	Location   string `gorm:"type:varchar(100)"`
	Status     string `gorm:"type:varchar(20)"`
	Type       string `gorm:"type:varchar(50)"`
}
