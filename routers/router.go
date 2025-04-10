package routers

import (
	"iotgo/controllers"

	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

func adminFilter(ctx *context.Context) {
	role := ctx.Input.Session("role")
	if role != "admin" {
		ctx.Redirect(302, "/login")
	}
}

func authFilter(ctx *context.Context) {
	path := ctx.Request.URL.Path

	if path == "/login" || path == "/logout" || path == "favicon.ico" || path == "/static/*" {
		return
	}

	userId := ctx.Input.Session("user_id")
	if userId == nil {
		ctx.Redirect(302, "/login")
	}
}

func init() {
	// Public Routes
	beego.Router("/login", &controllers.LoginController{}, "get:Get;post:Post")
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")

	// Apply authentication filter to all UI routes except login/logout
	beego.InsertFilter("/*", beego.BeforeRouter, authFilter)

	// Dashboard Routes
	beego.Router("/", &controllers.MainController{}, "get:Get")

	// RESTful Device routes
	beego.Router("/devices", &controllers.WebDeviceController{}, "get:List;post:Create")
	beego.Router("/devices/new", &controllers.WebDeviceController{}, "get:New")
	beego.Router("/devices/:id/edit", &controllers.WebDeviceController{}, "get:Edit;post:Update")
	beego.Router("/devices/:id/delete", &controllers.WebDeviceController{}, "post:Delete")
	beego.Router("/devices/:id", &controllers.WebDeviceController{}, "get:Show")
	beego.Router("/devices/:device_id/integrations/:id/edit", &controllers.WebIntegrationController{}, "get:Edit;post:Update")
	beego.Router("/devices/:device_id/integrations/new", &controllers.WebIntegrationController{}, "get:New;post:Create")
	beego.Router("/devices/:device_id/integrations", &controllers.WebIntegrationController{}, "post:Create")
	beego.Router("/devices/:device_id/integrations/:id/delete", &controllers.WebIntegrationController{}, "post:Delete")
	beego.Router("/devices/:device_id/integrations/:integration_id/browse", &controllers.WebDeviceController{}, "get:BrowseNodes")
	//	beego.Router("/devices/:device_id/integrations/:integration_id/record", &controllers.RecordingController{}, "post:RecordData")
	beego.Router("/devices/:device_id/integrations/:integration_id/snapshot", &controllers.SnapshotController{}, "post:TakeSnapshot")

	// RESTful Integration routes
	beego.Router("/integrations", &controllers.WebIntegrationController{}, "get:List")
	beego.Router("/integrations/:id", &controllers.WebIntegrationController{}, "get:Show")

	// Admin routes with admin role filter
	beego.InsertFilter("/admin/*", beego.BeforeRouter, adminFilter)
	beego.Router("/admin", &controllers.AdminController{}, "get:Dashboard")
	beego.Router("/admin/users", &controllers.AdminController{}, "get:GetUsers")
	beego.Router("/admin/users/create", &controllers.AdminController{}, "post:CreateUser")
	beego.Router("/admin/users/edit/:id", &controllers.AdminController{}, "get:EditUser")
	beego.Router("/admin/users/delete/:id", &controllers.AdminController{}, "get:DeleteUser")

	// API routes (RESTful)
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/devices", &controllers.DeviceController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/devices/:id", &controllers.DeviceController{}, "get:GetOne;put:Put;delete:Delete"),
			beego.NSRouter("/devices/:device_id/integrations", &controllers.IntegrationController{}, "get:GetAllByDevice;post:Post"),
			beego.NSRouter("/devices/:device_id/integrations/:id", &controllers.IntegrationController{}, "get:GetOne;put:Put;delete:Delete"),
		),
	)
	beego.AddNamespace(ns)
}
