package routers

import (
	"time"

	"github.com/cat-spmog/iothubgo/controllers"
	"github.com/cat-spmog/iothubgo/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) *gin.Engine {
	// Initialize session store
	store := gormsessions.NewStore(db, true, []byte("secret"))
	r := gin.Default()
	// Set up session middleware
	r.Use(sessions.Sessions("sessions", store))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowAllOrigins: true,
		AllowCredentials:  true,
		AllowHeaders:      []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:     []string{"Content-Length"},
		AllowMethods:      []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowOriginFunc:   func(origin string) bool { return origin == "http://localhost:3000" || origin == "http://localhost:3001" },
		MaxAge: 12 * time.Hour,
	}))
	// CORS setup with specific origins instead of default wildcard
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	config.AllowOriginFunc = func(origin string) bool {
		// Allow specific origins
		return origin == "http://localhost:3000" || origin == "http://localhost:3001"
	}
	// Allow all methods and headers
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 300 // Cache preflight requests for 5 minutes (300 seconds)

	r.Use(cors.New(config))

	// Setup controllers
	userController := controllers.NewUserController(db)
	authController := controllers.NewAuthController(db)
	deviceController := controllers.NewDeviceController(db)
	integrationController := controllers.NewIntegrationController(db)
	snapshotController := controllers.NewSnapshotController(db)
	errorLogController := controllers.NewErrorLogController(db)
	roleController := controllers.NewRoleController(db)

	apiV1Group := r.Group("/api/v1")

	// Auth routes (public)
	authRoutes := apiV1Group.Group("/auth")
	{
		authRoutes.POST("/login", authController.PostLogin)
		authRoutes.GET("/me", authController.GetMe)
		authRoutes.GET("/logout", authController.Logout)
	}

	// Routes that require authentication
	protectedRoutes := apiV1Group.Group("", middleware.AuthMiddleware())

	// User routes (admin only)
	usersGroup := protectedRoutes.Group("/users", middleware.RequireRole(db, "admin", "user"))
	{
		usersGroup.GET("/", userController.GetUsers)
		usersGroup.GET("/:id", userController.GetUser)
		usersGroup.POST("/", userController.CreateUser)
		usersGroup.PUT("/:id", userController.UpdateUser)
		usersGroup.DELETE("/:id", userController.DeleteUser)
	}

	// Role management routes (admin only)
	rolesGroup := protectedRoutes.Group("/roles", middleware.RequireRole(db, "admin", "role"))
	{
		rolesGroup.GET("/", roleController.GetRoles)
		rolesGroup.GET("/:id", roleController.GetRole)
		rolesGroup.POST("/", roleController.CreateRole)
		rolesGroup.PUT("/:id", roleController.UpdateRole)
		rolesGroup.DELETE("/:id", roleController.DeleteRole)

		// User-role assignment routes
		rolesGroup.POST("/users/:userId/roles/:roleId", roleController.AssignRoleToUser)
		rolesGroup.DELETE("/users/:userId/roles/:roleId", roleController.RemoveRoleFromUser)
		rolesGroup.GET("/users/:userId/roles", roleController.GetUserRoles)
	}

	// Device routes (requires write permission for creating/updating, read for viewing)
	devicesGroup := protectedRoutes.Group("/devices")
	{
		devicesGroup.POST("/", middleware.RequireRole(db, "write", "device"), deviceController.CreateDevice)
		devicesGroup.PUT("/:id", middleware.RequireRole(db, "write", "device"), deviceController.UpdateDevice)
		devicesGroup.DELETE("/:id", middleware.RequireRole(db, "admin", "device"), deviceController.DeleteDevice)
		devicesGroup.GET("/:id", middleware.RequireRole(db, "read", "device"), deviceController.GetDevice)
		devicesGroup.GET("/", middleware.RequireRole(db, "read", "device"), deviceController.GetAllDevices)
	}

	// Integration routes (requires appropriate permissions)
	integrationsGroup := protectedRoutes.Group("/integrations")
	{
		integrationsGroup.POST("/", middleware.RequireRole(db, "write", "integration"), integrationController.CreateIntegration)
		integrationsGroup.PUT("/:id", middleware.RequireRole(db, "write", "integration"), integrationController.UpdateIntegration)
		integrationsGroup.DELETE("/:id", middleware.RequireRole(db, "admin", "integration"), integrationController.DeleteIntegration)
		integrationsGroup.GET("/:id", middleware.RequireRole(db, "read", "integration"), integrationController.GetIntegration)
		integrationsGroup.GET("/", middleware.RequireRole(db, "read", "integration"), integrationController.GetAllIntegrations)
	}

	// Snapshot routes (requires appropriate permissions)
	snapshotsGroup := protectedRoutes.Group("/snapshots")
	{
		snapshotsGroup.POST("/", middleware.RequireRole(db, "write", "snapshot"), snapshotController.CreateSnapshot)
		snapshotsGroup.PUT("/:id", middleware.RequireRole(db, "write", "snapshot"), snapshotController.UpdateSnapshot)
		snapshotsGroup.DELETE("/:id", middleware.RequireRole(db, "admin", "snapshot"), snapshotController.DeleteSnapshot)
		snapshotsGroup.GET("/:id", middleware.RequireRole(db, "read", "snapshot"), snapshotController.GetSnapshot)
		snapshotsGroup.GET("/", middleware.RequireRole(db, "read", "snapshot"), snapshotController.GetAllSnapshots)
	}

	// ErrorLog routes (requires admin permissions)
	errorLogsGroup := protectedRoutes.Group("/error_logs", middleware.RequireRole(db, "admin", "errorlog"))
	{
		errorLogsGroup.POST("/", errorLogController.CreateErrorLog)
		errorLogsGroup.PUT("/:id", errorLogController.UpdateErrorLog)
		errorLogsGroup.DELETE("/:id", errorLogController.DeleteErrorLog)
		errorLogsGroup.GET("/:id", errorLogController.GetErrorLog)
		errorLogsGroup.GET("/", errorLogController.GetAllErrorLogs)
	}

	return r
}
