package model
import "gorm.io/gorm"
import "time"

type UserInteraction struct {
	gorm.Model
	UserID     uint      `gorm:"index"`
	DeviceID   uint      `gorm:"index"`
	SourceID   uint      `gorm:"index"`
	SinkID     uint      `gorm:"index"`
	Action     string    `gorm:"type:varchar(20);not null"`
	Timestamp  time.Time `gorm:"index"`
	Device     Device    `gorm:"foreignKey:DeviceID"`
	DataSource DataSource `gorm:"foreignKey:SourceID"`
	DataSink   DataSink   `gorm:"foreignKey:SinkID"`
}