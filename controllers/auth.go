package controllers

import (
	"iotgo/models"

	"github.com/beego/beego/v2/server/web"
)

type AuthController struct {
	web.Controller
}

func (c *AuthController) GetLogin() {
	c.Data["Title"] = "Login"
	c.TplName = "auth/login.tpl"
}

func (c *AuthController) PostLogin() {
	username := c.GetString("username")
	password := c.GetString("password")
	user, err := models.AuthenticateUser(username, password)
	if err != nil || user == nil {
		c.Ctx.WriteString("Invalid credentials")
		return
	}
	c.SetSession("user_id", user.Id)
	c.SetSession("role", user.Role)
	c.Redirect("/devices", 302)
}

func (c *AuthController) Logout() {
	c.DestroySession()
	c.Redirect("/login", 302)
}
