package main

import (
	"log"

	permify_gorm "github.com/Permify/permify-gorm"
	"github.com/cat-spmog/iothubgo/dal"
	model "github.com/cat-spmog/iothubgo/models"
	"github.com/cat-spmog/iothubgo/routers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Q dal.Query

func main() {
	dsn := "host=localhost user=postgres dbname=iotgo port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	// Auto-migrate models
	err = db.AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.SensorData{},
		&model.DataSink{},
		&model.DataSource{},
		&model.UserInteraction{},
	)

	// New initializer for Permify
	// If migration is true, it generate all tables in the database if they don't exist.
	permify, _ := permify_gorm.New(permify_gorm.Options{
		Migrate: true,
		DB:      db,
	})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := routers.Init(db)
	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
