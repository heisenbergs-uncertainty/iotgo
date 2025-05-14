package main

import (
	"app/dal"
	"app/middleware"
	"app/model"
	_ "app/routers"
	"app/seed"
	"log"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres dbname=iotgo sslmode=disable"), &gorm.Config{})
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

	beego.SetStaticPath("/static", "static")

	beego.BeeApp.InsertFilter("/*", beego.BeforeRouter, middleware.WebAuthMiddleware)

	beego.Run()
}
