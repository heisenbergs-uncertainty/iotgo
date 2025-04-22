package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for a password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AccessLevels defines permission hierarchy
var AccessLevels = map[string]int{
	"read":      1,
	"write":     2,
	"admin":     3,
	"superuser": 4,
	"root":      5,
}

// ErrUnauthorized is returned when a user doesn't have the required permissions
var ErrUnauthorized = errors.New("unauthorized: insufficient permissions")

// CheckPermission verifies if the provided access mode meets the required level
func CheckPermission(userAccessMode string, requiredLevel string) bool {
	userLevel, userOk := AccessLevels[userAccessMode]
	requiredVal, requiredOk := AccessLevels[requiredLevel]

	if !userOk || !requiredOk {
		return false
	}

	return userLevel >= requiredVal
}

// ValidAccessModes returns all valid access mode values
func ValidAccessModes() []string {
	modes := make([]string, 0, len(AccessLevels))
	for mode := range AccessLevels {
		modes = append(modes, mode)
	}
	return modes
}
