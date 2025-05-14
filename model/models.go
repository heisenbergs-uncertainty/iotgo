package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp with time zone" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// User represents a system user with metadata
type User struct {
	Model
	Email            string            `gorm:"size:100;uniqueIndex;not null"`
	Password         string            `gorm:"size:256;not null"`
	Name             string            `gorm:"size:100;not null"`
	Role             string            `gorm:"size:50;not null"`
	Avatar           *string           `gorm:"size:255"`
	Bio              *string           `gorm:"type:text"`
	UserInteractions []UserInteraction `gorm:"foreignKey:UserID"`
	ApiKeys          []ApiKey          `gorm:"foreignKey:UserID"`
	LastLogin        *time.Time        `gorm:"type:timestamp with time zone"`
	Metadata         string            `gorm:"type:jsonb;default:'{}'"` // JSON string for additional user metadata
}

func (user *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// UserInteraction represents user actions with metadata
type UserInteraction struct {
	Model
	UserID   uint      `gorm:"index;not null"`
	Action   string    `gorm:"size:100;not null"`
	Metadata string    `gorm:"type:jsonb;default:'{}'"` // JSON string for interaction details
	Occurred time.Time `gorm:"type:timestamp with time zone;not null"`
}

// ApiKey represents an API key with metadata
type ApiKey struct {
	Model
	Name      string     `gorm:"size:100;not null"`
	KeyID     string     `gorm:"uniqueIndex;size:36;not null"` // UUID as string
	Token     string     `gorm:"size:256;not null"`            // Actual API token
	KeyHash   string     `gorm:"size:256;index;not null"`      // Hashed token for secure storage
	IsActive  bool       `gorm:"default:true"`
	LastUsed  *time.Time `gorm:"type:timestamp with time zone"`
	ExpiresAt *time.Time `gorm:"type:timestamp with time zone"`
	UserID    uint       `gorm:"index;not null"`
	User      User       `gorm:"foreignKey:UserID"`
	Metadata  string     `gorm:"type:jsonb;default:'{}'"` // JSON string for API key metadata
}

// Site represents a physical location with metadata
type Site struct {
	Model
	Name        string   `gorm:"size:100;index;not null" json:"name"`
	Description *string  `gorm:"type:text" json:"description"`
	Address     string   `gorm:"size:255;not null" json:"address"`
	City        string   `gorm:"size:100;not null" json:"city"`
	State       string   `gorm:"size:50;not null" json:"state"`
	Country     string   `gorm:"size:50;not null" json:"country"`
	Devices     []Device `gorm:"foreignKey:SiteID" json:"devices"`
	Metadata    string   `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for site metadata
}

// ValueStream represents a product line or production flow with metadata
type ValueStream struct {
	Model
	Name        string   `gorm:"size:100;index;not null"`
	Description *string  `gorm:"type:text"`
	Type        string   `gorm:"size:50;not null"` // Manufacturing, Packaging, Logistics, etc.
	IsActive    bool     `gorm:"default:true"`
	Devices     []Device `gorm:"foreignKey:ValueStreamID"`
	Metadata    string   `gorm:"type:jsonb;default:'{}'"` // JSON string for value stream metadata
}

// Device represents an IoT device with metadata
type Device struct {
	Model
	Name          string       `gorm:"size:100;index;not null" json:"name"`
	SiteID        *uint        `gorm:"index" json:"site_id"`
	Site          *Site        `gorm:"foreignKey:SiteID" json:"site"`
	ValueStreamID *uint        `gorm:"index" json:"value_stream_id"`
	ValueStream   *ValueStream `gorm:"foreignKey:ValueStreamID" json:"value_stream"`
	Platforms     []Platform   `gorm:"many2many:device_platforms" json:"platforms"`
	Metadata      string       `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for device metadata
}

// DevicePlatform represents the many-to-many relationship between devices and platforms
type DevicePlatform struct {
	DeviceID    uint           `gorm:"primaryKey;index"`
	PlatformID  uint           `gorm:"primaryKey;index;uniqueIndex:idx_device_platform_alias"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	DeviceAlias string         `gorm:"size:100;uniqueIndex:idx_device_platform_alias;not null"` // Unique per platform
	Metadata    string         `gorm:"type:jsonb;default:'{}'"`                                 // JSON string for relationship metadata
}

// Platform represents an external system with metadata
type Platform struct {
	Model
	Name            string     `gorm:"size:100;index;not null"`
	Type            string     `gorm:"size:50;not null"` // CMMS, ERP, Database, Cloud, REST, OPC_UA, SDK, etc.
	ConnectionState string     `gorm:"size:50;default:'Disconnected'"`
	LastConnected   *time.Time `gorm:"type:timestamp with time zone"`
	OrganizationID  *int       `gorm:"index"`
	IsActive        bool       `gorm:"default:true"`
	Devices         []Device   `gorm:"many2many:device_platforms"`
	Resources       []Resource `gorm:"foreignKey:PlatformID"`
	Metadata        string     `gorm:"type:jsonb;default:'{}'"` // JSON string for platform-specific configuration
}

// Resource represents a specific interaction point for a platform with metadata
type Resource struct {
	Model
	PlatformID uint   `gorm:"index;not null"`
	Name       string `gorm:"size:100;index;not null"`
	Type       string `gorm:"size:50;not null"`        // e.g., 'rest_endpoint', 'opcua_node', 'sdk_method'
	Details    string `gorm:"type:jsonb;not null"`     // JSON string with type-specific details
	Metadata   string `gorm:"type:jsonb;default:'{}'"` // JSON string for resource metadata
}
