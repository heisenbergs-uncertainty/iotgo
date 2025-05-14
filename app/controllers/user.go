package controllers

import (
	"app/dal"
	"app/model"
	"encoding/json"
	"errors"
	"strconv"
)

type UserController struct {
	BaseController
}

// GetAll retrieves all users with pagination, filtering, and sorting
func (c *UserController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	roleFilter := c.GetString("role")
	nameSort := c.GetString("sort", "name")

	q := dal.Q
	query := q.User
	if roleFilter != "" {
		query.Where(q.User.Role.Eq(roleFilter))
	}

	if nameSort == "name" {
		query.Order(q.User.Name.Asc())
	} else {
		query.Order(q.User.Name.Desc())
	}

	users, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	total, err := query.Count()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.PaginatedResponse(users, total, limit, offset, nil)
}

// Get retrieves a single user by ID
func (c *UserController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	user, err := q.User.Where(q.User.ID.Eq(uint(id))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(user, nil)
}

// Post creates a new user with validation
func (c *UserController) Post() {
	var user model.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	// Validate required fields
	if user.Email == "" || user.Name == "" || user.Role == "" || user.Password == "" {
		c.JSONResponse(nil, errors.New("email, name, role, and password are required"))
		return
	}

	// Hash password
	if err := user.HashPassword(user.Password); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if err := dal.User.Create(&user); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(user, nil)
}

// Put updates an existing user with validation
func (c *UserController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var user model.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	// Validate required fields
	if user.Email == "" || user.Name == "" || user.Role == "" {
		c.JSONResponse(nil, errors.New("email, name, and role are required"))
		return
	}

	user.ID = uint(id)
	if user.Password != "" {
		if err := user.HashPassword(user.Password); err != nil {
			c.JSONResponse(nil, err)
			return
		}
	}

	q := dal.Q
	info, err := q.User.Where(q.User.ID.Eq(uint(id))).Updates(&user)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(user, nil)
}

// Delete removes a user by ID
func (c *UserController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.User.Where(q.User.ID.Eq(uint(id))).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "User deleted successfully"}, nil)
}
