package controllers

import (
	"html/template"
	"iotgo/models"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	admin "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/pagination"
)

// AdminController handles all admin-related operations
type AdminController struct {
	BaseController
}

// Prepare executes before each action in AdminController
func (c *AdminController) Prepare() {
	c.BaseController.Prepare()

	// Check if user is logged in and has admin role
	username := c.Data["User"]
	role := c.Data["UserRole"]

	if username == nil {
		c.Redirect("/login", 302)
		return
	}

	if role != "admin" {
		c.Redirect("/", 302)
		return
	}
}

func (c *AdminController) Dashboard() {
	start := time.Now()
	// Fetch necessary data for the dashboard
	userCount, err := models.GetUserCount()
	if err != nil {
		c.Data["UserCount"] = "Error fetching user count"
	} else {
		c.Data["UserCount"] = userCount
	}

	deviceCount, err := models.GetDeviceCount()
	if err != nil {
		c.Data["DeviceCount"] = "Error fetching device count"
	} else {
		c.Data["DeviceCount"] = deviceCount
	}

	// Fetch devices
	devices, err := models.GetAllDevices()
	if err != nil {
		c.Data["Error"] = "Error fetching devices: " + err.Error()
		c.Data["Redirect"] = "/admin"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin" // Retry the same action
		logs.Error("Error in Dashboard:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/admin", "&AdminController.Dashboard", time.Since(start))
		return
	}

	// Fetch request statistics
	stats := strings.Builder{}
	// beego.StatisticsMap.GetMap(stats)
	c.Data["RequestStats"] = stats.String()

	// Fetch health check status
	healthStatus := &strings.Builder{}
	// admin.GetHealthCheck(healthStatus)
	c.Data["HealthStatus"] = healthStatus.String()

	c.Data["Devices"] = devices
	c.Data["Title"] = "Admin Dashboard"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "admin/dashboard.html"

	admin.StatisticsMap.AddStatistics("GET", "/admin", "&AdminController.Dashboard", time.Since(start))
}

// GetUsers shows the user management page with pagination
func (c *AdminController) GetUsers() {
	start := time.Now()

	o := orm.NewOrm()
	// Get total user count
	total, err := models.GetUserCount()
	if err != nil {
		c.Data["Error"] = "Error fetching users: " + err.Error()
		c.Data["Redirect"] = "/admin"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users"
		logs.Error("Error in GetUsers:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/admin/users", "&AdminController.GetUsers", time.Since(start))
		return
	}

	// Initialize pagination
	perPage := 10
	paginator := pagination.SetPaginator(c.Ctx, perPage, total)

	// Fetch users for the current page
	var users []*models.User
	offset := (paginator.Page() - 1) * perPage
	_, err = o.QueryTable("user").Limit(perPage).Offset(int64(offset)).All(&users)
	if err != nil {
		c.Data["Error"] = "Error fetching users: " + err.Error()
		c.Data["Redirect"] = "/admin"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users"
		logs.Error("Error in GetUsers:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/admin/users", "&AdminController.GetUsers", time.Since(start))
		return
	}

	// Set pagination data
	c.Data["paginator"] = paginator

	// Get flash message from cookie if exists
	flash := c.Ctx.GetCookie("flash")
	if flash != "" {
		c.Data["Flash"] = flash
		c.Ctx.SetCookie("flash", "", -1) // Clear the flash cookie
	}

	c.Data["Title"] = "User Management"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "admin/users.html"

	admin.StatisticsMap.AddStatistics("GET", "/admin/users", "&AdminController.GetUsers", time.Since(start))
}

// CreateUser handles the creation of a new user
func (c *AdminController) CreateUser() {
	start := time.Now()

	username := c.GetString("username")
	password := c.GetString("password")
	role := c.GetString("role")

	// Create validation instance
	valid := validation.Validation{}
	valid.Required(username, "username").Message("Username is required")
	valid.MinSize(username, 3, "username").Message("Username must be at least 3 characters")
	valid.Required(password, "password").Message("Password is required")
	valid.MinSize(password, 6, "password").Message("Password must be at least 6 characters")
	valid.Required(role, "role").Message("Role is required")
	// valid.Match(role, `^(admin|user)$`, "role").Message("Role must be either 'admin' or 'user'")

	if valid.HasErrors() {
		c.Data["Errors"] = valid.Errors
		c.Data["Title"] = "User Management"
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
		c.TplName = "admin/users.html"
		admin.StatisticsMap.AddStatistics("POST", "/admin/users/create", "&AdminController.CreateUser", time.Since(start))
		return
	}

	if err := models.AddUser(username, password, role); err != nil {
		c.Data["Error"] = "Failed to add user: " + err.Error()
		c.Data["Redirect"] = "/admin/users"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users"
		logs.Error("Error in CreateUser:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/admin/users/create", "&AdminController.CreateUser", time.Since(start))
		return
	}

	c.SetFlash("User created successfully")
	c.Redirect("/admin/users", 302)
	admin.StatisticsMap.AddStatistics("POST", "/admin/users/create", "&AdminController.CreateUser", time.Since(start))
}

