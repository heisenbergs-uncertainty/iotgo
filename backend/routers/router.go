package routers

import (
	"app/controllers"
	"app/middleware"

	"github.com/beego/beego/v2/server/web"
)

func init() {
	authNs := web.NewNamespace("/auth",
		web.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
	)

	// API routes with auth filter
	apiNs := web.NewNamespace("/api",
		web.NSRouter("/users", &controllers.UserController{}, "get:GetAll;post:Post"),
		web.NSRouter("/users/:id", &controllers.UserController{}, "get:Get;put:Put;delete:Delete"),
		web.NSRouter("/users/:id/apikeys", &controllers.ApiKeyController{}, "get:GetAllByUser"),
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
		web.NSRouter("/platforms/:platform_id/resources", &controllers.ResourceController{}, "get:GetAll;post:Post"),
		// Resource routes
		web.NSRouter("/resources", &controllers.ResourceController{}, "get:GetAll;post:Post"),
		web.NSRouter("/resources/:id", &controllers.ResourceController{}, "get:Get;put:Put;delete:Delete"),
		web.NSRouter("/resources/:id/edit", &controllers.ResourceController{}, "get:GetResourceForEdit"),

		// Value Stream Routes
		web.NSRouter("/value-streams", &controllers.ValueStreamController{}, "get:GetAll;post:Post"),
		web.NSRouter("/value-streams/:id", &controllers.ValueStreamController{}, "get:Get;put:Put;delete:Delete"),

		// Resources Routes
		web.NSRouter("/resources/:id", &controllers.ResourceController{}, "get:Get;put:Put;delete:Delete"),
		web.NSRouter("/resources/:id/test", &controllers.ResourceController{}, "post:TestResource"),
	)
	apiNs.Filter("before", middleware.ApiAuthFilter)
	web.AddNamespace(apiNs, authNs)
}
