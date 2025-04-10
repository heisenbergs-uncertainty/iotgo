package controllers

import (
	"iotgo/models"
	"log"
)

// LoginController handles login-related actions
type LoginController struct {
	BaseController
}

// Get handles the login page request
func (c *LoginController) Get() {
	// Check if user is already logged in
	if c.GetSession("username") != nil {
		c.Redirect("/", 302)
		return
	}
	c.Data["Title"] = "Login"
	c.TplName = "auth/login.html"
	c.Layout = "" // Disable Navbar
}

// Post handles the login form submission
func (c *LoginController) Post() {
	username := c.GetString("username")
	password := c.GetString("password")

	user, err := models.AuthenticateUser(username, password)
	if err != nil {
		log.Printf("Login failed for user %s: %v", username, err)
		c.Data["Error"] = "Invalid username or password"
		c.Data["Title"] = "Login"
		c.TplName = "login.html"
		return
	}

	c.SetSession("username", username)
	c.SetSession("role", user.Role)
	c.SetSession("user_id", user.Id)
	c.Redirect("/", 302)
}

// Logout handles user logout
func (c *LoginController) Logout() {
	c.DestroySession()
	c.Redirect("/login", 302)
}
