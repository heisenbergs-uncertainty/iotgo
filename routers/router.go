package routers

import (
	"app/controllers"
	"app/middleware"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/login", &controllers.LoginController{})
	web.Router("/logout", &controllers.LogoutController{})
	web.Router("/dashboard", &controllers.DashboardController{})
	web.Router("/devices", &controllers.DeviceController{}, "get:ListDevices")
	web.Router("/devices/:device_id/platforms", &controllers.DevicePlatformController{}, "get:ListAssociations")
	web.Router("/sites", &controllers.SiteController{}, "get:ListSites")
	web.Router("/platforms", &controllers.PlatformController{}, "get:ListPlatforms")
	web.Router("/value_streams", &controllers.ValueStreamController{}, "get:ListValueStreams")
	web.Router("/api_key", &controllers.ApiKeyController{}, "get:ManageKeys")
	web.Router("/api_key/generate", &controllers.ApiKeyController{}, "post:Generate")
	web.Router("/api_key/revoke/:key_id", &controllers.ApiKeyController{}, "post:Revoke")

	// API routes with auth filter
	apiNs := web.NewNamespace("/api",
		web.NSRouter("/users", &controllers.UserController{}, "get:GetAll;post:Post"),
		web.NSRouter("/users/:id", &controllers.UserController{}, "get:Get;put:Put;delete:Delete"),
		// Site routes
		web.NSRouter("/sites", &controllers.SiteController{}, "get:GetAll;post:Post"),
		web.NSRouter("/sites/:id", &controllers.SiteController{}, "get:Get;put:Put;delete:Delete"),

		// Device routes
		web.NSRouter("/devices", &controllers.DeviceController{}, "get:GetAll;post:Post"),
		web.NSRouter("/devices/:id", &controllers.DeviceController{}, "get:Get;put:Put;delete:Delete"),
		web.NSRouter("/devices/:device_id/platforms", &controllers.DevicePlatformController{}, "get:GetAll;post:Post"),
		web.NSRouter("/devices/:device_id/platforms/:platform_id", &controllers.DevicePlatformController{}, "delete:Delete"),

		// Platform routes
		web.NSRouter("/platforms", &controllers.PlatformController{}, "get:GetAll;post:Post"),
		web.NSRouter("/platforms/:id", &controllers.PlatformController{}, "get:Get;put:Put;delete:Delete"),
		web.NSRouter("/platforms/:id/devices/:device_id/data", &controllers.PlatformController{}, "get:FetchDeviceData"),

		// Resource routes
		web.NSRouter("/resources", &controllers.ResourceController{}, "get:GetAll;post:Post"),
		web.NSRouter("/resources/:id", &controllers.ResourceController{}, "get:Get;put:Put;delete:Delete"),

		// Value Stream Routes
		web.NSRouter("/value_streams", &controllers.ValueStreamController{}, "get:GetAll;post:Post"),
		web.NSRouter("/value_streams/:id", &controllers.ValueStreamController{}, "get:Get;put:Put;delete:Delete"),
	)
	apiNs.Filter("before", middleware.ApiAuthFilter)
	web.AddNamespace(apiNs)
}
