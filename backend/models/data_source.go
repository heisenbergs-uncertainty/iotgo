package model

import "gorm.io/gorm"

type DataSource struct {
	gorm.Model
	SourceType        string `gorm:"type:varchar(20);not null"`
	ConnectionDetails string `gorm:"type:jsonb;not null"`
	Status            string `gorm:"type:varchar(20)"`
}