func (c *AdminController) EditUser() {
	start := time.Now()

	id, err := c.GetInt(":id")
	if err != nil {
		c.Data["Error"] = "Invalid user ID"
		c.Data["Redirect"] = "/admin/users"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users/edit/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in EditUser:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/admin/users/edit/"+strconv.Itoa(id), "&AdminController.EditUser", time.Since(start))
		return
	}

	user, err := models.GetUserById(id)
	if err != nil || user == nil {
		c.Data["Error"] = "User not found"
		c.Data["Redirect"] = "/admin/users"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users/edit/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in EditUser:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/admin/users/edit/"+strconv.Itoa(id), "&AdminController.EditUser", time.Since(start))
		return
	}

	if c.Ctx.Input.Method() == "POST" {
		username := c.GetString("username")
		role := c.GetString("role")

		// Server-side validation
		valid := validation.Validation{}
		valid.Required(username, "username").Message("Username is required")
		valid.MinSize(username, 3, "username").Message("Username must be at least 3 characters")
		valid.Required(role, "role").Message("Role is required")
		// valid.Match(role, `^(admin|user)$`, "role").Message("Role must be either 'admin' or 'user'")

		if valid.HasErrors() {
			c.Data["Errors"] = valid.Errors
			c.Data["User"] = user
			c.Data["Title"] = "Edit User"
			c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
			c.TplName = "admin/edit_user.html"
			admin.StatisticsMap.AddStatistics("POST", "/admin/users/edit/"+strconv.Itoa(id), "&AdminController.EditUser", time.Since(start))
			return
		}

		user.Username = username
		user.Role = role

		if err := models.UpdateUser(user); err != nil {
			c.Data["Error"] = "Failed to update user: " + err.Error()
			c.Data["Redirect"] = "/admin/users"
			c.Data["Title"] = "Error"
			c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
			c.Data["RetryURL"] = "/admin/users/edit/" + c.Ctx.Input.Param(":id")
			logs.Error("Error in EditUser:", err)
			c.TplName = "error.html"
			admin.StatisticsMap.AddStatistics("POST", "/admin/users/edit/"+strconv.Itoa(id), "&AdminController.EditUser", time.Since(start))
			return
		}

		c.SetFlash("User updated successfully")
		c.Redirect("/admin/users", 302)
		admin.StatisticsMap.AddStatistics("POST", "/admin/users/edit/"+strconv.Itoa(id), "&AdminController.EditUser", time.Since(start))
		return
	}

	c.Data["User"] = user
	c.Data["Title"] = "Edit User"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "admin/edit_user.html"
	admin.StatisticsMap.AddStatistics("GET", "/admin/users/edit/"+strconv.Itoa(id), "&AdminController.EditUser", time.Since(start))
}

// DeleteUser handles the deletion of a user
func (c *AdminController) DeleteUser() {
	start := time.Now()

	id, err := c.GetInt(":id")
	if err != nil {
		c.Data["Error"] = "Invalid user ID"
		c.Data["Redirect"] = "/admin/users"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users/delete/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in DeleteUser:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/admin/users/delete/"+strconv.Itoa(id), "&AdminController.DeleteUser", time.Since(start))
		return
	}

	if err := models.DeleteUser(id); err != nil {
		c.Data["Error"] = "Failed to delete user: " + err.Error()
		c.Data["Redirect"] = "/admin/users"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/admin/users/delete/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in DeleteUser:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/admin/users/delete/"+strconv.Itoa(id), "&AdminController.DeleteUser", time.Since(start))
		return
	}

	c.SetFlash("User deleted successfully")
	c.Redirect("/admin/users", 302)
	admin.StatisticsMap.AddStatistics("POST", "/admin/users/delete/"+strconv.Itoa(id), "&AdminController.DeleteUser", time.Since(start))
}
