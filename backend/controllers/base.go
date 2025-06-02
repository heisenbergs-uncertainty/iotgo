package controllers

import (
	"app/dal"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/beego/beego/v2/server/web"
	"golang.org/x/time/rate"
)

type BaseController struct {
	web.Controller
}

var (
	limiter   = rate.NewLimiter(rate.Limit(10), 100) // 10 requests per second, burst of 100
	limiterMu sync.Mutex
)

// Prepare runs before every request to check authentication
func (c *BaseController) Prepare() {
	limiterMu.Lock()
	if !limiter.Allow() {
		limiterMu.Unlock()
		c.JSONResponse(nil, errors.New("rate limit exceeded"))
		c.Abort("429")
		return
	}
	limiterMu.Unlock()

	// Set default UserRole to avoid nil in templates
	c.Data["UserRole"] = ""

	if strings.HasPrefix(c.Ctx.Request.URL.Path, "/api") {
		return
	}

	if c.Ctx.Request.URL.Path != "/login" {
		c.Layout = "base/base.html"

		userID := c.GetSession("user_id")
		if userID == nil {
			// Redirect to login if no user_id in session
			c.Redirect("/login", 302)
			return
		}

		q := dal.Q
		user, err := q.User.Where(q.User.ID.Eq(userID.(uint))).First()
		if err != nil || user.ID == 0 {
			log.Printf("Failed to fetch user: %v", err)
			c.DestroySession()
			c.Redirect("/login", 302)
			return
		}

		// Restrict admin-only routes
		if strings.HasPrefix(c.Ctx.Request.URL.Path, "/admin") && user.Role != "admin" {
			c.Redirect("/", 403)
			return
		}

		c.Data["UserID"] = userID
		c.Data["UserRole"] = user.Role
	}
}

// JSONResponse standardizes API responses
func (c *BaseController) JSONResponse(data interface{}, err error) {
	if err != nil {
		log.Printf("API error: %v", err)
		c.Data["json"] = map[string]interface{}{"error": err.Error(), "code": 400}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"data": data, "code": 200}
	c.ServeJSON()
}

// PaginatedResponse generates a paginated response for a given array of items
func (c *BaseController) PaginatedResponse(items interface{}, total int64, limit int, offset int, err error) {
	if err != nil {
		log.Printf("Paginated API error: %v", err)
		c.JSONResponse(nil, err)
		return
	}

	response := map[string]interface{}{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	c.JSONResponse(response, nil)
}
