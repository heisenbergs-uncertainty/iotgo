package model
import "gorm.io/gorm"
import "time"

// SensorData represents a record of sensor data in the system
type SensorData struct {
	gorm.Model
	DeviceID   uint      `gorm:"index"`
	SourceID   uint      `gorm:"index"`
	SinkID     uint      `gorm:"index"`
	Timestamp  time.Time `gorm:"index"`
	SensorType string    `gorm:"type:varchar(50)"`
	Value      float64
	Device     Device    `gorm:"foreignKey:DeviceID"`
	DataSource DataSource `gorm:"foreignKey:SourceID"`
	DataSink   DataSink   `gorm:"foreignKey:SinkID"`
}