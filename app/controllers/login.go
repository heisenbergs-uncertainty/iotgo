package controllers

import (
	"app/dal"

	"github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	web.Controller
}

func (c *LoginController) Get() {
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	email := c.GetString("email")
	password := c.GetString("password")

	q := dal.Q
	user, err := q.User.FindByEmail(email)
	if err != nil || user.ID == 0 {
		c.Data["Error"] = "Invalid email or password"
		c.TplName = "login.html"
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.Data["Error"] = "Invalid email or password"
		c.TplName = "login.html"
		return
	}

	c.SetSession("user_id", user.ID)
	c.Redirect("/", 302)
}
