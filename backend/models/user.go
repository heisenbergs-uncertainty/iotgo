package model

import (
	"errors"

	"github.com/cat-spmog/iothubgo/query"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string `gorm:"size:255"`
	Email          string `gorm:"type:varchar(100);unique_index"`
	HashedPassword []byte
	Active         bool
	Roles          []Role `gorm:"many2many:user_roles;"`
}

// BeforeCreate hook to validate user roles
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Check if the user has valid roles
	for _, role := range u.Roles {
		if role.AccessMode == "" || role.AccessType == "" {
			return errors.New("invalid role: accessMode and accessType are required")
		}
	}
	return nil
}

// SetNewPassword sets a new hashed password for the user
// SetNewPassword sets a new hashed password for the user
func (u *User) UpdatePassword(c *gin.Context, passwordString string) error {
	if passwordString == "" {
		return errors.New("password cannot be empty")
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the in-memory password
	u.HashedPassword = bcryptPassword

	// Update in database if ID exists
	if u.ID != 0 {
		q := query.User
		_, err = q.WithContext(c).Select(q.HashedPassword).
			Where(q.ID.Eq(u.ID)).
			Update(q.HashedPassword, bcryptPassword)
		return err
	}
	return nil
}

// CheckPassword verifies if the provided password matches the user's hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
	return err == nil
}

// HasPermission checks if the user has the required permission level for a specific resource type
func (u *User) HasPermission(requiredLevel string, resourceType string) bool {
	if len(u.Roles) == 0 {
		return false
	}

	for _, role := range u.Roles {
		// Wildcard resource type means the role applies to all resources
		if role.AccessType == "*" || role.AccessType == resourceType {
			if CheckPermission(role.AccessMode, requiredLevel) {
				return true
			}
		}
	}
	return false
}

// GetHighestPermissionLevel returns the highest permission level the user has
func (u *User) GetHighestPermissionLevel() string {
	if len(u.Roles) == 0 {
		return ""
	}

	// Define access levels - higher number means higher access
	accessLevels := map[string]int{
		"read":  1,
		"write": 2,
		"admin": 3,
	}

	highestLevel := 0
	highestLevelName := ""

	for _, role := range u.Roles {
		if level, exists := accessLevels[role.AccessMode]; exists && level > highestLevel {
			highestLevel = level
			highestLevelName = role.AccessMode
		}
	}

	return highestLevelName
}

// CreateUser creates a new user in the database
func CreateUser(db *gorm.DB, name *string, password *string, email *string) (uint, error) {
	if name == nil || password == nil || email == nil {
		return 0, errors.New("name, password, and email are required")
	}

	u := User{
		Name:   *name,
		Email:  *email,
		Active: true,
	}

	if err := u.SetNewPassword(*password); err != nil {
		return 0, err
	}

	// Create user in database
	if err := db.Create(&u).Error; err != nil {
		return 0, err
	}

	return u.ID, nil
}

// AuthenticateUser verifies a user's credentials and returns the User if successful
func AuthenticateUser(db *gorm.DB, email string, password string) (*User, error) {
	var user User
	if err := db.Preload("Roles").Where("email = ? AND active = ?", email, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Compare provided password with stored hash
	if !user.CheckPassword(password) {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

// CheckPermission checks if the provided access mode meets the required level
func CheckPermission(userAccessMode, requiredLevel string) bool {
	// Define access levels - higher number means higher access
	accessLevels := map[string]int{
		"read":  1,
		"write": 2,
		"admin": 3,
	}

	userLevel, userExists := accessLevels[userAccessMode]
	requiredLevelValue, requiredExists := accessLevels[requiredLevel]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevelValue
}
