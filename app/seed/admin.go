package seed

import (
	"app/model"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seedAdminUser creates an admin user and an associated API key if they don't exist
func AdminUser(db *gorm.DB) {
	// Check if admin user exists
	var count int64
	db.Model(&model.User{}).Where("role = ?", "admin").Count(&count)

	// If no admin user exists, create one
	if count == 0 {
		adminUser := model.User{
			Email:    "admin@iotgo.local",
			Name:     "Admin",
			Role:     "admin",
			Metadata: "{}",
		}

		// Set password for admin user
		err := adminUser.HashPassword("admin123") // Consider using a more secure password
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}

		// Create admin user in database
		result := db.Create(&adminUser)
		if result.Error != nil {
			log.Printf("Failed to create admin user: %v", result.Error)
			return
		}

		// Create superuser API key
		token := uuid.New().String()
		hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash API key: %v", err)
			return
		}

		metadata, _ := json.Marshal(map[string]interface{}{
			"scopes": []string{"read", "write", "superuser"},
		})

		apiKey := model.ApiKey{
			Name:      "Admin Superuser Key",
			KeyID:     uuid.New().String(),
			Token:     token,
			KeyHash:   string(hash),
			IsActive:  true,
			UserID:    adminUser.ID,
			Metadata:  string(metadata),
			ExpiresAt: nil, // No expiration for admin key
		}

		if err := db.Create(&apiKey).Error; err != nil {
			log.Printf("Failed to create admin API key: %v", err)
			return
		}

		log.Printf("Admin user created successfully with email: %s and API key: %s", adminUser.Email, token)
	} else {
		log.Print("Admin user already exists, skipping creation")
	}
}
