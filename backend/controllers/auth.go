package controllers

import (
	"app/dal"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	web.Controller
}

type Userlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Login() {
	user := &Userlogin{}
	if err := c.Ctx.BindJSON(user); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid request"}
		c.ServeJSON()
		return
	}

	logs.Info("Login attempt with email:", user.Email)

	if user.Email == "" || user.Password == "" {
		logs.Error("Email or password is empty")
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Email and password are required"}
		c.ServeJSON()
		return
	}

	q := dal.Q
	dbUser, err := q.User.FindByEmail(user.Email)
	if err != nil || dbUser.ID == 0 {
		logs.Error("User not found:", user.Email)
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid email or password"}
		c.ServeJSON()
		return
	}

	if err := dbUser.CheckPassword(user.Password); err != nil {
		logs.Error("Password mismatch for user:", user.Email)
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = map[string]string{"error": "Invalid email or password"}
		c.ServeJSON()
		return
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": dbUser.ID,
		"iss": "iotgo",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign with secret key (store securely, e.g., in environment variable)
	secretKey := []byte("your-secret-key") // Replace with env variable
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]string{"error": "Failed to generate token"}
		c.ServeJSON()
		return
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{"user": map[string]interface{}{"id": fmt.Sprint(dbUser.ID), "role": "superuser", "email": user.Email}, "token": tokenString}
	c.ServeJSON()
}
