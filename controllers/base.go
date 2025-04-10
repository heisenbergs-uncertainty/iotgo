package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

// BaseController is the base controller for all controllers
type BaseController struct {
	web.Controller
}

// Prepare is called prior to the controller method
func (c *BaseController) Prepare() {
	c.EnableXSRF = true
	c.XSRFExpire = 7200
	c.Data["User"] = c.GetSession("username")
	if role, ok := c.GetSession("role").(string); ok {
		c.Data["UserRole"] = role
	} else {
		c.Data["UserRole"] = ""
	}
	if username, ok := c.GetSession("username").(string); ok {
		c.Data["Username"] = username
	} else {
		c.Data["Username"] = "Guest"
	}
	// Update last activity timestamp in session
	c.Layout = "layout.html"
}

// SetFlash sets a flash message to be displayed on the next page
func (c *BaseController) SetFlash(msg string) {
	c.Ctx.SetCookie("flash", msg, 3)
}
