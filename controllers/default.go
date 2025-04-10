package controllers

import (
	"iotgo/models"
	"time"
)

// MainController handles the default routes
type MainController struct {
	BaseController
}

// Get handles the index page request
func (c *MainController) Get() {
	// Check if user is logged in
	if c.GetSession("username") == nil {
		c.Redirect("/login", 302)
		return
	}

	// Get dashboard data
	deviceCount, _ := models.GetDeviceCount()

	// Mock recent activities for demonstration
	recentActivities := []struct {
		Description string
		Time        string
	}{
		{"Device status updated", time.Now().Add(-5 * time.Minute).Format("15:04")},
		{"New user registered", time.Now().Add(-2 * time.Hour).Format("15:04")},
		{"System update completed", time.Now().Add(-24 * time.Hour).Format("Jan 02")},
	}

	c.Data["Title"] = "Dashboard"
	c.Data["DeviceCount"] = deviceCount
	c.Data["RecentActivities"] = recentActivities
	c.TplName = "index.html"
}
