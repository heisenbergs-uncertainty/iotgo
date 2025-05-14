package main

import (
	"app/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Define custom query interfaces for optimized database operations
type UserQuerier interface {
	// SELECT * FROM @@table WHERE email = @email LIMIT 1
	FindByEmail(email string) (gen.T, error)
	// SELECT * FROM @@table WHERE role = @role
	FindByRole(role string) ([]gen.T, error)
	// SELECT * FROM @@table WHERE last_login IS NOT NULL
	FindActiveUsers() ([]gen.T, error)
}

type DeviceQuerier interface {
	// SELECT * FROM @@table WHERE site_id = @siteID
	FindBySiteID(siteID uint) ([]gen.T, error)
	// SELECT * FROM @@table WHERE value_stream_id = @valueStreamID
	FindByValueStreamID(valueStreamID uint) ([]gen.T, error)
	// SELECT d.* FROM @@table d JOIN device_platforms dp ON d.id = dp.device_id WHERE dp.platform_id = @platformID
	FindDevicesOnPlatform(platformID uint) ([]gen.T, error)
}

type PlatformQuerier interface {
	// SELECT * FROM @@table WHERE type = @platformType
	FindByType(platformType string) ([]gen.T, error)
	// SELECT * FROM @@table WHERE is_active = true
	FindActivePlatforms() ([]gen.T, error)
}

type ResourceQuerier interface {
	// SELECT * FROM @@table WHERE platform_id = @platformID
	FindByPlatformID(platformID uint) ([]gen.T, error)
	// SELECT * FROM @@table WHERE type = @resourceType
	FindByType(resourceType string) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./dal",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	db, err := gorm.Open(postgres.Open("host=localhost user=postgres dbname=iotgo sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	g.UseDB(db)

	// Apply basic DAO API for all models
	g.ApplyBasic(
		model.User{},
		model.UserInteraction{},
		model.ApiKey{},
		model.Site{},
		model.ValueStream{},
		model.Device{},
		model.DevicePlatform{},
		model.Platform{},
		model.Resource{},
	)

	// Apply custom query interfaces to respective models
	g.ApplyInterface(func(UserQuerier) {}, model.User{})
	g.ApplyInterface(func(DeviceQuerier) {}, model.Device{})
	g.ApplyInterface(func(PlatformQuerier) {}, model.Platform{})
	g.ApplyInterface(func(ResourceQuerier) {}, model.Resource{})

	// Generate the code
	g.Execute()
}
