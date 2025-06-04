package main

import (
	"app/dal"
	"app/model"
	_ "app/routers"
	"app/seed"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	// Retrieve database connection details from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "iotgo")
	dbPortStr := getEnv("DB_PORT", "5432")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT environment variable: %v", err)
	}

	log.Printf("Attempting to connect to DB with parameters: HOST=%s, USER=%s, DBNAME=%s, PORT=%s", dbHost, dbUser, dbName, dbPortStr) //

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create tables
	db.AutoMigrate(&model.User{}, &model.Device{}, &model.ValueStream{}, &model.ApiKey{}, &model.Platform{}, &model.UserInteraction{}, &model.Site{}, &model.Resource{}, &model.DevicePlatform{})

	dal.SetDefault(db)

	// Seed admin user if no admin exists
	seed.AdminUser(db)

	// Initialize session
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionProvider = "memory"
	beego.BConfig.WebConfig.Session.SessionName = "iotgosessionid"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600

	beego.BConfig.CopyRequestBody = true

	beego.BConfig.Listen.EnableAdmin = true

	log := logs.NewLogger(10000)
	log.SetLogger("console")

	beego.Run()
}
