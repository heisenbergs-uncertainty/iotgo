package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"size:100;uniqueIndex"`
	Description string
	AccessMode  string  `gorm:"size:50"` // "read", "write", "admin"
	AccessType  string  `gorm:"size:50"` // Resource type like "device", "integration", "user", etc. or "*" for all
	Users       []*User `gorm:"many2many:user_roles;"`
}

// IsAdmin checks if the role has admin privileges
func (r *Role) IsAdmin() bool {
	return r.AccessMode == "admin"
}

// ResourceSpecific checks if the role applies to a specific resource
func (r *Role) ResourceSpecific() bool {
	return r.AccessType != "*"
}

// ValidateAccessMode ensures that the access mode is valid
func (r *Role) ValidateAccessMode() bool {
	validModes := map[string]bool{
		"read":  true,
		"write": true,
		"admin": true,
	}
	return validModes[r.AccessMode]
}
