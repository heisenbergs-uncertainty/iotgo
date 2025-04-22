package model

import "gorm.io/gorm"

// DataSource represents a data source in the system
type DataSink struct {
	gorm.Model
	SinkType          string `gorm:"type:varchar(20);not null"`
	ConnectionDetails string `gorm:"type:jsonb;not null"`
	Status            string `gorm:"type:varchar(20)"`
}