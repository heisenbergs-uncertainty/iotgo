package main

import (
	_ "github.com/lib/pq"
	"iotgo/models"
	_ "iotgo/routers"
	"iotgo/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Set Logs
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/iotgo.log","level":7}`)

	// Register the PostgreSQL database
	orm.RegisterDriver("postgres", orm.DRPostgres)
	orm.RegisterDataBase("default", "postgres", "postgres://postgres@localhost/spmog-iot?sslmode=disable")
	orm.RunSyncdb("default", false, true)

	// Register Template functions
	utils.RegisterTemplateFuncs()
}

func seedAdminUser() {
	if _, err := models.GetUserByUsername("admin"); err == nil {
		return // Admin already exists
	}
	if err := models.AddUser("admin", "admin123", "admin"); err != nil {
		panic("Failed to seed admin user: " + err.Error())
	}
}

func main() {
	debug, err := config.Bool("debug")
	if err == nil {
		// Enable ORM debug mode
		orm.Debug = debug
	}

	if beego.BConfig.RunMode == "dev" {
		seedAdminUser()
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	// Enable admin service
	beego.BConfig.Listen.EnableAdmin = true
	beego.BConfig.Listen.AdminAddr = "localhost"
	beego.BConfig.Listen.AdminPort = 8088

	// Initialize health checks
	initHealthChecks()

	beego.Run()
}
