package controllers

import (
	"app/dal"
)

type DashboardController struct {
	BaseController
}

func (c *DashboardController) Get() {
	// Fetch devices
	devices, err := dal.Device.Limit(10).Find()
	if err != nil {
		c.Data["Error"] = "Failed to load devices"
	}

	// Fetch platforms
	platforms, err := dal.Platform.Limit(10).Find()
	if err != nil {
		c.Data["Error"] = "Failed to load platforms"
	}

	c.Data["Devices"] = devices
	c.Data["Platforms"] = platforms
	c.TplName = "dashboard.html"
	c.Layout = "base/base.html"
}
