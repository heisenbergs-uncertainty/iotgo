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
	Email            string            `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password         string            `gorm:"size:256;not null" json:"password"`
	Name             string            `gorm:"size:100;not null" json:"name"`
	Role             string            `gorm:"size:50;not null" json:"role"` // Admin, User, Guest, etc.
	Avatar           *string           `gorm:"size:255" json:"avatar"`
	Bio              *string           `gorm:"type:text" json:"bio"`
	UserInteractions []UserInteraction `gorm:"foreignKey:UserID" json:"user_interactions"`
	ApiKeys          []ApiKey          `gorm:"foreignKey:UserID" json:"api_keys"`
	LastLogin        *time.Time        `gorm:"type:timestamp with time zone" json:"last_login"`
	Metadata         string            `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for additional user metadata
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
	UserID   uint      `gorm:"index;not null" json:"user_id"`
	Action   string    `gorm:"size:100;not null" json:"action"`
	Metadata string    `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for interaction details
	Occurred time.Time `gorm:"type:timestamp with time zone;not null" json:"occurred"`
}

// ApiKey represents an API key with metadata
type ApiKey struct {
	Model
	Name      string     `gorm:"size:100;not null" json:"name"`
	KeyID     string     `gorm:"uniqueIndex;size:36;not null" json:"key_id"` // UUID as string
	Token     string     `gorm:"size:256;not null" json:"token"`            // Actual API token
	KeyHash   string     `gorm:"size:256;index;not null" json:"key_hash"`      // Hashed token for secure storage
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	LastUsed  *time.Time `gorm:"type:timestamp with time zone" json:"last_used"`
	ExpiresAt *time.Time `gorm:"type:timestamp with time zone" json:"expires_at"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user"`
	Metadata  string     `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for API key metadata
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
	Name        string   `gorm:"size:100;index;not null" json:"name"`
	Description *string  `gorm:"type:text" json:"description"`
	Type        string   `gorm:"size:50;not null" json:"type"` // Manufacturing, Packaging, Logistics, etc.
	IsActive    bool     `gorm:"default:true" json:"is_active"`
	Devices     []Device `gorm:"foreignKey:ValueStreamID" json:"devices"`
	Metadata    string   `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for value stream metadata
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
	DeviceID    uint           `gorm:"primaryKey;index" json:"device_id"`
	PlatformID  uint           `gorm:"primaryKey;index;uniqueIndex:idx_device_platform_alias" json:"platform_id"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	DeviceAlias string         `gorm:"size:100;uniqueIndex:idx_device_platform_alias;not null" json:"device_alias"` // Unique per platform
	Metadata    string         `gorm:"type:jsonb;default:'{}'" json:"metadata"`                                 // JSON string for relationship metadata
}

// Platform represents an external system with metadata
type Platform struct {
	Model
	Name            string     `gorm:"size:100;index;not null" json:"name"`
	Type            string     `gorm:"size:50;not null" json:"type"` // CMMS, ERP, Database, Cloud, REST, OPC_UA, SDK, etc.
	ConnectionState string     `gorm:"size:50;default:'Disconnected'" json:"connection_state"`
	LastConnected   *time.Time `gorm:"type:timestamp with time zone" json:"last_connected"`
	OrganizationID  *int       `gorm:"index" json:"organization_id"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`
	Devices         []Device   `gorm:"many2many:device_platforms" json:"devices"`
	Resources       []Resource `gorm:"foreignKey:PlatformID" json:"resources"`
	Metadata        string     `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for platform-specific configuration
}

// Resource represents a specific interaction point for a platform with metadata
type Resource struct {
	Model
	PlatformID uint   `gorm:"index;not null" json:"platform_id"`
	Name       string `gorm:"size:100;index;not null" json:"name"`
	Type       string `gorm:"size:50;not null" json:"type"`        // e.g., 'rest_endpoint', 'opcua_node', 'sdk_method'
	Details    string `gorm:"type:jsonb;not null" json:"details"`     // JSON string with type-specific details
	Metadata   string `gorm:"type:jsonb;default:'{}'" json:"metadata"` // JSON string for resource metadata
}
